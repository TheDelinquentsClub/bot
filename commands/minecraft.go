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
			&discord.SubcommandOption{
				OptionName:  "user",
				Description: "Gets information about a user, using either a valid username or UUID",
				Options: []discord.CommandOptionValue{
					&discord.StringOption{
						OptionName:  "username",
						Description: "Get user information from a username",
					},
					&discord.StringOption{
						OptionName:  "uuid",
						Description: "Get user information from a UUID",
					},
				},
			},
		},
		Restricted: false,
		Run: func(e *gateway.InteractionCreateEvent, data *discord.CommandInteraction) {
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
					res := api.EditInteractionResponseData{
						Embeds: &[]discord.Embed{
							{
								Author: &discord.EmbedAuthor{
									Name: strings.ReplaceAll(data.Options[0].Options[0].Value.String(), "\"", ""),
									Icon: fmt.Sprintf("https://crafatar.com/avatars/%v?size=100", utils.GetUUID(data.Options[0].Options[0].Value.String())),
								},
								Timestamp:   discord.NowTimestamp(),
								Color:       utils.DiscordBlue,
								Title:       "Player Information",
								Description: "Showing info for player by username",
								Image: &discord.EmbedImage{
									URL: fmt.Sprintf("https://crafatar.com/renders/body/%v", utils.GetUUID(data.Options[0].Options[0].Value.String())),
								},
								Fields: []discord.EmbedField{
									{
										Name:   "Username",
										Inline: true,
										Value:  strings.ReplaceAll(data.Options[0].Options[0].Value.String(), "\"", ""),
									},
									{
										Name:   "UUID",
										Inline: true,
										Value:  utils.GetUUID(data.Options[0].Options[0].Value.String()),
									},
									{
										Name:   "Names",
										Inline: false,
										Value:  utils.GetNamesFromUsername(data.Options[0].Options[0].Value.String()),
									},
								},
							},
						},
						Components: discord.ComponentsPtr(
							&discord.ActionRowComponent{
								&discord.ButtonComponent{
									Style: discord.LinkButtonStyle(fmt.Sprintf("https://crafatar.com/capes/%v", utils.GetUUID(data.Options[0].Options[0].Value.String()))),
									Label: "Cape",
								},
								&discord.ButtonComponent{
									Style: discord.LinkButtonStyle(fmt.Sprintf("https://crafatar.com/skins/%v", utils.GetUUID(data.Options[0].Options[0].Value.String()))),
									Label: "Skin",
								},
							}),
					}

					if _, err := core.State.EditInteractionResponse(discord.AppID(utils.MustSnowflakeEnv(core.Config.APPID)), e.Token, res); err != nil {
						logger.Error(err)
					}
				case "uuid":
					res := api.EditInteractionResponseData{
						Embeds: &[]discord.Embed{
							{
								Author: &discord.EmbedAuthor{
									Name: strings.ReplaceAll(data.Options[0].Options[0].Value.String(), "\"", ""),
									Icon: fmt.Sprintf("https://crafatar.com/avatars/%v?size=100", data.Options[0].Options[0].Value.String()),
								},
								Timestamp:   discord.NowTimestamp(),
								Title:       "Player Information",
								Description: "Showing info for player by username",
								Image: &discord.EmbedImage{
									URL: fmt.Sprintf("https://crafatar.com/renders/body/%v", data.Options[0].Options[0].Value.String()),
								},
								Fields: []discord.EmbedField{
									{
										Name:   "Username",
										Inline: true,
										Value:  utils.GetUsername(data.Options[0].Options[0].Value.String()),
									},
									{
										Name:   "UUID",
										Inline: true,
										Value:  strings.ReplaceAll(data.Options[0].Options[0].Value.String(), "\"", ""),
									},
									{
										Name:   "Names",
										Inline: false,
										Value:  utils.GetNamesFromUUID(data.Options[0].Options[0].Value.String()),
									},
								},
							},
						},
						Components: discord.ComponentsPtr(
							&discord.ActionRowComponent{
								&discord.ButtonComponent{
									Style: discord.LinkButtonStyle(fmt.Sprintf("https://crafatar.com/capes/%v", strings.ReplaceAll(data.Options[0].Options[0].Value.String(), "\"", ""))),
									Label: "Cape",
								},
								&discord.ButtonComponent{
									Style: discord.LinkButtonStyle(fmt.Sprintf("https://crafatar.com/skins/%v", strings.ReplaceAll(data.Options[0].Options[0].Value.String(), "\"", ""))),
									Label: "Skin",
								},
							}),
					}

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
		Usage:       "/server <subcommand> [option]",
		Options: []discord.CommandOption{
			&discord.SubcommandGroupOption{
				OptionName:  "info",
				Description: "Get information about the server or a specific player",
				Subcommands: []*discord.SubcommandOption{
					{
						OptionName:  "server",
						Description: "Get information about the server",
					},
					{
						OptionName:  "player",
						Description: "Get information about a player",
						Options: []discord.CommandOptionValue{
							&discord.StringOption{
								OptionName:  "player",
								Description: "The specified player",
								Required:    true,
							},
						},
					},
				},
			},
			&discord.SubcommandGroupOption{
				OptionName:  "commands",
				Description: "Invoke a command from the comfort of discord!",
				Subcommands: []*discord.SubcommandOption{
					{
						OptionName:  "command",
						Description: "The command you wish to run",
						Options: []discord.CommandOptionValue{
							&discord.StringOption{
								OptionName:  "name",
								Required:    true,
								Description: "The name of the command (more coming!)",
								Choices: []discord.StringChoice{
									{
										Name:  "announce",
										Value: "announce",
									},
									{
										Name:  "kill",
										Value: "kill",
									},
									{
										Name:  "msg",
										Value: "msg",
									},
									{
										Name:  "restart",
										Value: "restart",
									},
								},
							},
							&discord.StringOption{
								OptionName:  "args",
								Description: "Any arguments the command my need (such as player or message)",
							},
						},
					},
				},
			},
		},
		Restricted: false,
		Run: func(e *gateway.InteractionCreateEvent, data *discord.CommandInteraction) {
			switch core.IsServerConnected {
			case true:
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

						core.ServerConn.Emit("getserverinfo")
						core.WSServer.OnEvent("/", "getserverinfo", func(s socketio.Conn, obj string) {
							err := json.Unmarshal([]byte(obj), &msg)
							if err != nil {
								logger.Error("Failed to parse JSON message")
							}
							res := api.EditInteractionResponseData{
								Embeds: &[]discord.Embed{
									{
										Title: "Showing server stats!",
										Fields: []discord.EmbedField{
											{
												Name:   "MOTD",
												Inline: false,
												Value:  msg.MOTD,
											},
											{
												Name:   "TPS",
												Inline: true,
												Value:  fmt.Sprint(msg.Tps),
											},
											{
												Name:   "Average tick time",
												Inline: true,
												Value:  fmt.Sprint(msg.AverageTickTime),
											},
											{
												Name:   "Max players",
												Inline: true,
												Value:  fmt.Sprint(msg.MaxPlayers),
											},
											{
												Name:   "Online players",
												Inline: true,
												Value:  fmt.Sprint(msg.OnlinePlayers),
											},
											{
												Name:   "Banned players",
												Inline: true,
												Value:  fmt.Sprint(msg.BannedPlayers),
											},
											{
												Name:   "All players",
												Inline: true,
												Value:  fmt.Sprint(msg.AllPlayers),
											},
											{
												Name:   "Game Version",
												Inline: true,
												Value:  strings.ReplaceAll(strings.ReplaceAll(fmt.Sprint(msg.Version), "git-Paper-313 (MC: ", ""), ")", ""),
											},
										},
										Timestamp: discord.NowTimestamp(),
										Color:     utils.DiscordBlue,
									},
								},
							}
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

						core.ServerConn.Emit("getplayerinfo", data.Options[0].Options[0].Options[0].Value)
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
								res = api.EditInteractionResponseData{
									Embeds: &[]discord.Embed{
										{
											Color:       utils.DiscordRed,
											Title:       "ERROR!",
											Description: msg.Error,
										},
									},
								}
							} else {
								switch msg.Online {
								case true:
									res = api.EditInteractionResponseData{
										Embeds: &[]discord.Embed{
											{
												Title:     fmt.Sprintf("Information for player %v", msg.DisplayName),
												Color:     utils.DiscordBlue,
												Timestamp: discord.NowTimestamp(),
												Image: &discord.EmbedImage{
													URL: fmt.Sprintf("https://crafatar.com/renders/body/%v", msg.UUID),
												},
												Fields: []discord.EmbedField{
													{
														Name:   "UUID",
														Inline: false,
														Value:  fmt.Sprint(msg.UUID),
													},
													{
														Name:   "Health",
														Inline: true,
														Value:  fmt.Sprintf("%v/%v", msg.Health, msg.MaxHealth),
													},
													{
														Name:   "First played",
														Inline: true,
														Value:  fmt.Sprint(time.UnixMilli(msg.FirstPlayed).UTC().Format(time.RFC822)),
													},
													{
														Name:   "Last seen",
														Inline: true,
														Value:  fmt.Sprint(time.UnixMilli(msg.LastSeen).UTC().Format(time.RFC822)),
													},
													{
														Name:   "Time played",
														Inline: true,
														Value:  utils.ConvertTickToDuration(msg.TimePlayed),
													},
													{
														Name:   "Is an operator?",
														Inline: true,
														Value:  fmt.Sprint(msg.IsOp),
													},
													{
														Name:   "Banned?",
														Inline: true,
														Value:  fmt.Sprint(msg.IsBanned),
													},
													{
														Name:   "Sleeping?",
														Inline: true,
														Value:  fmt.Sprint(msg.IsSleeping),
													},
													{
														Name:   "Flying?",
														Inline: true,
														Value:  fmt.Sprint(msg.IsFlying),
													},
													{
														Name:   "Sneaking?",
														Inline: true,
														Value:  fmt.Sprint(msg.IsSneaking),
													},
													{
														Name:   "Sprinting?",
														Inline: true,
														Value:  fmt.Sprint(msg.IsSprinting),
													},
													{
														Name:   "Online?",
														Inline: true,
														Value:  fmt.Sprint(msg.Online),
													},
													{
														Name:   "Gamemode",
														Inline: true,
														Value:  strings.Title(strings.ToLower(fmt.Sprint(msg.GameMode))),
													},
													{
														Name:   "Mobs killed",
														Inline: true,
														Value:  fmt.Sprint(msg.MobsKilled),
													},
													{
														Name:   "Animals bred",
														Inline: true,
														Value:  fmt.Sprint(msg.AnimalsBred),
													},
													{
														Name:   "Deaths",
														Inline: true,
														Value:  fmt.Sprint(msg.Deaths),
													},
													{
														Name:   "Games quit",
														Inline: true,
														Value:  fmt.Sprint(msg.GamesQuit),
													},
												},
											},
										},
									}
								case false:
									res = api.EditInteractionResponseData{
										Embeds: &[]discord.Embed{
											{
												Title:     fmt.Sprintf("Information for player %v", msg.DisplayName),
												Color:     utils.DiscordBlue,
												Timestamp: discord.NowTimestamp(),
												Image: &discord.EmbedImage{
													URL: fmt.Sprintf("https://crafatar.com/renders/body/%v", msg.UUID),
												},
												Fields: []discord.EmbedField{
													{
														Name:   "UUID",
														Inline: false,
														Value:  fmt.Sprint(msg.UUID),
													},
													{
														Name:   "Online?",
														Inline: true,
														Value:  fmt.Sprint(msg.Online),
													},
													{
														Name:   "First played",
														Inline: true,
														Value:  fmt.Sprint(time.UnixMilli(msg.FirstPlayed).UTC().Format(time.RFC822)),
													},
													{
														Name:   "Last seen",
														Inline: true,
														Value:  fmt.Sprint(time.UnixMilli(msg.LastSeen).UTC().Format(time.RFC822)),
													},
													{
														Name:   "Time played",
														Inline: true,
														Value:  utils.ConvertTickToDuration(msg.TimePlayed),
													},
													{
														Name:   "Is an operator?",
														Inline: true,
														Value:  fmt.Sprint(msg.IsOp),
													},
													{
														Name:   "Banned?",
														Inline: true,
														Value:  fmt.Sprint(msg.IsBanned),
													},
													{
														Name:   "Mobs killed",
														Inline: true,
														Value:  fmt.Sprint(msg.MobsKilled),
													},
													{
														Name:   "Animals bred",
														Inline: true,
														Value:  fmt.Sprint(msg.AnimalsBred),
													},
													{
														Name:   "Deaths",
														Inline: true,
														Value:  fmt.Sprint(msg.Deaths),
													},
													{
														Name:   "Games quit",
														Inline: true,
														Value:  fmt.Sprint(msg.GamesQuit),
													},
												},
											},
										},
									}
								}

							}
							if _, err := core.State.EditInteractionResponse(discord.AppID(utils.MustSnowflakeEnv(core.Config.APPID)), e.Token, res); err != nil {
								logger.Error(err)
							}
						})

					}
				case "commands":
					switch data.Options[0].Options[0].Name {
					case "command":
						switch strings.ReplaceAll(data.Options[0].Options[0].Options[0].Value.String(), "\"", "") {
						case "announce":
							var res api.EditInteractionResponseData
							if len(data.Options[0].Options[0].Options) != 2 {
								res = api.EditInteractionResponseData{
									Embeds: &[]discord.Embed{
										{
											Title:       "You didnt provide a message!",
											Description: "Usage: /server commands command `name: announce` `args: <message>`",
											Color:       utils.DiscordRed,
											Timestamp:   discord.NowTimestamp(),
										},
									},
								}
							}
							core.ServerConn.Emit("announce", strings.ReplaceAll(data.Options[0].Options[0].Options[1].Value.String(), "\"", ""))
							res = api.EditInteractionResponseData{
								Embeds: &[]discord.Embed{
									{
										Title:       "Making an announcement",
										Description: "This should make them listen",
										Color:       utils.DiscordGreen,
										Timestamp:   discord.NowTimestamp(),
									},
								},
							}
							if _, err := core.State.EditInteractionResponse(discord.AppID(utils.MustSnowflakeEnv(core.Config.APPID)), e.Token, res); err != nil {
								logger.Error(err)
							}
						case "kill":
							var (
								notfound = false
								res      api.EditInteractionResponseData
							)
							if len(data.Options[0].Options[0].Options) != 2 {
								res = api.EditInteractionResponseData{
									Embeds: &[]discord.Embed{
										{
											Title:       "You didnt provide a player!",
											Description: "Usage: /server commands command `name: announce` `args: <player>`",
											Color:       utils.DiscordRed,
											Timestamp:   discord.NowTimestamp(),
										},
									},
								}
							}
							core.ServerConn.Emit("kill", strings.ReplaceAll(data.Options[0].Options[0].Options[1].Value.String(), "\"", ""))
							core.WSServer.OnEvent("/", "playernotfound", func(s socketio.Conn) {
								notfound = true
								return
							})
							if notfound != true {
								res = api.EditInteractionResponseData{
									Embeds: &[]discord.Embed{
										{
											Title:       "Killing player....",
											Description: "This shouldn't take too long",
											Color:       utils.DiscordGreen,
											Timestamp:   discord.NowTimestamp(),
										},
									},
								}
							} else {
								res = api.EditInteractionResponseData{
									Embeds: &[]discord.Embed{
										{
											Title:       "Player not found!!",
											Description: fmt.Sprintf("`%v` is either offline or does not exist", strings.ReplaceAll(data.Options[0].Options[0].Options[1].Value.String(), "\"", "")),
											Color:       utils.DiscordRed,
											Timestamp:   discord.NowTimestamp(),
										},
									},
								}
							}
							if _, err := core.State.EditInteractionResponse(discord.AppID(utils.MustSnowflakeEnv(core.Config.APPID)), e.Token, res); err != nil {
								logger.Error(err)
							}
						case "msg":
							var (
								notfound = false
								res      api.EditInteractionResponseData
							)
							if len(data.Options[0].Options[0].Options) != 2 {
								res = api.EditInteractionResponseData{
									Embeds: &[]discord.Embed{
										{
											Title:       "You didnt provide a player or message!!",
											Description: "Usage: /server commands command `name: announce` `args: <player> <message>`",
											Color:       utils.DiscordRed,
											Timestamp:   discord.NowTimestamp(),
										},
									},
								}
							}
							core.ServerConn.Emit("msg", strings.ReplaceAll(data.Options[0].Options[0].Options[1].Value.String(), "\"", ""))
							core.WSServer.OnEvent("/", "playernotfound", func(s socketio.Conn) {
								notfound = true
								return
							})
							if notfound != true {
								res = api.EditInteractionResponseData{
									Embeds: &[]discord.Embed{
										{
											Title:       "Sending message...",
											Description: "Who sent that??",
											Color:       utils.DiscordGreen,
											Timestamp:   discord.NowTimestamp(),
										},
									},
								}
							} else {
								res = api.EditInteractionResponseData{
									Embeds: &[]discord.Embed{
										{
											Title:       "Player not found!!",
											Description: fmt.Sprintf("`%v` is either offline or does not exist", strings.ReplaceAll(data.Options[0].Options[0].Options[1].Value.String(), "\"", "")),
											Color:       utils.DiscordRed,
											Timestamp:   discord.NowTimestamp(),
										},
									},
								}
							}
							if _, err := core.State.EditInteractionResponse(discord.AppID(utils.MustSnowflakeEnv(core.Config.APPID)), e.Token, res); err != nil {
								logger.Error(err)
							}
						case "restart":
							core.ServerConn.Emit("restart")
							var res api.EditInteractionResponseData
							res = api.EditInteractionResponseData{
								Embeds: &[]discord.Embed{
									{
										Title:       "Restarting server...",
										Description: "This shouldn't take too long",
										Color:       utils.DiscordGreen,
										Timestamp:   discord.NowTimestamp(),
									},
								},
							}
							if _, err := core.State.EditInteractionResponse(discord.AppID(utils.MustSnowflakeEnv(core.Config.APPID)), e.Token, res); err != nil {
								logger.Error(err)
							}
						}
					}
				}
			case false:
				res := api.InteractionResponse{
					Type: api.MessageInteractionWithSource,
					Data: &api.InteractionResponseData{
						Embeds: &[]discord.Embed{
							{
								Title:       "Server is not online!",
								Color:       utils.DiscordRed,
								Description: "Try again later...",
								Timestamp:   discord.NowTimestamp(),
							},
						},
					},
				}
				if err := core.State.RespondInteraction(e.ID, e.Token, res); err != nil {
					logger.Error(err)
				}
			}
		},
	}
}
