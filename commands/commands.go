package commands

import (
	"fmt"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/kingultron99/tdcbot/core"
	"github.com/kingultron99/tdcbot/logger"
)

func AddHandlers() {
	core.State.AddHandler(func(e *gateway.InteractionCreateEvent) {
		switch data := e.Data.(type) {
		case *discord.CommandInteractionData:
			if cmd, ok := CommandsMap[data.Name]; ok {
				cmd.Run(e, data)
			}
		case *discord.ComponentInteractionData:
			break
		}
	})
}

var commands = []discord.Command{
	{
		Name:                "help",
		Description:         "Returns the available commands",
		NoDefaultPermission: false,
	},
	{
		Name:                "stats",
		Description:         "Returns the current statistics and host system information of GoTDC",
		NoDefaultPermission: false,
	},
}

func Register(appID discord.AppID, guildID discord.GuildID) {
	_, err := core.State.BulkOverwriteGuildCommands(appID, guildID, commands)
	if err != nil {
		logger.Error.Println(fmt.Sprintf("Failed to overwrite commands in TDC with err: %v", err))
	}
}
