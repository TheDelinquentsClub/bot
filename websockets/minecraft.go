package websockets

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/diamondburned/arikawa/v3/discord"
	socketio "github.com/googollee/go-socket.io"
	"github.com/kingultron99/tdcbot/core"
	"github.com/kingultron99/tdcbot/logger"
	"github.com/kingultron99/tdcbot/utils"
	"net/http"
	"strings"
)

func RegisterMinecraftHandlers() {
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
		body, err := json.Marshal(Message{
			Username: msgObj.Player,
			Avatar:   fmt.Sprintf("https://crafatar.com/avatars/%v", utils.GetUUID(msgObj.Player)),
			Embeds: &[]discord.Embed{
				{
					Title: "Advancement Made!",
					Fields: []discord.EmbedField{
						{
							Name:   msgObj.Advancement,
							Inline: true,
							Value:  "type: " + strings.Title(strings.ToLower(msgObj.Type)),
						},
					},
					Timestamp: discord.NowTimestamp(),
					Color:     utils.DefaultColour,
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
