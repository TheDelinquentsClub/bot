package websockets

import (
	"fmt"
	socketio "github.com/googollee/go-socket.io"
	"github.com/googollee/go-socket.io/engineio"
	"github.com/googollee/go-socket.io/engineio/transport"
	"github.com/googollee/go-socket.io/engineio/transport/polling"
	"github.com/googollee/go-socket.io/engineio/transport/websocket"
	"github.com/kingultron99/tdcbot/core"
	"github.com/kingultron99/tdcbot/logger"
	"log"
	"net/http"
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

func InitServer() {

	core.WSServer = Server

	core.WSServer.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		logger.Info("connected:", s.ID())
		return nil
	})

	core.WSServer.OnEvent("/", "hello", func(s socketio.Conn) {
		logger.Info(fmt.Sprintf("Hello from: %v", s))
	})

	RegisterMinecraftHandlers()
	RegisterWebsiteHandlers()

	core.WSServer.OnError("/", func(s socketio.Conn, e error) {
		logger.Error("meet error:", e)
	})

	core.WSServer.OnDisconnect("/", func(s socketio.Conn, reason string) {
		logger.Info("closed:", reason)
		if core.ServerConn != nil {
			if s.ID() == core.ServerConn.ID() {
				logger.Info("Server instance disconnected.")
				core.IsServerConnected = false
			}
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
