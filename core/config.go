package core

import (
	"encoding/json"
	"fmt"
	"github.com/kingultron99/tdcbot/logger"
	"os"
)

type configStruct struct {
	Token   string `json:"token"`
	Owner   string `json:"owner"`
	APPID   string `json:"appid"`
	GUILDID string `json:"guildid"`
	Version string `json:"version"`
	OwnerID string `json:"ownerid"`
}

var Config *configStruct

func InitConfig() {
	data, err := os.Open("config.json")
	if err != nil {
		logger.Error.Println(fmt.Sprintf("Error loading config: %v", err))
	} else {
		logger.Info.Println("Successfully loaded config!")
	}

	err = json.NewDecoder(data).Decode(&Config)
}
