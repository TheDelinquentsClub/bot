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
	"github.com/kingultron99/tdcbot/structs"
	"github.com/kingultron99/tdcbot/utils"
	"io/ioutil"
	"net/http"
	"strings"
)

func init() {
	MapCommands["question"] = structs.Command{
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
				logger.Error.Println(err)
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
				logger.Error.Println(fmt.Sprintf("Failed to get answer for query.\nerr: %v", err))
			}

			resp := api.EditInteractionResponseData{
				Embeds: &[]discord.Embed{
					{
						Title: fmt.Sprintf("Question: %v", query),
						Fields: []discord.EmbedField{
							{
								Name:  "Answer:",
								Value: answer,
							},
						},
						Color:     utils.DiscordGreen,
						Timestamp: discord.NowTimestamp(),
						Footer: &discord.EmbedFooter{
							Text: fmt.Sprintf("Quesition submitted by %v#%v", e.Member.User.Username, e.Member.User.Discriminator),
							Icon: e.Member.User.AvatarURL(),
						},
					},
				},
			}

			if _, err := core.State.EditInteractionResponse(discord.AppID(utils.MustSnowflakeEnv(core.Config.APPID)), e.Token, resp); err != nil {
				logger.Error.Println(err)
			}

		},
	}
	MapCommands["randomfact"] = structs.Command{
		Name:        "randomfact",
		Description: "Here, have a fact! But not just any old fact! Its a \"Random Fact\"! Just for you!",
		Group:       "misc",
		Usage:       "/randomFact",
		Options: []discord.CommandOption{
			{
				Type:        discord.CommandOptionType(3),
				Name:        "fact_of_the_day",
				Description: "Get the fact of the day!",
				Required:    false,
				Choices: []discord.CommandOptionChoice{
					{
						Name:  "Yes",
						Value: "yes",
					},
					{
						Name:  "No",
						Value: "no",
					},
				},
			},
		},
		OwnerOnly: false,
		Run: func(e *gateway.InteractionCreateEvent, data *discord.CommandInteractionData) {
			ack := api.InteractionResponse{
				Type: api.DeferredMessageInteractionWithSource,
			}

			if err := core.State.RespondInteraction(e.ID, e.Token, ack); err != nil {
				logger.Error.Println(err)
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
				fact = new(factStruct)
				url  string
			)

			if len(data.Options) == 1 {
				switch fmt.Sprint(data.Options[0]) {
				case "no":
					url = "https://uselessfacts.jsph.pl/random.json?language=en"
				case "yes":
					url = "https://uselessfacts.jsph.pl/today.json?language=en"
				}
			} else {
				url = "https://uselessfacts.jsph.pl/random.json?language=en"
			}

			resp, err := http.Get(url)
			if err != nil {
				logger.Error.Println(err)
			}

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				logger.Error.Println(err)
			}

			err = json.Unmarshal(body, &fact)

			if strings.Contains(fact.Text, "`") {
				fact.Text = strings.ReplaceAll(fact.Text, "`", "'")
			}

			res := api.EditInteractionResponseData{
				Embeds: &[]discord.Embed{
					{
						Title:     fact.Text,
						Color:     utils.DiscordBlue,
						Timestamp: discord.NowTimestamp(),
						Footer: &discord.EmbedFooter{
							Text: fmt.Sprintf("Requested by %v#%v", e.Member.User.Username, e.Member.User.Discriminator),
							Icon: e.Member.User.AvatarURL(),
						},
					},
				},
				Components: &[]discord.Component{
					&discord.ActionRowComponent{
						Components: []discord.Component{
							&discord.ButtonComponent{
								Label: "Source",
								Style: discord.LinkButton,
								URL:   fact.SourceURL,
							},
							&discord.ButtonComponent{
								Label: "Permalink",
								Style: discord.LinkButton,
								URL:   fact.Permalink,
							},
						},
					},
				},
			}

			if _, err := core.State.EditInteractionResponse(discord.AppID(utils.MustSnowflakeEnv(core.Config.APPID)), e.Token, res); err != nil {
				logger.Error.Println(err)
			}

		},
	}
}
