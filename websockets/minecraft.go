package websockets

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/utils/sendpart"
	socketio "github.com/googollee/go-socket.io"
	"github.com/kingultron99/tdcbot/core"
	"github.com/kingultron99/tdcbot/logger"
	"github.com/kingultron99/tdcbot/utils"
	"net/http"
	"os"
	"sort"
	"strings"
)

type MsgObj struct {
	Username string `json:"username"`
	Msg      string `json:"msg"`
}
type Advancement struct {
	Player      string
	Type        string
	Advancement string
	Icon        string
}

type Message struct {
	Username   string              `json:"username,omitempty"`
	Avatar     string              `json:"avatar_url,omitempty"`
	Content    string              `json:"content,omitempty"`
	Embeds     *[]discord.Embed    `json:"embeds,omitempty"`
	Components []discord.Component `json:"components,omitempty"`
	Files      []sendpart.File     `json:"files,omitempty"`
}

func RegisterMinecraftHandlers() {

	core.WSServer.OnEvent("/", "serverinstance", func(s socketio.Conn) {
		core.ServerConn = s
		core.IsServerConnected = true
	})

	core.WSServer.OnEvent("/", "playerchat", func(s socketio.Conn, msg string) {
		var msgObj MsgObj
		err := json.Unmarshal([]byte(msg), &msgObj)
		if err != nil {
			logger.Error("Failed to parse JSON message")
		}
		body, err := json.Marshal(Message{
			Username: msgObj.Username,
			Avatar:   fmt.Sprintf("https://crafatar.com/avatars/%v", utils.GetUUID(msgObj.Username)),
			Content:  msgObj.Msg,
		})
		if err != nil {
			logger.Error(err)
		}
		resp := bytes.NewBuffer(body)

		http.Post(
			core.Config.Webhook,
			"application/json",
			resp)
	})
	core.WSServer.OnEvent("/", "playeradvancement", func(s socketio.Conn, msg string) {
		var msgObj Advancement
		err := json.Unmarshal([]byte(msg), &msgObj)
		if err != nil {
			logger.Error("Failed to parse JSON message")
		}

		index := sort.StringSlice.Search(core.ItemIcons, strings.ToLower(msgObj.Icon))

		utils.GenerateAdvancement(utils.BasePath+core.ItemIcons[index], msgObj.Type, msgObj.Advancement)

		wd, _ := os.Getwd()

		image, err := os.ReadFile(wd + "/assets/generated/advancement.png")
		if err != nil {
			logger.Error(err)
		}
		reader := bytes.NewReader(image)

		body := api.SendMessageData{
			Embeds: []discord.Embed{
				{
					Author: &discord.EmbedAuthor{
						Name: msgObj.Player,
						Icon: fmt.Sprintf("https://crafatar.com/avatars/%v", utils.GetUUID(msgObj.Player)),
					},
					Color:     utils.DefaultColour,
					Title:     fmt.Sprint(msgObj.Player + " Did a thing!"),
					Timestamp: discord.NowTimestamp(),
				},
			},
			Files: []sendpart.File{
				{
					Name:   fmt.Sprintf("%v.png", msgObj.Type),
					Reader: reader,
				},
			},
		}

		_, err = core.State.SendMessageComplex(discord.ChannelID(utils.MustSnowflakeEnv(core.Config.BridgeChannelID)), body)
		if err != nil {
			logger.Error("Failed to send advancement message: ", err)
		}
	})
	core.WSServer.OnEvent("/", "playerjoin", func(s socketio.Conn, msg string) {

		body, err := json.Marshal(Message{
			Username: "TDC Bot",
			Avatar:   "https://cdn.discordapp.com/avatars/769753889960361994/a4876fb3b263409750a0b93feb619386.webp?size=128",
			Embeds: &[]discord.Embed{
				{
					Title:     fmt.Sprintf("%v joined the server!", msg),
					Color:     utils.DiscordGreen,
					Timestamp: discord.NowTimestamp(),
				},
			},
		})
		if err != nil {
			logger.Error(err)
		}
		resp := bytes.NewBuffer(body)

		http.Post(
			core.Config.Webhook,
			"application/json",
			resp)
	})
	core.WSServer.OnEvent("/", "playerleft", func(s socketio.Conn, player string, reason string) {
		body, err := json.Marshal(Message{
			Username: "TDC Bot",
			Avatar:   "https://cdn.discordapp.com/avatars/769753889960361994/a4876fb3b263409750a0b93feb619386.webp?size=128",
			Embeds: &[]discord.Embed{
				{
					Title: fmt.Sprintf("%v left the server!", player),
					Fields: []discord.EmbedField{
						{
							Name:   "Reason:",
							Inline: false,
							Value:  strings.Title(strings.ToLower(reason)),
						},
					},
					Color:     utils.DiscordRed,
					Timestamp: discord.NowTimestamp(),
				},
			},
		})
		if err != nil {
			logger.Error(err)
		}
		resp := bytes.NewBuffer(body)

		http.Post(
			core.Config.Webhook,
			"application/json",
			resp)
	})
}
