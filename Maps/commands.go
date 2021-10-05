package Maps

import (
	"fmt"
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/kingultron99/tdcbot/core"
	"github.com/kingultron99/tdcbot/logger"
	"github.com/kingultron99/tdcbot/utils"
)

func AddHandlers() {
	core.State.AddHandler(func(e *gateway.InteractionCreateEvent) {
		switch data := e.Data.(type) {
		case *discord.CommandInteractionData:
			if cmd, ok := MapCommands[data.Name]; ok {
				cmd.Run(e, data)
			}
		case *discord.ComponentInteractionData:
			if cmd, ok := MapComponents[data.CustomID]; ok {
				cmd.Run(e, data)
			}
		}
	})
}

func Register(appID discord.AppID, guildID discord.GuildID) {

	var commands []discord.Command

	for _, command := range MapCommands {
		commands = append(commands, discord.Command{
			Type:                discord.CommandType(1),
			Name:                command.Name,
			Description:         command.Description,
			Options:             command.Options,
			NoDefaultPermission: command.OwnerOnly,
		})
	}

	_, err := core.State.BulkOverwriteGuildCommands(appID, guildID, commands)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to overwrite commands in TDC with err: %v", err))
	}

	registeredCommands, err := core.State.GuildCommands(appID, guildID)
	if err != nil {
		logger.Error(err)
	}

	for _, command := range registeredCommands {
		logger.Debug(command, " Registered!")
		if command.NoDefaultPermission == true {
			core.State.BatchEditCommandPermissions(appID, guildID, []api.BatchEditCommandPermissionsData{
				{
					ID: command.ID,
					Permissions: []discord.CommandPermissions{
						{
							ID:         utils.MustSnowflakeEnv(core.Config.OwnerID),
							Type:       2,
							Permission: true,
						},
					},
				},
			})
			logger.Info("Successfully updated", command.Name, "permissions")
		}
	}
}
