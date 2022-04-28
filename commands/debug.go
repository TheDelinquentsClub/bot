package commands

import (
	"fmt"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/kingultron99/tdcbot/core"
	"github.com/kingultron99/tdcbot/logger"
	"github.com/kingultron99/tdcbot/utils"
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
		OwnerOnly:   false,
		Usage:       "/stats",
		Run: func(e *gateway.InteractionCreateEvent, data *discord.CommandInteractionData) {
			var m runtime.MemStats
			runtime.ReadMemStats(&m)
			res := utils.NewEmbed().
				SetTitle("GoTDC Stats").
				SetDescription("GoTDC is a bot made specifically for \"The Delinquents Club\" discord.").
				SetColor(discord.DefaultEmbedColor).
				AddField("GoTDC Version", true, core.Config.Version).
				AddField("GoLang Version", true, strings.Trim(runtime.Version(), "go")).
				AddField("№ of GoRoutines", true, fmt.Sprintf("%v", runtime.NumGoroutine())).
				AddField("Uptime", true, utils.GetDurationString(time.Since(core.TimeNow))).
				AddField("OS", true, runtime.GOOS).
				AddField("Memory Used", false, fmt.Sprintf("using %v MB / %v MB\n%v MB garbage collected. next GC cycle at %v MB.\ncurrent number of GC Cycles: %v", utils.BToMb(m.Alloc), utils.BToMb(m.Sys), utils.BToMb(m.GCSys), utils.BToMb(m.NextGC), m.NumGC)).
				SetFooter(e.Member.User.Username, e.Member.User.AvatarURL()).
				MakeResponse()

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
		OwnerOnly:   true,
		Run: func(e *gateway.InteractionCreateEvent, data *discord.CommandInteractionData) {
			logger.Info(e.Member.User.Username, "triggered a GC cycle!")
			runtime.GC()

			res := utils.NewEmbed().
				SetTitle(":wastebasket: triggered a GC cycle!").SetFooter(e.Member.User.Username, e.Member.User.AvatarURL()).MakeResponse()

			if err := core.State.RespondInteraction(e.ID, e.Token, res); err != nil {
				logger.Error(err)
			}
		},
	}
	MapCommands["kill"] = Command{
		Name:        "kill",
		Description: "Kills the bots process.",
		OwnerOnly:   true,
		Usage:       "/kill",
		Group:       "debug",
		Run: func(e *gateway.InteractionCreateEvent, data *discord.CommandInteractionData) {

			res := utils.NewEmbed().
				SetTitle("Bye Bye :wave:").
				SetDescription("Killing the process...").
				SetColor(utils.DiscordGreen).
				SetFooter(fmt.Sprintf("Killed by %v#%v", e.Member.User.Username, e.Member.User.Discriminator), e.Member.User.AvatarURL()).
				MakeResponse()

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
}
