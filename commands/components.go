package commands

import (
	"fmt"
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/kingultron99/tdcbot/core"
	"github.com/kingultron99/tdcbot/logger"
	"github.com/kingultron99/tdcbot/structs"
)

func init() {
	MapComponents["help_select_category"] = structs.Component{
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
				logger.Error.Println(err)
			}
		},
	}
}
