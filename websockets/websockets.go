package websockets

import (
	"bytes"
	"encoding/json"
	"fmt"
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

func InitServer() {

	core.WSServer = socketio.NewServer(nil)

	core.WSServer.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		logger.Info("connected:", s.ID())
		return nil
	})

	core.WSServer.OnEvent("/", "playerchat", func(s socketio.Conn, msg string) {
		var msgObj MsgObj
		err := json.Unmarshal([]byte(msg), &msgObj)
		if err != nil {
			logger.Error("Failed to parse JSON message")
		}
		embed, err := utils.NewEmbed().
			SetText(msgObj.Msg).
			MakeWebhookText(msgObj.Username, fmt.Sprintf("https://crafatar.com/avatars/%v", utils.GetUUID(msgObj.Username)))
		if err != nil {
			logger.Error(err)
		}
		resp := bytes.NewBuffer(embed)

		http.Post(
			core.Config.Webhook,
			"application/json",
			resp)
	})
	core.WSServer.OnEvent("/", "playerjoin", func(s socketio.Conn, msg string) {
		embed, err := utils.NewEmbed().
			SetTitle(fmt.Sprintf("%v joined the server!", msg)).
			SetTimestamp().
			SetColor(utils.DiscordGreen).
			MakeWebhookEmbed("TDC Bot", "https://cdn.discordapp.com/avatars/769753889960361994/a4876fb3b263409750a0b93feb619386.webp?size=128")
		if err != nil {
			logger.Error(err)
		}
		resp := bytes.NewBuffer(embed)

		http.Post(
			core.Config.Webhook,
			"application/json",
			resp)
	})
	core.WSServer.OnEvent("/", "playerleft", func(s socketio.Conn, player string, reason string) {
		embed, err := utils.NewEmbed().
			SetTitle(fmt.Sprintf("%v left the server!", player)).
			AddField("Reason:", false, strings.Title(strings.ToLower(reason))).
			SetTimestamp().
			SetColor(utils.DiscordRed).
			MakeWebhookEmbed("TDC Bot", "https://cdn.discordapp.com/avatars/769753889960361994/a4876fb3b263409750a0b93feb619386.webp?size=128")
		if err != nil {
			logger.Error(err)
		}
		resp := bytes.NewBuffer(embed)

		http.Post(
			core.Config.Webhook,
			"application/json",
			resp)
	})

	core.WSServer.OnError("/", func(s socketio.Conn, e error) {
		fmt.Println("meet error:", e)
	})

	core.WSServer.OnDisconnect("/", func(s socketio.Conn, reason string) {
		fmt.Println("closed", reason)
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
