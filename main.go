package main

import (
	"context"
	"fmt"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/gateway/shard"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/kingultron99/tdcbot/commands"
	"github.com/kingultron99/tdcbot/core"
	"github.com/kingultron99/tdcbot/logger"
	"github.com/kingultron99/tdcbot/utils"
)

func main() {
	logger.Info.Println("----------------------------------------------------------------")

	core.InitConfig()

	newShard := state.NewShardFunc(func(m *shard.Manager, s *state.State) {
		// Add the needed Gateway intents.
		s.AddIntents(gateway.IntentGuildMessages)
		s.AddIntents(gateway.IntentDirectMessages)

		core.State = s
	})

	m, err := shard.NewManager(fmt.Sprint("Bot ", core.Config.Token), newShard)
	if err != nil {
		logger.Error.Println(fmt.Sprintf("failed to create shard manager: %v", err))
	}

	if err := m.Open(context.Background()); err != nil {
		logger.Error.Println(fmt.Sprintf("failed to connect shards: %v", err))
	}
	defer func(m *shard.Manager) {
		err := m.Close()
		if err != nil {
		}
	}(m)

	var shardNum int

	m.ForEach(func(s shard.Shard) {
		botState := s.(*state.State)

		u, err := botState.Me()
		if err != nil {
			logger.Error.Println(fmt.Sprintf("failed to get myself: %v", err))
		}

		logger.Info.Println(fmt.Sprintf("Shard %d/%d started as %s", shardNum, m.NumShards()-1, u.Tag()))

		shardNum++
	})

	commands.AddHandlers()
	commands.Register(discord.AppID(utils.MustSnowflakeEnv(core.Config.APPID)), discord.GuildID(utils.MustSnowflakeEnv(core.Config.GUILDID)))

	// Block forever.
	select {}
}
