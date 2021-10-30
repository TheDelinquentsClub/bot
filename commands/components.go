package commands

import (
	"fmt"
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/kingultron99/tdcbot/core"
	"github.com/kingultron99/tdcbot/logger"
	"github.com/kingultron99/tdcbot/utils"
)

func init() {

	MapComponents["help_select_category"] = Component{
		Run: func(e *gateway.InteractionCreateEvent, data *discord.ComponentInteractionData) {
			var fields, embedColour = GetCommands(data.Values[0])

			res := api.InteractionResponse{
				Type: api.UpdateMessage,
				Data: &api.InteractionResponseData{
					Embeds: &[]discord.Embed{
						{
							Title:       fmt.Sprint("Help! | ", data.Values[0]),
							Description: fmt.Sprint("Here are all the available commands in the ", data.Values[0], " category"),
							Fields:      fields,
							Color:       embedColour,
							Timestamp:   discord.NowTimestamp(),
							Footer: &discord.EmbedFooter{
								Text: fmt.Sprintf("Requested by %v#%v", e.Member.User.Username, e.Member.User.Discriminator),
								Icon: e.Member.User.AvatarURL(),
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

	//the next few component interaction functions are purely for the poll command
	MapComponents["first_option"] = Component{
		Run: func(e *gateway.InteractionCreateEvent, data *discord.ComponentInteractionData) {
			if Voted[e.Member.User.ID] != "" {
				logger.Warn(e.Member.User.Username, "has already voted!")
				res := api.InteractionResponse{
					Type: api.MessageInteractionWithSource,
					Data: &api.InteractionResponseData{
						Flags: api.EphemeralResponse,
						Embeds: &[]discord.Embed{
							{
								Title:       "You've already voted!!",
								Description: fmt.Sprintf("You voted for `%v`", Voted[e.Member.User.ID]),
								Color:       utils.DiscordRed,
							},
						},
					},
				}
				if err := core.State.RespondInteraction(e.ID, e.Token, res); err != nil {
					logger.Error(err)
				}
				return
			}
			res := api.InteractionResponse{
				Type: api.MessageInteractionWithSource,
				Data: &api.InteractionResponseData{
					Flags: api.EphemeralResponse,
					Embeds: &[]discord.Embed{
						{
							Title: "Vote successful!",
							Color: utils.DiscordGreen,
						},
					},
				},
			}
			if err := core.State.RespondInteraction(e.ID, e.Token, res); err != nil {
				logger.Error(err)
			}

			Voted[e.Member.User.ID] = OptionsMap[data.CustomID]
			Scores[data.CustomID]++
			Update()
		},
	}
	MapComponents["second_option"] = Component{

		Run: func(e *gateway.InteractionCreateEvent, data *discord.ComponentInteractionData) {
			if Voted[e.Member.User.ID] != "" {
				logger.Warn(e.Member.User.Username, "has already voted!")
				res := api.InteractionResponse{
					Type: api.MessageInteractionWithSource,
					Data: &api.InteractionResponseData{
						Flags: api.EphemeralResponse,
						Embeds: &[]discord.Embed{
							{
								Title:       "You've already voted!!",
								Description: fmt.Sprintf("You voted for `%v`", Voted[e.Member.User.ID]),
								Color:       utils.DiscordRed,
							},
						},
					},
				}
				if err := core.State.RespondInteraction(e.ID, e.Token, res); err != nil {
					logger.Error(err)
				}
				return
			}
			res := api.InteractionResponse{
				Type: api.MessageInteractionWithSource,
				Data: &api.InteractionResponseData{
					Flags: api.EphemeralResponse,
					Embeds: &[]discord.Embed{
						{
							Title: "Vote successful!",
							Color: utils.DiscordGreen,
						},
					},
				},
			}
			if err := core.State.RespondInteraction(e.ID, e.Token, res); err != nil {
				logger.Error(err)
			}

			Voted[e.Member.User.ID] = OptionsMap[data.CustomID]
			Scores[data.CustomID]++

			Update()
		},
	}
	MapComponents["third_option"] = Component{

		Run: func(e *gateway.InteractionCreateEvent, data *discord.ComponentInteractionData) {
			if Voted[e.Member.User.ID] != "" {
				logger.Warn(e.Member.User.Username, "has already voted!")
				res := api.InteractionResponse{
					Type: api.MessageInteractionWithSource,
					Data: &api.InteractionResponseData{
						Flags: api.EphemeralResponse,
						Embeds: &[]discord.Embed{
							{
								Title:       "You've already voted!!",
								Description: fmt.Sprintf("You voted for `%v`", Voted[e.Member.User.ID]),
								Color:       utils.DiscordRed,
							},
						},
					},
				}
				if err := core.State.RespondInteraction(e.ID, e.Token, res); err != nil {
					logger.Error(err)
				}
				return
			}
			res := api.InteractionResponse{
				Type: api.MessageInteractionWithSource,
				Data: &api.InteractionResponseData{
					Flags: api.EphemeralResponse,
					Embeds: &[]discord.Embed{
						{
							Title: "Vote successful!",
							Color: utils.DiscordGreen,
						},
					},
				},
			}
			if err := core.State.RespondInteraction(e.ID, e.Token, res); err != nil {
				logger.Error(err)
			}

			Voted[e.Member.User.ID] = OptionsMap[data.CustomID]
			Scores[data.CustomID]++

			Update()
		},
	}
	MapComponents["fourth_option"] = Component{

		Run: func(e *gateway.InteractionCreateEvent, data *discord.ComponentInteractionData) {
			if Voted[e.Member.User.ID] != "" {
				logger.Warn(e.Member.User.Username, "has already voted!")
				res := api.InteractionResponse{
					Type: api.MessageInteractionWithSource,
					Data: &api.InteractionResponseData{
						Flags: api.EphemeralResponse,
						Embeds: &[]discord.Embed{
							{
								Title:       "You've already voted!!",
								Description: fmt.Sprintf("You voted for `%v`", Voted[e.Member.User.ID]),
								Color:       utils.DiscordRed,
							},
						},
					},
				}
				if err := core.State.RespondInteraction(e.ID, e.Token, res); err != nil {
					logger.Error(err)
				}
				return
			}
			res := api.InteractionResponse{
				Type: api.MessageInteractionWithSource,
				Data: &api.InteractionResponseData{
					Flags: api.EphemeralResponse,
					Embeds: &[]discord.Embed{
						{
							Title: "Vote successful!",
							Color: utils.DiscordGreen,
						},
					},
				},
			}
			if err := core.State.RespondInteraction(e.ID, e.Token, res); err != nil {
				logger.Error(err)
			}

			Voted[e.Member.User.ID] = OptionsMap[data.CustomID]
			Scores[data.CustomID]++

			Update()
		},
	}
}
