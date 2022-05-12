package core

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/diamondburned/arikawa/v3/state"
	socketio "github.com/googollee/go-socket.io"
	"github.com/kingultron99/tdcbot/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type configStruct struct {
	Token           string `json:"token"`
	APPID           string `json:"appid"`
	GUILDID         string `json:"guildid"`
	WolframID       string `json:"wolframid"`
	Version         string `json:"version"`
	CreatorID       string `json:"creatorid"`
	OwnerRole       string `json:"ownerRole"`
	BotBreakerRole  string `json:"botBreakerRole"`
	Webhook         string `json:"webhook"`
	BridgeChannelID string `json:"bridgeChannelId"`
}

var (
	Config            *configStruct
	State             *state.State
	TimeNow           time.Time
	Logg              *zap.Logger
	WSServer          *socketio.Server
	ServerConn        socketio.Conn
	ItemIcons         []string
	IsServerConnected = false
	clear             map[string]func() //create a map for storing clear funcs
)

func init() {
	TimeNow = time.Now()

	clear = make(map[string]func()) //Initialize it
	clear["linux"] = func() {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

// Initialise sets up the logger and calls for "setupCloseHandler" and "initConfig"
func Initialise() {
	value, ok := clear[runtime.GOOS] //runtime.GOOS -> linux, windows, darwin etc.
	if ok {                          //if we defined a clear func for that platform:
		value() //we execute it
	} else { //unsupported platform
		panic("Your platform is unsupported! I can't clear terminal screen :(")
	}

	writerSync := logger.GetLogWriter()
	encoder := logger.GetEncoder()

	core := zapcore.NewCore(encoder, writerSync, zapcore.DebugLevel)
	Logg = zap.New(core)

	defer Logg.Sync()

	zap.ReplaceGlobals(Logg)

	logger.Info("Initialised Logger!")

	setupCloseHandler(Logg)
	initConfig()

}

// initConfig Initialises the bots config
func initConfig() {

	data, err := os.Open("config.json")

	if os.IsNotExist(err) {
		logger.Error("No config was found!")
		byte, err := json.Marshal(configStruct{
			APPID:           "",
			BridgeChannelID: "",
			BotBreakerRole:  "",
			GUILDID:         "",
			CreatorID:       "",
			OwnerRole:       "",
			Token:           "",
			Version:         "",
			Webhook:         "",
			WolframID:       "",
		})
		if err != nil {
			logger.Error("Failed to marshal Config")
		}
		os.WriteFile("./config.json", byte, os.ModePerm)
		logger.Info("Generated Config.json template in root project root directory")
		logger.Info("Please fill Config with required fields.")
		os.Exit(1)
	} else if err != nil {
		logger.Error(fmt.Sprintf("Error loading config: %v", err))
		os.Exit(1)
	} else {
		logger.Print(fmt.Sprintf(" ::::::::   ::::::::       ::::::::: ::::::::    ::::::::\n:+:    :+: :+:    :+:         :+:    :+:   :+:  :+:    :+:\n+:+    +:+ +:+    +:+   (:o   +:+    +:+    +:+ +:+\n+#+        +#+    +#+ +#+#+#+ +#+    +#+    +#+ +#+\n+#+   #+#+ +#+    +#+         +#+    +#+    +#+ +#+\n#+#    #+# #+#    #+#         #+#    #+#   #+#  #+#    #+#\n ########   ########          ###    ########    ########\n                                              v%v\n", Config.Version))
		logger.Info("Successfully loaded config!")
	}

	_ = json.NewDecoder(data).Decode(&Config)
}

func setupCloseHandler(logg *zap.Logger) {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		WSServer.BroadcastToNamespace("/", "shutdown")
		logger.Info("Beginning shutdown process")
		logg.Sync()
		logger.Debug("Flushed log buffer")
		WSServer.Close()
		logger.Info("Closed WS server")
		logger.Info("Goodbye!")
		os.Exit(0)
	}()
}
