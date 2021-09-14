package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/gateway/shard"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/kingultron99/tdcbot/commands"
	"github.com/kingultron99/tdcbot/core"
	"github.com/kingultron99/tdcbot/logger"
	"log"
	"os"
)

type ConfigStruct struct {
	Prefix  string `json:"prefix"`
	Token   string `json:"token"`
	Owner   string `json:"owner"`
	APPID   string `json:"appid"`
	GUILDID string `json:"guildid"`
	Version string `json:"version"`
}

var Config *ConfigStruct

func main() {
	logger.Info.Println("----------------------------------------------------------------")

	data, err := os.Open("config.json")
	if err != nil {
		logger.Error.Println(fmt.Sprintf("Error loading config: %v", err))
	} else {
		logger.Info.Println("Successfully loaded config!")
	}

	err = json.NewDecoder(data).Decode(&Config)

	newShard := state.NewShardFunc(func(m *shard.Manager, s *state.State) {
		// Add the needed Gateway intents.
		s.AddIntents(gateway.IntentGuildMessages)
		s.AddIntents(gateway.IntentDirectMessages)

		core.State = s
	})

	m, err := shard.NewManager(fmt.Sprint("Bot ", Config.Token), newShard)
	if err != nil {
		logger.Error.Println(fmt.Sprintf("failed to create shard manager: %v", err))
	}

	if err := m.Open(context.Background()); err != nil {
		logger.Error.Println(fmt.Sprintf("failed to connect shards: %v", err))
	}
	defer m.Close()

	var shardNum int

	m.ForEach(func(s shard.Shard) {
		state := s.(*state.State)

		u, err := state.Me()
		if err != nil {
			logger.Error.Println(fmt.Sprintf("failed to get myself: %v", err))
		}

		logger.Info.Println(fmt.Sprintf("Shard %d/%d started as %s", shardNum, m.NumShards()-1, u.Tag()))

		shardNum++
	})

	commands.AddHandlers()
	commands.Register(discord.AppID(mustSnowflakeEnv(Config.APPID)), discord.GuildID(mustSnowflakeEnv(Config.GUILDID)))

	// Block forever.
	select {}
}

func mustSnowflakeEnv(env string) discord.Snowflake {
	s, err := discord.ParseSnowflake(env)
	if err != nil {
		log.Fatalf("Invalid snowflake for $%s: %v", env, err)
	}
	return s
}
