package main

import (
	"bufio"
	"context"
	"github.com/TheDelinquentsClub/bot/util"
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/session/shard"
	"github.com/diamondburned/arikawa/v3/state"
	"go.uber.org/zap"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func init() {
	zap.ReplaceGlobals(zap.New(util.LoggerCore))
}

func main() {
	logger := zap.S()
	config, err := util.NewConfig("config.json")
	if err != nil {
		logger.Fatal("Failed to parse configuration file:", err.Error())
	}

	manager, err := shard.NewManager(config.Token, state.NewShardFunc(func(m *shard.Manager, s *state.State) {
		self, err := s.Me()
		if err != nil {
			m.ForEach(func(shard shard.Shard) {
				_ = shard.Close()
			})
			log.Fatalf("Failed to get own user for shard #%d: %s", m.NumShards(), err.Error())
		}

		s.AddHandler(func(msg *gateway.MessageCreateEvent) {
			// Simply checks only this for now, for testing.
			if msg.ChannelID == config.MinecraftChannelID {
				if msg.Author.Bot || msg.Author.Discriminator == "0000" {
					return
				}
				_, _ = s.SendMessageComplex(msg.ChannelID, api.SendMessageData{
					Embeds: []discord.Embed{
						{
							Author: &discord.EmbedAuthor{
								Name: msg.Author.Username,
								Icon: msg.Author.AvatarURL(),
							},
							Description: msg.Content,
							Timestamp:   discord.NowTimestamp(),
						},
					},
				})
			}
		})

		s.AddHandler(func(interaction *gateway.InteractionCreateEvent) {
			// Nothing to do, yet.
		})

		s.AddIntents(
			gateway.IntentGuilds |
				gateway.IntentGuildMessages |
				gateway.IntentDirectMessages |
				gateway.IntentGuildVoiceStates,
		)

		if err := s.Open(context.Background()); err != nil {
			m.ForEach(func(shard shard.Shard) {
				_ = shard.Close()
			})
			log.Fatalf("Failed to start Gateway connection for shard %d: %s", m.NumShards(), err.Error())
		}

		logger.Infof("Started shard #%d as %s", m.NumShards(), self.Tag())
	}))
	if err != nil {
		logger.Fatal("Failed to create a shard manager:", err.Error())
	}

	inputReader := bufio.NewReader(os.Stdin)
	killCh := make(chan os.Signal, 1)
	signal.Notify(killCh, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case <-killCh:
			logger.Info("Are you sure you want to terminate? (Y/N)")
			text := func() string {
				text, _ := inputReader.ReadString('\n')

				l := len(text) - 1
				if l == -1 || l == 0 {
					return ""
				}
				if text[l-1] == '\r' {
					l -= 1
				}

				return text[:l]
			}()

			switch text {
			case "Y", "y", "Yes", "yes":
				logger.Info("Shutting down process...")
				manager.ForEach(func(shard shard.Shard) {
					_ = shard.Close()
				})
				_ = util.LoggerCore.Sync()
				os.Exit(0)
			case "N", "n", "No", "no":
				logger.Info("Continuing...")
				continue
			default:
				logger.Error("I didn't understand your input, continuing.")
				continue
			}
		}
	}
}
