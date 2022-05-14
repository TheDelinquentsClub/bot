package main

import (
	"context"
	"fmt"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/session/shard"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/kingultron99/tdcbot/commands"
	"github.com/kingultron99/tdcbot/core"
	"github.com/kingultron99/tdcbot/logger"
	"github.com/kingultron99/tdcbot/utils"
	"github.com/kingultron99/tdcbot/websockets"
)

func main() {
	core.Initialise()
	newShard := state.NewShardFunc(func(m *shard.Manager, s *state.State) {
		s.AddIntents(gateway.IntentGuilds)
		s.AddIntents(gateway.IntentGuildMessages)
		s.AddIntents(gateway.IntentDirectMessages)
		core.State = s
	})
	m, err := shard.NewManager(fmt.Sprint("Bot "+core.Config.Token), newShard)
	if err != nil {
		logger.Fatal("Failed to make shard manager!", err)
	}

	if err := m.Open(context.Background()); err != nil {
		logger.Fatal("Failed to connect shards!", err)
	}
	defer m.Close()
	var shardNum int

	m.ForEach(func(s shard.Shard) {
		state := s.(*state.State)
		u, err := state.Me()
		if err != nil {
			logger.Fatal("Failed to get myself!", err)
		}
		logger.Info(fmt.Sprintf("Shard %d/%d started as %s", shardNum, m.NumShards()-1, u.Tag()))
		shardNum++
	})

	// Message event handler for sending messages to the minecraft server
	core.State.AddHandler(func(c *gateway.MessageCreateEvent) {
		if core.IsServerConnected {
			if c.Message.ChannelID.String() == core.Config.BridgeChannelID {
				if c.Author.Bot == false {
					core.ServerConn.Emit("discordmessage", c.Message.Content, c.Member.User.Username)
				}
			}
		}

	})

	utils.MapIcons()

	commands.AddHandlers()
	commands.Register(discord.AppID(utils.MustSnowflakeEnv(core.Config.APPID)), discord.GuildID(utils.MustSnowflakeEnv(core.Config.GUILDID)))
	go websockets.InitServer()

	// Block forever
	select {}
}
