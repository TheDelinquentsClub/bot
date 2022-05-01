package commands

import (
	"encoding/json"
	"fmt"
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	socketio "github.com/googollee/go-socket.io"
	"github.com/kingultron99/tdcbot/core"
	"github.com/kingultron99/tdcbot/logger"
	"github.com/kingultron99/tdcbot/utils"
	"strings"
	"time"
)

func init() {
	MapCommands["minecraft"] = Command{
		Name:        "minecraft",
		Description: "Gets information regarding the game or a player from the public minecraft APIs",
		Group:       "minecraft",
		Usage:       "/minecraft <option>",
		Options: []discord.CommandOption{
			{
				Type:        1,
				Name:        "user",
				Description: "Gets information about a user, using either a valid username or UUID",
				Options: []discord.CommandOption{
					{
						Type:        3,
						Name:        "username",
						Description: "Get user information from a username",
					},
					{
						Type:        3,
						Name:        "uuid",
						Description: "Get user information from a UUID",
					},
				},
			},
		},
		OwnerOnly: false,
		Exclude:   false,
		Run: func(e *gateway.InteractionCreateEvent, data *discord.CommandInteractionData) {
			ack := api.InteractionResponse{
				Type: api.DeferredMessageInteractionWithSource,
			}

			if err := core.State.RespondInteraction(e.ID, e.Token, ack); err != nil {
				logger.Error(err)
			}

			switch data.Options[0].Name {
			case "user":
				switch data.Options[0].Options[0].Name {
				case "username":
					res := utils.NewEmbed().SetColor(utils.DiscordBlue).
						SetAuthor(strings.ReplaceAll(data.Options[0].Options[0].Value.String(), "\"", ""), fmt.Sprintf("https://crafatar.com/avatars/%v?size=100", utils.GetUUID(data.Options[0].Options[0].Value.String()))).
						SetTitle("Player Information").
						SetDescription("Showing infor for player by username").
						SetImage(fmt.Sprintf("https://crafatar.com/renders/body/%v", utils.GetUUID(data.Options[0].Options[0].Value.String()))).
						AddField("Username", true, strings.ReplaceAll(data.Options[0].Options[0].Value.String(), "\"", "")).
						AddField("UUID", true, utils.GetUUID(data.Options[0].Options[0].Value.String())).
						AddField("Names", false, utils.GetNamesFromUsername(data.Options[0].Options[0].Value.String())).
						AddURLButton("Cape", fmt.Sprintf("https://crafatar.com/capes/%v", utils.GetUUID(data.Options[0].Options[0].Value.String()))).
						AddURLButton("Skin", fmt.Sprintf("https://crafatar.com/skins/%v", utils.GetUUID(data.Options[0].Options[0].Value.String()))).
						SetTimestamp().
						EditInteraction()

					if _, err := core.State.EditInteractionResponse(discord.AppID(utils.MustSnowflakeEnv(core.Config.APPID)), e.Token, res); err != nil {
						logger.Error(err)
					}
				case "uuid":
					res := utils.NewEmbed().SetColor(utils.DiscordBlue).
						SetAuthor(strings.ReplaceAll(data.Options[0].Options[0].Value.String(), "\"", ""), fmt.Sprintf("https://crafatar.com/avatars/%v?size=100", utils.GetUUID(data.Options[0].Options[0].Value.String()))).
						SetTitle("Player Information").
						SetDescription("Showing infor for player by UUID").
						SetImage(fmt.Sprintf("https://crafatar.com/renders/body/%v", strings.ReplaceAll(data.Options[0].Options[0].Value.String(), "\"", ""))).
						AddField("Username", true, utils.GetUsername(data.Options[0].Options[0].Value.String())).
						AddField("UUID", true, strings.ReplaceAll(data.Options[0].Options[0].Value.String(), "\"", "")).
						AddField("Names", false, utils.GetNamesFromUUID(data.Options[0].Options[0].Value.String())).
						AddURLButton("Cape", fmt.Sprintf("https://crafatar.com/capes/%v", strings.ReplaceAll(data.Options[0].Options[0].Value.String(), "\"", ""))).
						AddURLButton("Skin", fmt.Sprintf("https://crafatar.com/skins/%v", strings.ReplaceAll(data.Options[0].Options[0].Value.String(), "\"", ""))).
						SetTimestamp().
						EditInteraction()

					if _, err := core.State.EditInteractionResponse(discord.AppID(utils.MustSnowflakeEnv(core.Config.APPID)), e.Token, res); err != nil {
						logger.Error(err)
					}
				}
			}
		},
	}
	MapCommands["server"] = Command{
		Name:        "server",
		Description: "Commands to interact with the minecraft server!",
		Group:       "minecraft",
		Usage:       "/server <option>",
		Options: []discord.CommandOption{
			{
				Type:        2,
				Name:        "info",
				Description: "Get information about the server or a specific player",
				Options: []discord.CommandOption{
					{
						Type:        1,
						Name:        "server",
						Description: "Get information about the server",
					},
					{
						Type:        1,
						Name:        "player",
						Description: "Get information about a player",
						Options: []discord.CommandOption{
							{
								Type:        3,
								Name:        "player",
								Description: "The specified player",
								Required:    true,
							},
						},
					},
				},
			},
		},
		OwnerOnly: false,
		Exclude:   false,
		Run: func(e *gateway.InteractionCreateEvent, data *discord.CommandInteractionData) {
			ack := api.InteractionResponse{
				Type: api.DeferredMessageInteractionWithSource,
			}

			if err := core.State.RespondInteraction(e.ID, e.Token, ack); err != nil {
				logger.Error(err)
			}
			switch data.Options[0].Name {
			case "info":
				switch data.Options[0].Options[0].Name {
				case "server":
					type Info struct {
						Tps             float32 `json:"tps"`
						AverageTickTime float32 `json:"averageTickTime"`
						OnlinePlayers   int     `json:"onlinePlayers"`
						AllPlayers      int     `json:"allPlayers"`
						BannedPlayers   int     `json:"bannedPlayers"`
						MOTD            string  `json:"motd"`
						Version         string  `json:"version"`
						MaxPlayers      int     `json:"maxPlayers"`
					}
					var msg Info

					core.WSServer.BroadcastToNamespace("/", "getserverinfo")
					core.WSServer.OnEvent("/", "getserverinfo", func(s socketio.Conn, obj string) {
						err := json.Unmarshal([]byte(obj), &msg)
						if err != nil {
							logger.Error("Failed to parse JSON message")
						}
						res := utils.NewEmbed().
							SetTitle("Showing server stats!").
							AddField("MOTD", true, fmt.Sprint(msg.MOTD)).
							AddField("TPS", true, fmt.Sprint(msg.Tps)).
							AddField("Average tick time", true, fmt.Sprint(msg.AverageTickTime)).
							AddField("Max players", true, fmt.Sprint(msg.MaxPlayers)).
							AddField("Online players", true, fmt.Sprint(msg.OnlinePlayers)).
							AddField("Banned players", true, fmt.Sprint(msg.BannedPlayers)).
							AddField("All players", true, fmt.Sprint(msg.AllPlayers)).
							AddField("Version", true, strings.ReplaceAll(strings.ReplaceAll(fmt.Sprint(msg.Version), "git-Paper-313 (MC: ", ""), ")", "")).
							SetTimestamp().
							SetColor(utils.DiscordBlue).
							EditInteraction()

						if _, err := core.State.EditInteractionResponse(discord.AppID(utils.MustSnowflakeEnv(core.Config.APPID)), e.Token, res); err != nil {
							logger.Error(err)
						}
					})
				case "player":
					type Info struct {
						Error        string  `json:"error,omitempty"`
						DisplayName  string  `json:"displayName,omitempty"`
						MaxHealth    float32 `json:"maxHealth,omitempty"`
						Health       float32 `json:"health,omitempty"`
						IsFlying     bool    `json:"isFlying,omitempty"`
						IsSleeping   bool    `json:"isSleeping,omitempty"`
						IsSneaking   bool    `json:"isSneaking,omitempty"`
						IsSprinting  bool    `json:"isSprinting,omitempty"`
						FirstPlayed  int64   `json:"firstPlayed,omitempty"`
						GameMode     string  `json:"gameMode,omitempty"`
						IsOp         bool    `json:"isOp,omitempty"`
						Online       bool    `json:"online,omitempty"`
						UUID         string  `json:"UUID,omitempty"`
						MobsKilled   int     `json:"mobsKilled,omitempty"`
						ItemsDropped int     `json:"itemsDropped,omitempty"`
						AnimalsBred  int     `json:"animalsBred,omitempty"`
						Deaths       int     `json:"deaths,omitempty"`
						GamesQuit    int     `json:"gamesQuit,omitempty"`
						TimePlayed   int     `json:"timePlayed,omitempty"`
						LastSeen     int64   `json:"lastSeen,omitempty"`
						IsBanned     bool    `json:"isBanned,omitempty"`
					}

					core.WSServer.BroadcastToNamespace("/", "getplayerinfo", data.Options[0].Options[0].Options[0].Value)
					core.WSServer.OnEvent("/", "getplayerinfo", func(s socketio.Conn, obj string) {
						var (
							res api.EditInteractionResponseData
							msg Info
						)
						err := json.Unmarshal([]byte(obj), &msg)
						if err != nil {
							logger.Error("Failed to parse JSON message: ", err)
						}
						if msg.Error != "" {
							res = utils.NewEmbed().
								SetColor(utils.DiscordRed).
								SetTitle("ERROR!").
								SetDescription(msg.Error).EditInteraction()
						} else {
							switch msg.Online {
							case true:
								res = utils.NewEmbed().
									SetTitle(fmt.Sprintf("Information for player %v", msg.DisplayName)).
									SetColor(utils.DiscordBlue).
									AddField("UUID", false, fmt.Sprint(msg.UUID)).
									AddField("Health", true, fmt.Sprintf("%v/%v", msg.Health, msg.MaxHealth)).
									AddField("First played", true, fmt.Sprint(time.UnixMilli(msg.FirstPlayed).UTC().Format(time.RFC822))).
									AddField("Last seen", true, fmt.Sprint(time.UnixMilli(msg.LastSeen).UTC().Format(time.RFC822))).
									AddField("Time played", true, utils.ConvertTickToDuration(msg.TimePlayed)).
									AddField("Is an operator?", true, fmt.Sprint(msg.IsOp)).
									AddField("Banned?", true, fmt.Sprint(msg.IsBanned)).
									AddField("Sleeping?", true, fmt.Sprint(msg.IsSleeping)).
									AddField("Flying?", true, fmt.Sprint(msg.IsFlying)).
									AddField("Sneaking?", true, fmt.Sprint(msg.IsSneaking)).
									AddField("Sprinting?", true, fmt.Sprint(msg.IsSprinting)).
									AddField("Online?", true, fmt.Sprint(msg.Online)).
									AddField("Gamemode", true, strings.Title(strings.ToLower(fmt.Sprint(msg.GameMode)))).
									AddField("Mobs killed", true, fmt.Sprint(msg.MobsKilled)).
									AddField("Animals bred", true, fmt.Sprint(msg.AnimalsBred)).
									AddField("Deaths", true, fmt.Sprint(msg.Deaths)).
									AddField("Games quit", true, fmt.Sprint(msg.GamesQuit)).
									SetImage(fmt.Sprintf("https://crafatar.com/renders/body/%v", msg.UUID)).
									SetTimestamp().
									EditInteraction()
							case false:
								res = utils.NewEmbed().
									SetTitle(fmt.Sprintf("Information for player %v", msg.DisplayName)).
									SetColor(utils.DiscordBlue).
									AddField("UUID", false, fmt.Sprint(msg.UUID)).
									AddField("Online?", true, fmt.Sprint(msg.Online)).
									AddField("First played", true, fmt.Sprint(time.UnixMilli(msg.FirstPlayed).UTC().Format(time.RFC822))).
									AddField("Last seen", true, fmt.Sprint(time.UnixMilli(msg.LastSeen).UTC().Format(time.RFC822))).
									AddField("Time played", true, fmt.Sprintf("%v", utils.ConvertTickToDuration(msg.TimePlayed))).
									AddField("Is an operator?", true, fmt.Sprint(msg.IsOp)).
									AddField("Banned?", true, fmt.Sprint(msg.IsBanned)).
									AddField("Mobs killed", true, fmt.Sprint(msg.MobsKilled)).
									AddField("Animals bred", true, fmt.Sprint(msg.AnimalsBred)).
									AddField("Deaths", true, fmt.Sprint(msg.Deaths)).
									AddField("Games quit", true, fmt.Sprint(msg.GamesQuit)).
									SetImage(fmt.Sprintf("https://crafatar.com/renders/body/%v", msg.UUID)).
									SetTimestamp().
									EditInteraction()
							}

						}
						if _, err := core.State.EditInteractionResponse(discord.AppID(utils.MustSnowflakeEnv(core.Config.APPID)), e.Token, res); err != nil {
							logger.Error(err)
						}
					})

				}

			}
		},
	}
}
