package commands

import (
	"fmt"
	"github.com/diamondburned/arikawa/v3/api"
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
				switch cmd.OwnerOnly {
				case true:
					if e.Member.User.ID.String() == core.Config.OwnerID {
						cmd.Run(e, data)
					} else {
						noPermission(e)
					}
				case false:
					cmd.Run(e, data)
				default:
					cmd.Run(e, data)
				}
			}
		case *discord.ComponentInteractionData:
			break
		}
	})
}

func Register(appID discord.AppID, guildID discord.GuildID) {
	_, err := core.State.BulkOverwriteGuildCommands(appID, guildID, commands)
	if err != nil {
		logger.Error.Println(fmt.Sprintf("Failed to overwrite commands in TDC with err: %v", err))
	}
}

func noPermission(e *gateway.InteractionCreateEvent) {

	res := api.InteractionResponse{
		Type: api.MessageInteractionWithSource,
		Data: &api.InteractionResponseData{
			Embeds: &[]discord.Embed{
				{
					Title:       "Insufficient Permissions!",
					Description: "You do not have permission to execute this command!",
					Color:       0xFF0000,
					Timestamp:   discord.NowTimestamp(),
					Footer: &discord.EmbedFooter{
						Text: e.Member.User.Username,
						Icon: e.Member.User.AvatarURL(),
					},
				},
			},
		},
	}

	if err := core.State.RespondInteraction(e.ID, e.Token, res); err != nil {
		logger.Error.Println(err)
	}
}
