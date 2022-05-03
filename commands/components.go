package commands

import (
	"fmt"
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

			res := utils.NewEmbed().
				SetTitle(fmt.Sprint("Help! | ", data.Values[0])).
				SetDescription(fmt.Sprint("Here are all the available commands in the ", data.Values[0], " category")).
				AddFields(fields).
				SetColor(embedColour).
				SetFooter(fmt.Sprintf("Requested by %v#%v", e.Member.User.Username, e.Member.User.Discriminator), e.Member.User.AvatarURL()).
				UpdateResponse()

			if err := core.State.RespondInteraction(e.ID, e.Token, res); err != nil {
				logger.Error(err)
			}
		},
	}

}
