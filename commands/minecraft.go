package commands

import (
	"encoding/json"
	"fmt"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"strings"
	"time"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	socketio "github.com/googollee/go-socket.io"
	"github.com/kingultron99/tdcbot/core"
	"github.com/kingultron99/tdcbot/logger"
	"github.com/kingultron99/tdcbot/utils"
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
				if len(data.Options[0].Options) != 1 {
					res := api.EditInteractionResponseData{
						Embeds: &[]discord.Embed{
							{
								Title:       "You didnt provide a user or UUID!!",
								Description: "Usage: `/minecraft user [username|UUID]`",
								Color:       utils.DiscordRed,
								Timestamp:   discord.NowTimestamp(),
							},
						},
					}
					if _, err := core.State.EditInteractionResponse(discord.AppID(utils.MustSnowflakeEnv(core.Config.APPID)), e.Token, res); err != nil {
						logger.Error(err)
					}
				} else {
					switch utils.GetUUID(data.Options[0].Options[0].Value.String()) {
					case "":
						res := api.EditInteractionResponseData{
							Embeds: &[]discord.Embed{
								{
									Title:       "Player does not exist!",
									Description: "Make sure you entered the correct UUID or username",
									Color:       utils.DiscordRed,
									Timestamp:   discord.NowTimestamp(),
								},
							},
						}
						if _, err := core.State.EditInteractionResponse(discord.AppID(utils.MustSnowflakeEnv(core.Config.APPID)), e.Token, res); err != nil {
							logger.Error(err)
						}
					default:
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
			&discord.SubcommandOption{
				OptionName:  "verify",
				Description: "Link your discord and minecraft using the unique code provided!",
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
				var hasPerms bool
				for _, d := range e.Member.RoleIDs {
					if d.String() == core.Config.CreatorID || d.String() == core.Config.BotBreakerRole {
						hasPerms = true
					} else {
						hasPerms = false
					}
				}

				switch data.Options[0].Name {
				case "info":
					ack := api.InteractionResponse{
						Type: api.DeferredMessageInteractionWithSource,
					}

					if err := core.State.RespondInteraction(e.ID, e.Token, ack); err != nil {
						logger.Error(err)
					}
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
					ack := api.InteractionResponse{
						Type: api.DeferredMessageInteractionWithSource,
					}

					if err := core.State.RespondInteraction(e.ID, e.Token, ack); err != nil {
						logger.Error(err)
					}
					switch data.Options[0].Options[0].Name {
					case "command":
						switch strings.ReplaceAll(data.Options[0].Options[0].Options[0].Value.String(), "\"", "") {
						case "announce":
							var res api.EditInteractionResponseData
							if hasPerms != true {
								res = api.EditInteractionResponseData{
									Embeds: &[]discord.Embed{
										{
											Title:     "You don't have permission to execute this command!",
											Color:     utils.DiscordRed,
											Timestamp: discord.NowTimestamp(),
										},
									},
								}
							} else if len(data.Options[0].Options[0].Options) != 2 {
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
							} else {
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
							}
							if _, err := core.State.EditInteractionResponse(discord.AppID(utils.MustSnowflakeEnv(core.Config.APPID)), e.Token, res); err != nil {
								logger.Error(err)
							}
						case "kill":
							var (
								notfound = false
								res      api.EditInteractionResponseData
							)
							if hasPerms != true {
								res = api.EditInteractionResponseData{
									Embeds: &[]discord.Embed{
										{
											Title:     "You don't have permission to execute this command!",
											Color:     utils.DiscordRed,
											Timestamp: discord.NowTimestamp(),
										},
									},
								}
							} else if len(data.Options[0].Options[0].Options) != 2 {
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
							} else {
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
							} else {
								core.ServerConn.Emit("msg", e.Member.User.ID, strings.ReplaceAll(data.Options[0].Options[0].Options[1].Value.String(), "\"", ""))
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
							}
							if _, err := core.State.EditInteractionResponse(discord.AppID(utils.MustSnowflakeEnv(core.Config.APPID)), e.Token, res); err != nil {
								logger.Error(err)
							}
						case "restart":
							var res api.EditInteractionResponseData
							if hasPerms != true {
								res = api.EditInteractionResponseData{
									Embeds: &[]discord.Embed{
										{
											Title:     "You don't have permission to execute this command!",
											Color:     utils.DiscordRed,
											Timestamp: discord.NowTimestamp(),
										},
									},
								}
							} else {
								core.ServerConn.Emit("restart")
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
							}
							if _, err := core.State.EditInteractionResponse(discord.AppID(utils.MustSnowflakeEnv(core.Config.APPID)), e.Token, res); err != nil {
								logger.Error(err)
							}
						}
					}
				case "verify":
					logger.Debug("begining player registration!")
					var res api.InteractionResponse
					// Check if user has already claimed an account.
					check1, err := core.DB.Query("SELECT Discord_UUID, MC_UUID from players where Discord_UUID = ?", e.Member.User.ID.String())
					if err != nil {
						logger.Error("failed to check for existing accounts tied to:", e.Member.User.ID.String(), "error:", err)
					}
					for check1.Next() {
						var (
							mc        string
							discordid string
						)
						check1.Scan(&discordid, &mc)

						if discordid != "" {
							res = api.InteractionResponse{
								Type: api.MessageInteractionWithSource,
								Data: &api.InteractionResponseData{
									Flags: api.EphemeralResponse,
									Embeds: &[]discord.Embed{
										{
											Color:       utils.DiscordRed,
											Title:       "You've already claimed an account!",
											Description: "If this is an error please let <@148203660088705025> know immediately in <#974302886945243156>",
											Fields: []discord.EmbedField{
												{
													Name:   "Minecraft username:",
													Value:  utils.GetUsername(mc),
													Inline: true,
												},
												{
													Name:   "Minecraft UUID",
													Value:  mc,
													Inline: true,
												},
											},
										},
									},
								},
							}
							if err := core.State.RespondInteraction(e.ID, e.Token, res); err != nil {
								logger.Error(err)
							}
							return
						}
					}
					check1.Close()

					res = api.InteractionResponse{
						Type: api.ModalResponse,
						Data: &api.InteractionResponseData{
							Title:    option.NewNullableString("Account verification"),
							Content:  option.NewNullableString("Enter the your provided code here to verify your account!"),
							CustomID: option.NewNullableString("verification"),
							Components: discord.ComponentsPtr(
								&discord.ActionRowComponent{
									&discord.TextInputComponent{
										Style:       discord.TextInputShortStyle,
										Label:       "Enter code here",
										Required:    true,
										CustomID:    "verification_code",
										Placeholder: option.NewNullableString("code"),
									},
								}),
						},
					}

					if err := core.State.RespondInteraction(e.ID, e.Token, res); err != nil {
						logger.Error(err)
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
	MapCommands["Player Information"] = Command{
		Name:        "Player Information",
		Description: "",
		Group:       "minecraft",
		Usage:       "Right click on a user and select `Apps > Player Information` from the menu!",
		Options:     nil,
		Type:        discord.UserCommand,
		Restricted:  false,
		Run: func(e *gateway.InteractionCreateEvent, data *discord.CommandInteraction) {
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
			me, _ := core.State.Me()
			for id, _ := range data.Resolved.Users {
				if id == me.ID {
					resp := api.InteractionResponse{
						Type: api.MessageInteractionWithSource,
						Data: &api.InteractionResponseData{
							Embeds: &[]discord.Embed{
								{
									Color:       utils.DiscordRed,
									Title:       "Woah there!",
									Description: "You wont be able to get any player stats about me... I cant play minecraft!\n\nYou can always try someone else though...",
									Timestamp:   discord.NowTimestamp(),
								},
							},
						},
					}
					if err := core.State.RespondInteraction(e.ID, e.Token, resp); err != nil {
						logger.Error(err)
					}
					return
				}
				res, err := core.DB.Query("SELECT MC_UUID from players where Discord_UUID = ? and status != 'PENDING'", id)
				if err != nil {
					logger.Error(err)
				}

				var uuid string

				for res.Next() {
					res.Scan(&uuid)
				}

				if core.IsServerConnected == true {
					if uuid == "" {
						logger.Error("Failed to verify player: Already linked to another account!")
						resp := api.InteractionResponse{
							Type: api.MessageInteractionWithSource,
							Data: &api.InteractionResponseData{
								Embeds: &[]discord.Embed{
									{
										Color:       utils.DiscordRed,
										Title:       "ERROR!",
										Description: "This user has either not linked their accounts or has not played on the server before",
										Timestamp:   discord.NowTimestamp(),
									},
								},
							},
						}
						if err := core.State.RespondInteraction(e.ID, e.Token, resp); err != nil {
							logger.Error(err)
						}
						return
					}
					ack := api.InteractionResponse{
						Type: api.DeferredMessageInteractionWithSource,
					}

					if err := core.State.RespondInteraction(e.ID, e.Token, ack); err != nil {
						logger.Error(err)
					}
					core.ServerConn.Emit("getplayerinfo", uuid)
					core.WSServer.OnEvent("/", "getplayerinfo", func(s socketio.Conn, obj string) {
						var (
							res api.EditInteractionResponseData
							msg *Info
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
				} else {
					resp := api.InteractionResponse{
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
					if err := core.State.RespondInteraction(e.ID, e.Token, resp); err != nil {
						logger.Error(err)
					}
				}
			}

			//core.ServerConn.Emit("getplayerinfo", data.Options[0].Options[0].Options[0].Value)
			//core.WSServer.OnEvent("/", "getplayerinfo", func(s socketio.Conn, obj string) {
			//	var (
			//		res api.EditInteractionResponseData
			//		msg Info
			//	)
			//	err := json.Unmarshal([]byte(obj), &msg)
			//	if err != nil {
			//		logger.Error("Failed to parse JSON message: ", err)
			//	}
			//	if msg.Error != "" {
			//		res = api.EditInteractionResponseData{
			//			Embeds: &[]discord.Embed{
			//				{
			//					Color:       utils.DiscordRed,
			//					Title:       "ERROR!",
			//					Description: msg.Error,
			//				},
			//			},
			//		}
			//	} else {
			//		switch msg.Online {
			//		case true:
			//			res = api.EditInteractionResponseData{
			//				Embeds: &[]discord.Embed{
			//					{
			//						Title:     fmt.Sprintf("Information for player %v", msg.DisplayName),
			//						Color:     utils.DiscordBlue,
			//						Timestamp: discord.NowTimestamp(),
			//						Image: &discord.EmbedImage{
			//							URL: fmt.Sprintf("https://crafatar.com/renders/body/%v", msg.UUID),
			//						},
			//						Fields: []discord.EmbedField{
			//							{
			//								Name:   "UUID",
			//								Inline: false,
			//								Value:  fmt.Sprint(msg.UUID),
			//							},
			//							{
			//								Name:   "Health",
			//								Inline: true,
			//								Value:  fmt.Sprintf("%v/%v", msg.Health, msg.MaxHealth),
			//							},
			//							{
			//								Name:   "First played",
			//								Inline: true,
			//								Value:  fmt.Sprint(time.UnixMilli(msg.FirstPlayed).UTC().Format(time.RFC822)),
			//							},
			//							{
			//								Name:   "Last seen",
			//								Inline: true,
			//								Value:  fmt.Sprint(time.UnixMilli(msg.LastSeen).UTC().Format(time.RFC822)),
			//							},
			//							{
			//								Name:   "Time played",
			//								Inline: true,
			//								Value:  utils.ConvertTickToDuration(msg.TimePlayed),
			//							},
			//							{
			//								Name:   "Is an operator?",
			//								Inline: true,
			//								Value:  fmt.Sprint(msg.IsOp),
			//							},
			//							{
			//								Name:   "Banned?",
			//								Inline: true,
			//								Value:  fmt.Sprint(msg.IsBanned),
			//							},
			//							{
			//								Name:   "Sleeping?",
			//								Inline: true,
			//								Value:  fmt.Sprint(msg.IsSleeping),
			//							},
			//							{
			//								Name:   "Flying?",
			//								Inline: true,
			//								Value:  fmt.Sprint(msg.IsFlying),
			//							},
			//							{
			//								Name:   "Sneaking?",
			//								Inline: true,
			//								Value:  fmt.Sprint(msg.IsSneaking),
			//							},
			//							{
			//								Name:   "Sprinting?",
			//								Inline: true,
			//								Value:  fmt.Sprint(msg.IsSprinting),
			//							},
			//							{
			//								Name:   "Online?",
			//								Inline: true,
			//								Value:  fmt.Sprint(msg.Online),
			//							},
			//							{
			//								Name:   "Gamemode",
			//								Inline: true,
			//								Value:  strings.Title(strings.ToLower(fmt.Sprint(msg.GameMode))),
			//							},
			//							{
			//								Name:   "Mobs killed",
			//								Inline: true,
			//								Value:  fmt.Sprint(msg.MobsKilled),
			//							},
			//							{
			//								Name:   "Animals bred",
			//								Inline: true,
			//								Value:  fmt.Sprint(msg.AnimalsBred),
			//							},
			//							{
			//								Name:   "Deaths",
			//								Inline: true,
			//								Value:  fmt.Sprint(msg.Deaths),
			//							},
			//							{
			//								Name:   "Games quit",
			//								Inline: true,
			//								Value:  fmt.Sprint(msg.GamesQuit),
			//							},
			//						},
			//					},
			//				},
			//			}
			//		case false:
			//			res = api.EditInteractionResponseData{
			//				Embeds: &[]discord.Embed{
			//					{
			//						Title:     fmt.Sprintf("Information for player %v", msg.DisplayName),
			//						Color:     utils.DiscordBlue,
			//						Timestamp: discord.NowTimestamp(),
			//						Image: &discord.EmbedImage{
			//							URL: fmt.Sprintf("https://crafatar.com/renders/body/%v", msg.UUID),
			//						},
			//						Fields: []discord.EmbedField{
			//							{
			//								Name:   "UUID",
			//								Inline: false,
			//								Value:  fmt.Sprint(msg.UUID),
			//							},
			//							{
			//								Name:   "Online?",
			//								Inline: true,
			//								Value:  fmt.Sprint(msg.Online),
			//							},
			//							{
			//								Name:   "First played",
			//								Inline: true,
			//								Value:  fmt.Sprint(time.UnixMilli(msg.FirstPlayed).UTC().Format(time.RFC822)),
			//							},
			//							{
			//								Name:   "Last seen",
			//								Inline: true,
			//								Value:  fmt.Sprint(time.UnixMilli(msg.LastSeen).UTC().Format(time.RFC822)),
			//							},
			//							{
			//								Name:   "Time played",
			//								Inline: true,
			//								Value:  utils.ConvertTickToDuration(msg.TimePlayed),
			//							},
			//							{
			//								Name:   "Is an operator?",
			//								Inline: true,
			//								Value:  fmt.Sprint(msg.IsOp),
			//							},
			//							{
			//								Name:   "Banned?",
			//								Inline: true,
			//								Value:  fmt.Sprint(msg.IsBanned),
			//							},
			//							{
			//								Name:   "Mobs killed",
			//								Inline: true,
			//								Value:  fmt.Sprint(msg.MobsKilled),
			//							},
			//							{
			//								Name:   "Animals bred",
			//								Inline: true,
			//								Value:  fmt.Sprint(msg.AnimalsBred),
			//							},
			//							{
			//								Name:   "Deaths",
			//								Inline: true,
			//								Value:  fmt.Sprint(msg.Deaths),
			//							},
			//							{
			//								Name:   "Games quit",
			//								Inline: true,
			//								Value:  fmt.Sprint(msg.GamesQuit),
			//							},
			//						},
			//					},
			//				},
			//			}
			//		}
			//
			//	}
			//	if _, err := core.State.EditInteractionResponse(discord.AppID(utils.MustSnowflakeEnv(core.Config.APPID)), e.Token, res); err != nil {
			//		logger.Error(err)
			//	}
			//})
		},
	}
}
