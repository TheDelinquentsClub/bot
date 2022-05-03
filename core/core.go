package core

import (
	"encoding/json"
	"fmt"
	"github.com/diamondburned/arikawa/v3/state"
	socketio "github.com/googollee/go-socket.io"
	"github.com/kingultron99/tdcbot/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"
)

type configStruct struct {
	Token           string `json:"token"`
	Owner           string `json:"owner"`
	APPID           string `json:"appid"`
	GUILDID         string `json:"guildid"`
	WolframID       string `json:"wolframid"`
	Version         string `json:"version"`
	OwnerID         string `json:"ownerid"`
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
	IsServerConnected = false
)

func init() {
	TimeNow = time.Now()
}

// Initialise sets up the logger and calls for "setupCloseHandler" and "initConfig"
func Initialise() {
	cmd := exec.Command("cmd", "/c", "cls", "clear")
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		log.Fatal("Failed to clear terminal!")
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

	logger.Print(fmt.Sprintf(" ::::::::   ::::::::       ::::::::: ::::::::    ::::::::\n:+:    :+: :+:    :+:         :+:    :+:   :+:  :+:    :+:\n+:+    +:+ +:+    +:+   (:o   +:+    +:+    +:+ +:+\n+#+        +#+    +#+ +#+#+#+ +#+    +#+    +#+ +#+\n+#+   #+#+ +#+    +#+         +#+    +#+    +#+ +#+\n#+#    #+# #+#    #+#         #+#    #+#   #+#  #+#    #+#\n ########   ########          ###    ########    ########\n                                              %v\n", Config.Version))
}

// initConfig Initialises the bots config
func initConfig() {

	data, err := os.Open("config.json")
	if err != nil {
		logger.Error(fmt.Sprintf("Error loading config: %v", err))
	} else {
		logger.Info("Successfully loaded config!")
	}

	err = json.NewDecoder(data).Decode(&Config)
}

func setupCloseHandler(logg *zap.Logger) {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		logger.Info("Beginning shutdown process")
		logg.Sync()
		logger.Debug("Flushed log buffer")
		//WSServer.Close()
		//logger.Info("Closed WS server")
		logger.Info("Goodbye!")
		os.Exit(0)
	}()
}
