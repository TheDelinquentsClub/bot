package componentInteractions

import (
	"fmt"
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/kingultron99/tdcbot/Maps"
	"github.com/kingultron99/tdcbot/componentInteractions/pollStuff"
	"github.com/kingultron99/tdcbot/core"
	"github.com/kingultron99/tdcbot/logger"
	"github.com/kingultron99/tdcbot/structs"
	"strings"
)

func init() {

	Maps.MapComponents["help_select_category"] = structs.Component{
		Run: func(e *gateway.InteractionCreateEvent, data *discord.ComponentInteractionData) {
			var fields, embedColour = Maps.GetCommands(data.Values[0])

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
	Maps.MapComponents["first_option"] = structs.Component{

		Run: func(e *gateway.InteractionCreateEvent, data *discord.ComponentInteractionData) {
			if pollStuff.Voted[e.Member.User.ID] != "" {
				logger.Warn(e.Member.User.Username, "has already voted!")
				res := api.InteractionResponse{
					Type: api.MessageInteractionWithSource,
					Data: &api.InteractionResponseData{
						Flags: api.EphemeralResponse,
						Embeds: &[]discord.Embed{
							{
								Title:       "You've already voted!!",
								Description: fmt.Sprintf("You voted for %v", pollStuff.Voted[e.Member.User.ID]),
							},
						},
					},
				}
				if err := core.State.RespondInteraction(e.ID, e.Token, res); err != nil {
					logger.Error(err)
				}
				return
			}

			pollStuff.Voted[e.Member.User.ID] = strings.ReplaceAll(fmt.Sprint(data.Values), "\"", "")
			pollStuff.Item1++

			logger.Info(pollStuff.Voted[e.Member.User.ID])

			pollStuff.Update()
		},
	}
	Maps.MapComponents["second_option"] = structs.Component{

		Run: func(e *gateway.InteractionCreateEvent, data *discord.ComponentInteractionData) {
			if pollStuff.Voted[e.Member.User.ID] != "" {
				logger.Warn(e.Member.User.Username, "has already voted!")
				res := api.InteractionResponse{
					Type: api.MessageInteractionWithSource,
					Data: &api.InteractionResponseData{
						Flags: api.EphemeralResponse,
						Embeds: &[]discord.Embed{
							{
								Title:       "You've already voted!!",
								Description: fmt.Sprintf("You voted for %v", pollStuff.Voted[e.Member.User.ID]),
							},
						},
					},
				}
				if err := core.State.RespondInteraction(e.ID, e.Token, res); err != nil {
					logger.Error(err)
				}
				return
			}

			pollStuff.Voted[e.Member.User.ID] = strings.ReplaceAll(fmt.Sprint(data.Values), "\"", "")
			pollStuff.Item2++

			logger.Info(pollStuff.Voted[e.Member.User.ID])

			pollStuff.Update()
		},
	}
	Maps.MapComponents["third_option"] = structs.Component{

		Run: func(e *gateway.InteractionCreateEvent, data *discord.ComponentInteractionData) {
			if pollStuff.Voted[e.Member.User.ID] != "" {
				logger.Warn(e.Member.User.Username, "has already voted!")
				res := api.InteractionResponse{
					Type: api.MessageInteractionWithSource,
					Data: &api.InteractionResponseData{
						Flags: api.EphemeralResponse,
						Embeds: &[]discord.Embed{
							{
								Title:       "You've already voted!!",
								Description: fmt.Sprintf("You voted for %v", pollStuff.Voted[e.Member.User.ID]),
							},
						},
					},
				}
				if err := core.State.RespondInteraction(e.ID, e.Token, res); err != nil {
					logger.Error(err)
				}
				return
			}

			pollStuff.Voted[e.Member.User.ID] = strings.ReplaceAll(fmt.Sprint(data.Values), "\"", "")
			pollStuff.Item3++

			logger.Info(pollStuff.Voted[e.Member.User.ID])

			pollStuff.Update()
		},
	}
	Maps.MapComponents["fourth_option"] = structs.Component{

		Run: func(e *gateway.InteractionCreateEvent, data *discord.ComponentInteractionData) {
			if pollStuff.Voted[e.Member.User.ID] != "" {
				logger.Warn(e.Member.User.Username, "has already voted!")
				res := api.InteractionResponse{
					Type: api.MessageInteractionWithSource,
					Data: &api.InteractionResponseData{
						Flags: api.EphemeralResponse,
						Embeds: &[]discord.Embed{
							{
								Title:       "You've already voted!!",
								Description: fmt.Sprintf("You voted for %v", pollStuff.Voted[e.Member.User.ID]),
							},
						},
					},
				}
				if err := core.State.RespondInteraction(e.ID, e.Token, res); err != nil {
					logger.Error(err)
				}
				return
			}

			pollStuff.Voted[e.Member.User.ID] = strings.ReplaceAll(fmt.Sprint(data.Values), "\"", "")
			pollStuff.Item4++

			logger.Info(pollStuff.Voted[e.Member.User.ID])

			pollStuff.Update()
		},
	}
}
