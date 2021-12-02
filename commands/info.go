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
		OwnerOnly:   false,
		Usage:       "/help",
		Run: func(e *gateway.InteractionCreateEvent, data *discord.CommandInteractionData) {
			res := api.InteractionResponse{
				Type: api.MessageInteractionWithSource,
				Data: &api.InteractionResponseData{
					Embeds: &[]discord.Embed{
						{
							Title:       "Help!",
							Description: "Select a category to see available commands!",
							Color:       utils.DefaultColour,
							Footer: &discord.EmbedFooter{
								Text: e.Member.User.Username,
								Icon: e.Member.User.AvatarURL(),
							},
							Timestamp: discord.NowTimestamp(),
						},
					},
					Components: &[]discord.Component{
						&discord.ActionRowComponent{
							Components: []discord.Component{
								&discord.SelectComponent{
									CustomID: "help_select_category",
									Options: []discord.SelectComponentOption{
										{
											Label:       "Miscellaneous",
											Value:       "misc",
											Description: "returns all miscellaneous commands",
										},
										{
											Label:       "Info",
											Value:       "info",
											Description: "returns all info commands",
										},
										{
											Label:       "Utility",
											Value:       "utility",
											Description: "returns all  utility commands",
										},
										{
											Label:       "Debug",
											Value:       "debug",
											Description: "returns all debug commands",
										},
										{
											Label:       "Music",
											Value:       "music",
											Description: "returns all music commands",
										},
									},
									Placeholder: "Select a command category!",
									Disabled:    false,
								},
							},
						},
					},
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
		OwnerOnly:   false,
		Run: func(e *gateway.InteractionCreateEvent, data *discord.CommandInteractionData) {
			res := api.InteractionResponse{
				Type: api.MessageInteractionWithSource,
				Data: &api.InteractionResponseData{
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
					Components: &[]discord.Component{
						&discord.ActionRowComponent{
							Components: []discord.Component{
								&discord.ButtonComponent{
									Label: "Source",
									URL:   "https://github.com/kingultron99/TDC-Bot",
									Style: discord.LinkButton,
								},
							},
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
