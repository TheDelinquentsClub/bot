package pollStuff

import (
	"fmt"
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/kingultron99/tdcbot/commands"
	"github.com/kingultron99/tdcbot/core"
	"github.com/kingultron99/tdcbot/logger"
	"github.com/kingultron99/tdcbot/utils"
	"strings"
)

var Voted = make(map[discord.UserID]string)

var (
	Item1 int
	Item2 int
	Item3 int
	Item4 int
)

func GetGraph(item string) string {
	str := strings.Builder{}
	switch item {
	case "first_option":
		for i := 0; i < Item1; i++ {
			str.WriteString("|")
		}
	case "second_option":
		for i := 0; i < Item1; i++ {
			str.WriteString("|")
		}
	case "third_option":
		for i := 0; i < Item1; i++ {
			str.WriteString("|")
		}
	case "fourth_option":
		for i := 0; i < Item1; i++ {
			str.WriteString("|")
		}
	default:
		str.WriteString("")
	}
	return str.String()
}

func Update() {
	res := api.EditInteractionResponseData{
		Embeds: &[]discord.Embed{
			{
				Title:  commands.Data,
				Fields: commands.Fields,
			},
		},
		Components: &[]discord.Component{
			&discord.ActionRowComponent{
				Components: commands.Components,
			},
		},
	}
	if _, err := core.State.EditInteractionResponse(discord.AppID(utils.MustSnowflakeEnv(core.Config.APPID)), fmt.Sprint(commands.InteractionID), res); err != nil {
		logger.Error(err)
	}
}
