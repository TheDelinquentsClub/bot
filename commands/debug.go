package commands

import (
	"fmt"
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"github.com/kingultron99/tdcbot/Maps"
	"github.com/kingultron99/tdcbot/core"
	"github.com/kingultron99/tdcbot/logger"
	"github.com/kingultron99/tdcbot/structs"
	"github.com/kingultron99/tdcbot/utils"
	"runtime"
	"strings"
	"time"
)

func init() {
	Maps.MapCommands["stats"] = structs.Command{
		Name:        "stats",
		Description: "Returns the current statistics and host system information of GoTDC",
		Group:       "debug",
		OwnerOnly:   false,
		Usage:       "/stats",
		Run: func(e *gateway.InteractionCreateEvent, data *discord.CommandInteractionData) {
			var m runtime.MemStats
			runtime.ReadMemStats(&m)

			res := api.InteractionResponse{
				Type: api.MessageInteractionWithSource,
				Data: &api.InteractionResponseData{
					Embeds: &[]discord.Embed{
						{
							Title:       "GoTDC Stats",
							Description: "GoTDC is a bot made specifically for \"The Delinquents Club\" discord.",
							Color:       utils.DefaultColour,
							Fields: []discord.EmbedField{
								{
									Name:   "GoTDC Version",
									Value:  core.Config.Version,
									Inline: true,
								},
								{
									Name:   "GoLang Version",
									Value:  strings.Trim(runtime.Version(), "go"),
									Inline: true,
								},
								{
									Name:   "№ of GoRoutines",
									Value:  fmt.Sprintf("%v", runtime.NumGoroutine()),
									Inline: true,
								},
								{
									Name:   "Uptime",
									Value:  utils.GetDurationString(time.Since(core.TimeNow)),
									Inline: true,
								},
								{
									Name:   "OS",
									Value:  runtime.GOOS,
									Inline: true,
								},
								{
									Name:  "Memory Used",
									Value: fmt.Sprintf("using %v MB / %v MB\n%v MB garbage collected. next GC cycle at %v MB.\ncurrent number of GC Cycles: %v", utils.BToMb(m.Alloc), utils.BToMb(m.Sys), utils.BToMb(m.GCSys), utils.BToMb(m.NextGC), m.NumGC),
								},
							},
							Footer: &discord.EmbedFooter{
								Text: e.Member.User.Username,
								Icon: e.Member.User.AvatarURL(),
							},
							Timestamp: discord.NowTimestamp(),
						},
					},
				},
			}
			if err := core.State.RespondInteraction(e.ID, e.Token, res); err != nil {
				logger.Error(err)
			}
		},
	}
	Maps.MapCommands["gc"] = structs.Command{
		Name:        "gc",
		Description: "Triggers a garbage collection cycle",
		Usage:       "/gc",
		Group:       "debug",
		OwnerOnly:   true,
		Run: func(e *gateway.InteractionCreateEvent, data *discord.CommandInteractionData) {
			logger.Info(e.Member.User.Username, "triggered a GC cycle!")
			runtime.GC()

			res := api.InteractionResponse{
				Type: api.MessageInteractionWithSource,
				Data: &api.InteractionResponseData{
					Content: option.NewNullableString(":wastebasket: triggered a GC cycle!"),
				},
			}

			if err := core.State.RespondInteraction(e.ID, e.Token, res); err != nil {
				logger.Error(err)
			}
		},
	}
	Maps.MapCommands["kill"] = structs.Command{
		Name:        "kill",
		Description: "Kills the bots process.",
		OwnerOnly:   true,
		Usage:       "/kill",
		Group:       "debug",
		Run: func(e *gateway.InteractionCreateEvent, data *discord.CommandInteractionData) {
			res := api.InteractionResponse{
				Type: api.MessageInteractionWithSource,
				Data: &api.InteractionResponseData{
					Embeds: &[]discord.Embed{
						{
							Title:       "Bye Bye! :wave:",
							Description: "Killing the process....",
							Color:       utils.DiscordGreen,
							Footer: &discord.EmbedFooter{
								Text: fmt.Sprintf("Killed by %v#%v", e.Member.User.Username, e.Member.User.Discriminator),
								Icon: e.Member.User.AvatarURL(),
							},
							Timestamp: discord.NowTimestamp(),
						},
					},
				},
			}
			if err := core.State.RespondInteraction(e.ID, e.Token, res); err != nil {
				logger.Error(err)
			}

			logger.Fatal(fmt.Sprintf("User %v#%v killed the process", e.Member.User.Username, e.Member.User.Discriminator))
		},
	}
}
