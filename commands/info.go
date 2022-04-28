package commands

import (
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/kingultron99/tdcbot/core"
	"github.com/kingultron99/tdcbot/logger"
	"github.com/kingultron99/tdcbot/utils"
)

func init() {
	MapCommands["help"] = Command{
		Name:        "help",
		Description: "Returns a list of available commands",
		Group:       "info",
		OwnerOnly:   false,
		Usage:       "/help",
		Run: func(e *gateway.InteractionCreateEvent, data *discord.CommandInteractionData) {
			res := utils.NewEmbed().
				SetTitle("Help!").
				SetDescription("Select a category to see available commands").
				SetColor(utils.DefaultColour).
				SetFooter(e.Member.User.Username, e.Member.User.AvatarURL()).
				AddSelectComponent("help_select_category", "Select a command category", false).
				AddOption("Miscellaneous", "misc", "returns all miscellaneous commands", &discord.ButtonEmoji{}, false).
				AddOption("Info", "info", "returns all info commands", &discord.ButtonEmoji{}, false).
				AddOption("Utility", "utility", "returns all  utility commands", &discord.ButtonEmoji{}, false).
				AddOption("Debug", "debug", "returns all debug commands", &discord.ButtonEmoji{}, false).
				MakeSelectComponent().
				MakeResponse()

			if err := core.State.RespondInteraction(e.ID, e.Token, res); err != nil {
				logger.Error(err)
			}

		},
	}
	MapCommands["info"] = Command{
		Name:        "info",
		Description: "Information about GoTDC",
		Group:       "info",
		Usage:       "/info",
		OwnerOnly:   false,
		Run: func(e *gateway.InteractionCreateEvent, data *discord.CommandInteractionData) {
			res := utils.NewEmbed().
				SetTitle("GoTDC Info").
				SetDescription("GoTDC is a bot Developed by king_ultron99 for the sole use in \"the Delinquents Club\" server. GoTDC will be used as a tool for server management and eventually as a discord-based gateway for the associated minecraft server").
				SetColor(utils.DefaultColour).
				SetFooter(e.Member.User.Username, e.Member.User.AvatarURL()).
				AddURLButton("Source", "https://github.com/kingultron99/TDC-Bot").
				MakeResponse()

			if err := core.State.RespondInteraction(e.ID, e.Token, res); err != nil {
				logger.Error(err)
			}
		},
	}
}
