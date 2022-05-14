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
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"
)

func init() {
	// TODO: record statistics to sqlite db
	MapCommands["question"] = Command{
		Name:        "question",
		Description: "Massive amounts of knowledge at your fingertips",
		Group:       "misc",
		Usage:       "/question <question>",
		Restricted:  false,
		Options: []discord.CommandOption{
			&discord.StringOption{
				OptionName:  "question",
				Description: "What would you like to know?",
				Required:    true,
			},
			&discord.StringOption{
				OptionName:  "measurement_system",
				Description: "Metric (default) or Imperial",
				Required:    false,
				Choices: []discord.StringChoice{
					{
						Name:  "Imperial",
						Value: "Imperial",
					},
					{
						Name:  "Metric",
						Value: "Metric",
					},
				},
			},
		},
		Run: func(e *gateway.InteractionCreateEvent, data *discord.CommandInteraction) {
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
				color       discord.Color
				res         api.EditInteractionResponseData
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
				res = api.EditInteractionResponseData{
					Embeds: &[]discord.Embed{
						{
							Title:       "Failed to get answer for query.",
							Description: fmt.Sprint(err),
							Timestamp:   discord.NowTimestamp(),
							Color:       utils.DiscordRed,
							Footer: &discord.EmbedFooter{
								Text: fmt.Sprintf("Quesition submitted by %v#%v", e.Member.User.Username, e.Member.User.Discriminator),
								Icon: e.Member.User.AvatarURL(),
							},
						},
					},
				}
			} else {
				if utf8.RuneCountInString(answer) > 1024 {
					logger.Warn("message over allowed limit")
					answer = answer[:1020] + "..."
					logger.Info(answer)
				}
				color = utils.DiscordGreen
				if answer == "No short answer available" {
					color = utils.DiscordRed
				}

				stmt, err := core.DB.Prepare("UPDATE stats SET questions = questions + 1")
				if err != nil {
					logger.Error("failed to update Questions stat")
				}

				stmt.Exec()
				stmt.Close()

				res = api.EditInteractionResponseData{
					Embeds: &[]discord.Embed{
						{
							Title: fmt.Sprintf("Question: %v", query),
							Footer: &discord.EmbedFooter{
								Text: fmt.Sprintf("Quesition submitted by %v#%v", e.Member.User.Username, e.Member.User.Discriminator),
								Icon: e.Member.User.AvatarURL(),
							},
							Color:     color,
							Timestamp: discord.NowTimestamp(),
							Fields: []discord.EmbedField{
								{
									Name:   "Answer",
									Value:  answer,
									Inline: false,
								},
							},
						},
					},
				}
			}

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
			&discord.BooleanOption{
				OptionName:  "fotd",
				Description: "Get the fact of the day!",
				Required:    true,
			},
		},
		Restricted: false,
		Run: func(e *gateway.InteractionCreateEvent, data *discord.CommandInteraction) {

			stmt, err := core.DB.Prepare("UPDATE stats SET facts = facts + 1")
			if err != nil {
				logger.Error("failed to update facts stat")
			}
			stmt.Exec()
			stmt.Close()

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

			b, err := strconv.ParseBool(fmt.Sprint(data.Options[0]))
			if err != nil {
				logger.Error(err)
			}
			boolean = b

			switch boolean {
			case true:
				url = "https://uselessfacts.jsph.pl/today.json?language=en"
			case false:
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

			res := api.EditInteractionResponseData{
				Embeds: &[]discord.Embed{
					{
						Description: "**" + fact.Text + "**",
						Color:       utils.DiscordBlue,
						Footer: &discord.EmbedFooter{
							Text: fmt.Sprintf("Requested by %v#%v", e.Member.User.Username, e.Member.User.Discriminator),
							Icon: e.Member.User.AvatarURL(),
						},
						Timestamp: discord.NowTimestamp(),
					},
				},
				Components: discord.ComponentsPtr(
					&discord.ActionRowComponent{
						&discord.ButtonComponent{
							Label: "Source",
							Style: discord.LinkButtonStyle(fact.SourceURL),
						},
						&discord.ButtonComponent{
							Label: "Permalink",
							Style: discord.LinkButtonStyle(fact.Permalink),
						},
					}),
			}

			if _, err := core.State.EditInteractionResponse(discord.AppID(utils.MustSnowflakeEnv(core.Config.APPID)), e.Token, res); err != nil {
				logger.Error(err)
			}

		},
	}
}
