package commands

import (
	"github.com/diamondburned/arikawa/v3/api"
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
		Restricted:  false,
		Usage:       "/help",
		Run: func(e *gateway.InteractionCreateEvent, data *discord.CommandInteraction) {

			res := api.InteractionResponse{
				Type: api.MessageInteractionWithSource,
				Data: &api.InteractionResponseData{
					Flags: api.EphemeralResponse,
					Embeds: &[]discord.Embed{
						{
							Title:       "Help!",
							Description: "Select a category to see available commands",
							Color:       utils.DefaultColour,
							Footer: &discord.EmbedFooter{
								Text: e.Member.User.Username,
								Icon: e.Member.User.AvatarURL(),
							},
							Timestamp: discord.NowTimestamp(),
						},
					},
					Components: discord.ComponentsPtr(
						&discord.ActionRowComponent{
							&discord.SelectComponent{
								CustomID:    "help_select_category",
								Placeholder: "Select a command category",
								Options: []discord.SelectOption{
									{
										Label:       "Miscellaneous",
										Value:       "misc",
										Description: "returns all miscellaneous commands",
										Default:     false,
									},
									{
										Label:       "Info",
										Value:       "info",
										Description: "returns all info commands",
										Default:     false,
									},
									{
										Label:       "Utility",
										Value:       "utility",
										Description: "returns all utility commands",
										Default:     false,
									},
									{
										Label:       "Debug",
										Value:       "debug",
										Description: "returns all debug commands",
										Default:     false,
									},
									{
										Label:       "Minecraft",
										Value:       "minecraft",
										Description: "returns all minecraft commands",
										Default:     false,
									},
								},
							},
						}),
				},
			}

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
		Restricted:  false,
		Run: func(e *gateway.InteractionCreateEvent, data *discord.CommandInteraction) {

			res := api.InteractionResponse{
				Type: api.MessageInteractionWithSource,
				Data: &api.InteractionResponseData{
					Flags: api.EphemeralResponse,
					Embeds: &[]discord.Embed{
						{
							Title:       "GoTDC Info",
							Description: "GoTDC is a bot Developed by king_ultron99 for the sole use in \"the Delinquents Club\" server. GoTDC will be used as a tool for server management and eventually as a discord-based gateway for the associated minecraft server",
							Color:       utils.DefaultColour,
							Footer: &discord.EmbedFooter{
								Text: e.Member.User.Username,
								Icon: e.Member.User.AvatarURL(),
							},
							Timestamp: discord.NowTimestamp(),
						},
					},
				},
			}

			if err := core.State.RespondInteraction(e.ID, e.Token, res); err != nil {
				logger.Error(err)
			}
		},
	}
}
