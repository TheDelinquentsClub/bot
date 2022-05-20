package commands

import (
	"bytes"
	"fmt"
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/utils/sendpart"
	"github.com/kingultron99/tdcbot/core"
	"github.com/kingultron99/tdcbot/logger"
	"github.com/kingultron99/tdcbot/utils"
	"github.com/pbnjay/memory"
	"os"
	"runtime"
	"strings"
	"time"
)

func init() {
	MapCommands["stats"] = Command{
		Name:        "stats",
		Description: "Returns the current statistics and host system information of GoTDC",
		Group:       "debug",
		Restricted:  false,
		Usage:       "/stats",
		Run: func(e *gateway.InteractionCreateEvent, data *discord.CommandInteraction) {
			var m runtime.MemStats
			runtime.ReadMemStats(&m)

			res := api.InteractionResponse{
				Type: api.MessageInteractionWithSource,
				Data: &api.InteractionResponseData{
					Flags: api.EphemeralResponse,
					Embeds: &[]discord.Embed{
						{
							Title:       "GoTDC Stats",
							Description: "GoTDC is a bot made specifically for \"The Delinquents Club\" discord.",
							Timestamp:   discord.NowTimestamp(),
							Color:       utils.DefaultColour,
							Footer: &discord.EmbedFooter{
								Text: e.Member.User.Username,
								Icon: e.Member.User.AvatarURL(),
							},
							Fields: []discord.EmbedField{
								{
									Name:   "GoTDC Version",
									Value:  core.Config.Version,
									Inline: true,
								},
								{
									Name:   "Golang Version",
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
									Name:   "PID",
									Inline: true,
									Value:  fmt.Sprint(os.Getpid()),
								},
								{
									Name: "Memory Used",
									Value: fmt.Sprintf(
										"using %v MB / %v MB\n%v MB garbage collected. next GC cycle at %v MB.\ncurrent number of GC Cycles: %v",
										utils.BToMb(m.Alloc),
										utils.BToMb(m.Sys),
										utils.BToMb(m.GCSys),
										utils.BToMb(m.NextGC),
										m.NumGC),
									Inline: false,
								},
								{
									Name:   "System memory",
									Inline: true,
									Value:  fmt.Sprintf("%v / %v", memory.FreeMemory(), memory.TotalMemory()),
								},
								{
									Name:   "Cores",
									Inline: true,
									Value:  fmt.Sprint(runtime.NumCPU()),
								},
							},
						},
					},
				},
			}

			if err := core.State.RespondInteraction(e.ID, e.Token, res); err != nil {
				logger.Error(err)
			}
		},
	}
	MapCommands["gc"] = Command{
		Name:        "gc",
		Description: "Triggers a garbage collection cycle",
		Usage:       "/gc",
		Group:       "debug",
		Restricted:  true,
		Run: func(e *gateway.InteractionCreateEvent, data *discord.CommandInteraction) {
			logger.Info(e.Member.User.Username, "triggered a GC cycle!")
			runtime.GC()

			res := api.InteractionResponse{
				Type: api.MessageInteractionWithSource,
				Data: &api.InteractionResponseData{
					Flags: api.EphemeralResponse,
					Embeds: &[]discord.Embed{
						{
							Title:     ":wastebasket: triggered a GC cycle!",
							Footer:    &discord.EmbedFooter{Text: e.Member.User.Username, Icon: e.Member.User.AvatarURL()},
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
	MapCommands["kill"] = Command{
		Name:        "kill",
		Description: "Kills the bots process.",
		Restricted:  true,
		Usage:       "/kill",
		Group:       "debug",
		Run: func(e *gateway.InteractionCreateEvent, data *discord.CommandInteraction) {

			res := api.InteractionResponse{
				Type: api.MessageInteractionWithSource,
				Data: &api.InteractionResponseData{
					Flags: api.EphemeralResponse,
					Embeds: &[]discord.Embed{
						{
							Title:       "Bye bye :wave:",
							Description: "Killing the process...",
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

			logger.Info(fmt.Sprintf("User %v#%v triggered bot shutdown", e.Member.User.Username, e.Member.User.Discriminator))
			core.Logg.Sync()
			logger.Debug("Flushed log buffer")
			logger.Info("Goodbye!")
			os.Exit(0)

		},
	}
	MapCommands["logs"] = Command{
		Name:        "logs",
		Description: "sends the latest log as a file",
		Group:       "debug",
		Usage:       "/logs",
		Options:     nil,
		Restricted:  false,
		Run: func(e *gateway.InteractionCreateEvent, data *discord.CommandInteraction) {
			logger.Info(fmt.Sprintf("%v#%v requested the latest log file!", e.Member.User.Username, e.Member.User.Discriminator))

			logfile, err := os.ReadFile(logger.LogFile.Name())
			if err != nil {
				logger.Error(err)
			}
			reader := bytes.NewReader(logfile)

			res := api.InteractionResponse{
				Type: api.MessageInteractionWithSource,
				Data: &api.InteractionResponseData{
					Embeds: &[]discord.Embed{
						{
							Title: "Here's the latest log file!",
							Color: utils.DiscordGreen,
							Footer: &discord.EmbedFooter{
								Text: fmt.Sprintf("Requested by %v#%v", e.Member.User.Username, e.Member.User.Discriminator),
								Icon: e.Member.User.AvatarURL(),
							},
							Timestamp: discord.NowTimestamp(),
						},
					},
					Files: []sendpart.File{
						{
							Name:   "Log-Latest.log",
							Reader: reader,
						},
					},
				},
			}

			if err = core.State.RespondInteraction(e.ID, e.Token, res); err != nil {
				logger.Error(err)
			}
		},
	}
}
