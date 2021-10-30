package commands

import (
	"fmt"
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/kingultron99/tdcbot/core"
	"github.com/kingultron99/tdcbot/logger"
	"github.com/kingultron99/tdcbot/utils"
	"strings"
	time2 "time"
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
		Value:  fmt.Sprintf("`%v`", GetGraph(item)),
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
				Title:       Title,
				Description: fmt.Sprintf("This poll will end in <t:%v:R>", duration),
				Fields:      fields,
				Color:       utils.DiscordBlue,
			},
		},
		Components: &[]discord.Component{
			&discord.ActionRowComponent{
				Components: Components,
			},
		},
	}
	if _, err := core.State.EditInteractionResponse(discord.AppID(utils.MustSnowflakeEnv(core.Config.APPID)), Interactiontoken, res); err != nil {
		logger.Error(err)
	}
}

func startTimer(time int) {
	logger.Info(time)
	time2.AfterFunc(time2.Second*time2.Duration(time), func() {
		logger.Info("Poll timer has ended!")
		endPoll()
	})
}

func endPoll() {
	var (
		option       string
		bignum, temp int
	)

	for key, keyval := range Scores {
		if keyval > temp {
			temp = keyval
			bignum = temp
			option = key
		}
	}

	if bignum == 0 {
		end := api.EditInteractionResponseData{
			Embeds: &[]discord.Embed{
				{
					Title:     fmt.Sprintf("There was no winner..."),
					Timestamp: discord.NowTimestamp(),
					Color:     utils.DiscordRed,
				},
			},
			Components: &[]discord.Component{},
		}
		if _, err := core.State.EditInteractionResponse(discord.AppID(utils.MustSnowflakeEnv(core.Config.APPID)), Interactiontoken, end); err != nil {
			logger.Error("Failed to end poll: ", err)
		}
	} else {
		end := api.EditInteractionResponseData{
			Embeds: &[]discord.Embed{
				{
					Title:     fmt.Sprintf("%v has won with a total of %v votes!", OptionsMap[option], bignum),
					Timestamp: discord.NowTimestamp(),
					Color:     utils.DiscordGreen,
				},
			},
			Components: &[]discord.Component{},
		}
		if _, err := core.State.EditInteractionResponse(discord.AppID(utils.MustSnowflakeEnv(core.Config.APPID)), Interactiontoken, end); err != nil {
			logger.Error("Failed to end poll: ", err)
		}
	}
}
