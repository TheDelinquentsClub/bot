package commands

import (
	"encoding/json"
	"fmt"
	"github.com/Krognol/go-wolfram"
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/kingultron99/tdcbot/core"
	"github.com/kingultron99/tdcbot/logger"
	"github.com/kingultron99/tdcbot/utils"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"
)

func init() {
	MapCommands["question"] = Command{
		Name:        "question",
		Description: "Massive amounts of knowledge at your fingertips",
		Group:       "misc",
		Usage:       "/question <question>",
		OwnerOnly:   false,
		Options: []discord.CommandOption{
			{
				Type:        discord.CommandOptionType(3),
				Name:        "question",
				Description: "What would you like to know?",
				Required:    true,
			},
			{
				Type:        discord.CommandOptionType(3),
				Name:        "measurement_system",
				Description: "Metric (default) or Imperial",
				Choices: []discord.CommandOptionChoice{
					{
						Name:  "Imperial",
						Value: "Imperial",
					},
					{
						Name:  "Metric",
						Value: "Metric",
					},
				},
				Required: false,
			},
		},
		Run: func(e *gateway.InteractionCreateEvent, data *discord.CommandInteractionData) {
			ack := api.InteractionResponse{
				Type: api.DeferredMessageInteractionWithSource,
			}

			if err := core.State.RespondInteraction(e.ID, e.Token, ack); err != nil {
				logger.Error(err)
			}

			waAPI := &wolfram.Client{AppID: core.Config.WolframID}
			var (
				query       = fmt.Sprint(data.Options[0])
				measurement wolfram.Unit
			)
			if len(data.Options) > 1 {
				m := fmt.Sprint(data.Options[1])
				switch m {
				case "Imperial":
					measurement = wolfram.Imperial
				case "Metric":
					measurement = wolfram.Metric
				}
			} else {
				measurement = wolfram.Metric
			}

			answer, err := waAPI.GetShortAnswerQuery(query, measurement, 10)
			if err != nil {
				logger.Error(fmt.Sprintf("Failed to get answer for query.\nerr: %v", err))
			}

			if utf8.RuneCountInString(answer) > 1024 {
				logger.Warn("message over allowed limit")
				answer = answer[:1020] + "..."
				logger.Info(answer)
			}

			res := utils.NewEmbed().
				SetTitle(fmt.Sprintf("Question: %v", query)).
				AddField("Answer", false, answer).
				SetColor(utils.DiscordGreen).
				SetFooter(fmt.Sprintf("Quesition submitted by %v#%v", e.Member.User.Username, e.Member.User.Discriminator), e.Member.User.AvatarURL()).
				SetTimestamp().
				EditInteraction()

			if _, err := core.State.EditInteractionResponse(discord.AppID(utils.MustSnowflakeEnv(core.Config.APPID)), e.Token, res); err != nil {
				logger.Error(err)
			}

		},
	}
	MapCommands["randomfact"] = Command{
		Name:        "randomfact",
		Description: "Here, have a fact! But not just any old fact! Its a \"Random Fact\"! Just for you!",
		Group:       "misc",
		Usage:       "/randomFact",
		Options: []discord.CommandOption{
			{
				Type:        discord.CommandOptionType(5),
				Name:        "fact_of_the_day",
				Description: "Get the fact of the day!",
				Required:    false,
			},
		},
		OwnerOnly: false,
		Run: func(e *gateway.InteractionCreateEvent, data *discord.CommandInteractionData) {
			ack := api.InteractionResponse{
				Type: api.DeferredMessageInteractionWithSource,
			}

			if err := core.State.RespondInteraction(e.ID, e.Token, ack); err != nil {
				logger.Error(err)
			}

			type factStruct struct {
				Id        string `json:"id"`
				Text      string `json:"text"`
				Source    string `json:"source"`
				Language  string `json:"language"`
				Permalink string `json:"permalink"`
				SourceURL string `json:"source_url"`
			}

			var (
				fact    = new(factStruct)
				url     string
				boolean bool
			)

			if len(data.Options) != 0 {
				b, err := strconv.ParseBool(fmt.Sprint(data.Options[0]))
				if err != nil {
					logger.Error(err)
				}
				boolean = b
			}

			if len(data.Options) == 1 {
				switch boolean {
				case true:
					url = "https://uselessfacts.jsph.pl/today.json?language=en"
				case false:
					url = "https://uselessfacts.jsph.pl/random.json?language=en"
				}
			} else {
				url = "https://uselessfacts.jsph.pl/random.json?language=en"
			}

			resp, err := http.Get(url)
			if err != nil {
				logger.Error(err)
			}

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				logger.Error(err)
			}

			err = json.Unmarshal(body, &fact)

			if strings.Contains(fact.Text, "`") {
				fact.Text = strings.ReplaceAll(fact.Text, "`", "'")
			}

			res := utils.NewEmbed().
				SetTitle(fact.Text).
				SetColor(utils.DiscordBlue).
				SetFooter(fmt.Sprintf("Requested by %v#%v", e.Member.User.Username, e.Member.User.Discriminator), e.Member.User.AvatarURL()).
				AddURLButton("Source", fact.SourceURL).
				AddURLButton("Permalink", fact.Permalink).
				SetTimestamp().
				EditInteraction()

			if _, err := core.State.EditInteractionResponse(discord.AppID(utils.MustSnowflakeEnv(core.Config.APPID)), e.Token, res); err != nil {
				logger.Error(err)
			}

		},
	}
	MapCommands["logs"] = Command{
		Name:        "logs",
		Description: "sends the latest log as a file",
		Group:       "misc",
		Usage:       "/logs",
		Options:     nil,
		OwnerOnly:   false,
		Run: func(e *gateway.InteractionCreateEvent, data *discord.CommandInteractionData) {
			logger.Info(fmt.Sprintf("%v#%v requested the latest log file!", e.Member.User.Username, e.Member.User.Discriminator))

			logfile, err := os.ReadFile(logger.LogFile.Name())
			if err != nil {
				logger.Error(err)
			}

			res := utils.NewEmbed().
				SetTitle("Here's the latest log file!").
				SetColor(utils.DiscordGreen).
				SetFooter(fmt.Sprintf("Requested by %v#%v", e.Member.User.Username, e.Member.User.Discriminator), e.Member.User.AvatarURL()).
				AddFile(logger.LogFile.Name(), logfile).
				SetTimestamp().
				MakeResponse()

			if err = core.State.RespondInteraction(e.ID, e.Token, res); err != nil {
				logger.Error(err)
			}
		},
	}
}
