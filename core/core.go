package core

import (
	"encoding/json"
	"fmt"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/kingultron99/tdcbot/logger"
	"github.com/lukasl-dev/waterlink"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"
)

type configStruct struct {
	Token     string `json:"token"`
	Owner     string `json:"owner"`
	APPID     string `json:"appid"`
	GUILDID   string `json:"guildid"`
	WolframID string `json:"wolframid"`
	Version   string `json:"version"`
	OwnerID   string `json:"ownerid"`
}

var (
	Config  *configStruct
	State   *state.State
	TimeNow time.Time
	Logg    *zap.Logger
	Conn    waterlink.Connection
	Update  *gateway.VoiceServerUpdateEvent
)

func init() {
	TimeNow = time.Now()
}

// Initialise sets up the logger and calls for "setupCloseHandler" and "initConfig"
func Initialise() {
	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	cmd.Run()

	writerSync := logger.GetLogWriter()
	encoder := logger.GetEncoder()

	core := zapcore.NewCore(encoder, writerSync, zapcore.DebugLevel)
	Logg = zap.New(core)

	defer Logg.Sync()

	zap.ReplaceGlobals(Logg)

	setupCloseHandler(Logg)
	initConfig()

	logger.Print(fmt.Sprintf(" ::::::::   ::::::::       ::::::::: ::::::::    ::::::::\n:+:    :+: :+:    :+:         :+:    :+:   :+:  :+:    :+:\n+:+    +:+ +:+    +:+   (:o   +:+    +:+    +:+ +:+\n+#+        +#+    +#+ +#+#+#+ +#+    +#+    +#+ +#+\n+#+   #+#+ +#+    +#+         +#+    +#+    +#+ +#+\n#+#    #+# #+#    #+#         #+#    #+#   #+#  #+#    #+#\n ########   ########          ###    ########    ########\n                                              %v\n", Config.Version))
	logger.Debug("Initialised Logger!")
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
		logger.Info("Goodbye!")
		os.Exit(0)
	}()
}
