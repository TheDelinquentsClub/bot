package commands

import (
	"fmt"
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/kingultron99/tdcbot/core"
	"github.com/kingultron99/tdcbot/logger"
)

type Command struct {
	Name  string
	Description  string
	Usage string
	Run   func(e *gateway.InteractionCreateEvent, data *discord.CommandInteractionData)
}

var CommandsMap = make(map[string]Command)

func AddHandlers() {
	core.State.AddHandler(func(e *gateway.InteractionCreateEvent) {
		switch data := e.Data.(type) {
		case *discord.CommandInteractionData:
			logger.Info.Println(CommandsMap[data.Name])
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
		Name: "help",
		Description: "help",
		NoDefaultPermission: false,
	},
}

func Register(appID discord.AppID, guildID discord.GuildID) {
	_, err := core.State.BulkOverwriteGuildCommands(appID, guildID, commands)
	if err != nil {
		logger.Error.Println(fmt.Sprintf("Failed to overwrite commands in TDC with err: %v", err))
	}
}

func init() {
	CommandsMap["help"] = Command{
		Name:        "help",
		Description: "Returns a list of available commands",
		Usage:       "/help [page]",
		Run: func(e *gateway.InteractionCreateEvent, data *discord.CommandInteractionData) *api.InteractionResponseData {
			logger.Info.Println("Received a help interaction!")

			var res = api.InteractionResponseData{
				Embeds:
			}

			return res
		},
	}
}
