package websockets

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/utils/sendpart"
	socketio "github.com/googollee/go-socket.io"
	"github.com/googollee/go-socket.io/engineio"
	"github.com/googollee/go-socket.io/engineio/transport"
	"github.com/googollee/go-socket.io/engineio/transport/polling"
	"github.com/googollee/go-socket.io/engineio/transport/websocket"
	"github.com/kingultron99/tdcbot/core"
	"github.com/kingultron99/tdcbot/logger"
	"github.com/kingultron99/tdcbot/utils"
	"log"
	"net/http"
	"strings"
)

var (
	Server = socketio.NewServer(&engineio.Options{
		Transports: []transport.Transport{
			&polling.Transport{
				CheckOrigin: allowOriginFunc,
			},
			&websocket.Transport{
				CheckOrigin: allowOriginFunc,
			},
		},
	})
	allowOriginFunc = func(r *http.Request) bool {
		return true
	}
)

type MsgObj struct {
	Username string `json:"username"`
	Msg      string `json:"msg"`
}

type Message struct {
	Username   string              `json:"username,omitempty"`
	Avatar     string              `json:"avatar_url,omitempty"`
	Content    string              `json:"content,omitempty"`
	Embeds     *[]discord.Embed    `json:"embeds,omitempty"`
	Components []discord.Component `json:"components,omitempty"`
	Files      []sendpart.File     `json:"files,omitempty"`
}

func InitServer() {

	core.WSServer = socketio.NewServer(nil)

	core.WSServer.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		logger.Info("connected:", s.ID())
		return nil
	})

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

	core.WSServer.OnError("/", func(s socketio.Conn, e error) {
		logger.Error("meet error:", e)
	})

	core.WSServer.OnDisconnect("/", func(s socketio.Conn, reason string) {
		logger.Info("closed:", reason)
		if s.ID() == core.ServerConn.ID() {
			logger.Info("Server instance disconnected.")
			core.IsServerConnected = false
		}
	})

	go func() {
		if err := core.WSServer.Serve(); err != nil {
			log.Fatalf("socketio listen error: %s\n", err)
		}
	}()
	defer core.WSServer.Close()

	http.Handle("/socket.io/", core.WSServer)
	logger.Info("Serving at localhost:3000...")
	logger.Fatal(http.ListenAndServe(":3000", nil))
}
