package commands

import (
	"fmt"
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/kingultron99/tdcbot/core"
	"github.com/kingultron99/tdcbot/logger"
	"github.com/kingultron99/tdcbot/utils"
	"strings"
)

var Voted = make(map[discord.UserID]string)
var Scores = make(map[string]int)

func GetGraph(item string) string {

	str := strings.Builder{}

	if Scores[item] == 0 {
		str.WriteString("No votes...")
	} else {
		for i := 0; i < Scores[item]; i++ {
			str.WriteString("|")
		}
	}

	logger.Debug(str.String())
	return str.String()
}

func GenFields(name string, item string) discord.EmbedField {
	field := discord.EmbedField{
		Name:   name,
		Value:  fmt.Sprintf("%v", GetGraph(item)),
		Inline: false,
	}

	return field
}

func Update() {
	var fields []discord.EmbedField
	for _, option := range Data {
		item := strings.ReplaceAll(fmt.Sprint(option.Name), "\"", "")
		name := strings.Title(strings.ReplaceAll(fmt.Sprint(option.Value), "\"", ""))
		fields = append(fields, GenFields(name, item))
	}

	res := api.EditInteractionResponseData{
		Embeds: &[]discord.Embed{
			{
				Title:  Title,
				Fields: fields,
			},
		},
		Components: &[]discord.Component{
			&discord.ActionRowComponent{
				Components: Components,
			},
		},
	}
	if _, err := core.State.EditInteractionResponse(discord.AppID(utils.MustSnowflakeEnv(core.Config.APPID)), InteractionID, res); err != nil {
		logger.Error(err)
	}
}
