package commands

import (
	"encoding/json"
	"fmt"
	"github.com/Krognol/go-wolfram"
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"github.com/kingultron99/tdcbot/core"
	"github.com/kingultron99/tdcbot/logger"
	"github.com/kingultron99/tdcbot/utils"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"
)

var (
	Fields        []discord.EmbedField
	Components    []discord.Component
	Title         string
	Data          []discord.InteractionOption
	InteractionID string
	OptionsMap    map[string]string
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
				fact = new(factStruct)
				url  string
				bool bool
			)

			if len(data.Options) != 0 {
				b, err := strconv.ParseBool(fmt.Sprint(data.Options[0]))
				if err != nil {
					logger.Error(err)
				}
				bool = b
			}

			if len(data.Options) == 1 {
				switch bool {
				case true:
					url = "https://uselessfacts.jsph.pl/random.json?language=en"
				case false:
					url = "https://uselessfacts.jsph.pl/today.json?language=en"
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
				logger.Error(err)
			}

		},
	}
	MapCommands["pollattempt"] = Command{
		Name:        "pollattempt",
		Description: "This is the first attempt at a poll command. planning to add vote visualisation as a ASCII graph",
		Group:       "misc",
		Usage:       "/pollattempt <title> <duration> <option1> <option2> [option3] [option4] ",
		Options: []discord.CommandOption{
			{
				Type:        discord.CommandOptionType(3),
				Name:        "title",
				Description: "What is this poll about?",
				Required:    true,
			},
			{
				Type:        discord.CommandOptionType(4),
				Name:        "duration",
				Description: "How long (in seconds) should this poll last for?",
				Required:    true,
			},
			{
				Type:        discord.CommandOptionType(3),
				Name:        "first_option",
				Description: " ",
				Required:    true,
			},
			{
				Type:        discord.CommandOptionType(3),
				Name:        "second_option",
				Description: " ",
				Required:    true,
			},
			{
				Type:        discord.CommandOptionType(3),
				Name:        "third_option",
				Description: " ",
				Required:    false,
			},
			{
				Type:        discord.CommandOptionType(3),
				Name:        "fourth_option",
				Description: " ",
				Required:    false,
			},
		},
		OwnerOnly: false,
		Run: func(e *gateway.InteractionCreateEvent, data *discord.CommandInteractionData) {

			oldres := api.EditInteractionResponseData{
				Content:    option.NewNullableString("**Poll is no-longer valid!**"),
				Components: &[]discord.Component{},
			}
			if _, err := core.State.EditInteractionResponse(discord.AppID(utils.MustSnowflakeEnv(core.Config.APPID)), InteractionID, oldres); err != nil {
				logger.Error("Failed to invalidate old poll!")
			}

			for key, _ := range OptionsMap {
				delete(OptionsMap, key)
				logger.Debug("deleted", key)
			}
			for key, _ := range Voted {
				delete(Voted, key)
				logger.Debug("deleted", key)
			}
			for key, _ := range Scores {
				delete(Scores, key)
				logger.Debug("deleted", key)
			}
			Fields = []discord.EmbedField{}
			Components = []discord.Component{}
			Title = ""
			Data = []discord.InteractionOption{}
			InteractionID = ""

			OptionsMap = make(map[string]string)
			Data = data.Options[2:]
			for _, option := range data.Options[2:] {
				item := strings.ReplaceAll(fmt.Sprint(option.Name), "\"", "")
				name := strings.Title(strings.ReplaceAll(fmt.Sprint(option.Value), "\"", ""))
				OptionsMap[item] = name
				Fields = append(Fields, GenFields(name, item))

				Components = append(Components, utils.GenButtonComponents(option))
			}

			Title = fmt.Sprint(data.Options[0])
			InteractionID = e.Token

			res := api.InteractionResponse{
				Type: api.MessageInteractionWithSource,
				Data: &api.InteractionResponseData{
					Embeds: &[]discord.Embed{
						{
							Title:  Title,
							Fields: Fields,
						},
					},
					Components: &[]discord.Component{
						&discord.ActionRowComponent{
							Components: Components,
						},
					},
				},
			}
			if err := core.State.RespondInteraction(e.ID, e.Token, res); err != nil {
				logger.Error(err)
			}
		},
	}
}
