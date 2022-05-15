package websockets

import (
	"encoding/json"
	"fmt"
	socketio "github.com/googollee/go-socket.io"
	"github.com/kingultron99/tdcbot/core"
	"github.com/kingultron99/tdcbot/logger"
	"github.com/kingultron99/tdcbot/utils"
	"os"
	"runtime"
	"strings"
	"time"
)

func RegisterWebsiteHandlers() {
	core.WSServer.OnEvent("/", "websiteinstance", func(s socketio.Conn) {
		core.Websites = append(core.Websites, s.ID())
	})

	core.WSServer.OnEvent("/", "botinfo", func(s socketio.Conn) {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		type info struct {
			BotVersion    string `json:"botVersion"`
			GoVersion     string `json:"goVersion"`
			GoRoutines    string `json:"goroutines"`
			CurrentUpTime string `json:"currentUpTime"`
			OS            string `json:"os"`
			PID           string `json:"pid"`
			Memory        string `json:"memory"`
		}
		var createRes = info{
			BotVersion:    core.Config.Version,
			GoVersion:     strings.Trim(runtime.Version(), "go"),
			GoRoutines:    fmt.Sprintf("%v", runtime.NumGoroutine()),
			CurrentUpTime: utils.GetDurationString(time.Since(core.TimeNow)),
			OS:            runtime.GOOS,
			PID:           fmt.Sprint(os.Getpid()),
			Memory: fmt.Sprintf("using %v MB / %v MB\n%v MB garbage collected. next GC cycle at %v MB.\ncurrent number of GC Cycles: %v",
				utils.BToMb(m.Alloc),
				utils.BToMb(m.Sys),
				utils.BToMb(m.GCSys),
				utils.BToMb(m.NextGC),
				m.NumGC),
		}
		res, err := json.Marshal(createRes)
		if err != nil {
			logger.Error("Failed to marshal bot info json: ", err)
		}
		s.Emit("botinfo", res)
	})

	core.WSServer.OnEvent("/", "stats", func(s socketio.Conn) {
		type info struct {
			Facts     string `json:"facts"`
			Questions string `json:"questions"`
		}
		var createRes info
		query, _ := core.DB.Query("SELECT facts, questions from stats")

		for query.Next() {
			query.Scan(&createRes.Facts, &createRes.Questions)
		}
		query.Close()

		res, err := json.Marshal(createRes)
		if err != nil {
			logger.Error("Failed to marshal stats json: ", err)
		}
		s.Emit("stats", res)
	})

}
