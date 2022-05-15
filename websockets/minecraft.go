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
		msgObj.Advancement = strings.ReplaceAll(strings.ReplaceAll(msgObj.Advancement, ".title", ""), ".", "_")

		index := sort.StringSlice.Search(core.ItemIcons, strings.ToLower(msgObj.Icon))
		utils.GenerateAdvancement(utils.BasePath+core.ItemIcons[index], msgObj.Type, utils.GetLocale(msgObj.Advancement, "title"), msgObj.Player)

		wd, _ := os.Getwd()

		image, err := os.ReadFile(wd + fmt.Sprintf("/assets/generated/%v/%v_advancement.png", msgObj.Player, msgObj.Player))
		if err != nil {
			logger.Error(err)
		}
		reader := bytes.NewReader(image)

		body := api.SendMessageData{
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
		type Player struct {
			Uuid     string `json:"uuid"`
			Username string `json:"username"`
		}

		var (
			DBRes  string
			status string
			ID     string
			rand   = utils.RandString(10)
			player Player
		)

		err := json.Unmarshal([]byte(msg), &player)
		if err != nil {
			logger.Error("Failed to parse JSON message")
		}

		query, _ := core.DB.Query("SELECT MC_UUID, status, ID FROM players WHERE MC_UUID = ?", player.Uuid)
		for query.Next() {
			query.Scan(&DBRes, &status, &ID)
		}
		query.Close()

		if DBRes == "" && status != "PENDING" {
			stmt, _ := core.DB.Prepare("INSERT INTO players (MC_UUID, ID) VALUES (?, ?)")
			stmt.Exec(player.Uuid, rand)
			core.ServerConn.Emit("msg", "<gradient:#D8B4FE:#9333EA>TDC</gradient>", fmt.Sprintf("%v Hey there!", player.Username))
			core.ServerConn.Emit("msg", "<gradient:#D8B4FE:#9333EA>TDC</gradient>", fmt.Sprintf("%v We noticed you're new here!", player.Username))
			core.ServerConn.Emit("msg", "<gradient:#D8B4FE:#9333EA>TDC</gradient>", fmt.Sprintf("%v I've added your player ID to my database.", player.Username))
			core.ServerConn.Emit("msg", "<gradient:#D8B4FE:#9333EA>TDC</gradient>", fmt.Sprintf("%v jump onto our discord and type /server verify and enter your code: <hover:show_text:'Click to copy to clipboard!'><color:gold><click:copy_to_clipboard:'%v'>%v</click></color></hover> to verify this account with your discord", player.Username, rand, rand))
			core.ServerConn.Emit("msg", "<gradient:#D8B4FE:#9333EA>TDC</gradient>", fmt.Sprintf("%v This will give you and other players have access to some tools on discord that wont work otherwise!", player.Username))
		} else if status == "PENDING" {
			core.ServerConn.Emit("msg", "<gradient:#D8B4FE:#9333EA>TDC</gradient>", fmt.Sprintf("%v Just a friendly reminder to verify your account!", player.Username))
			core.ServerConn.Emit("msg", "<gradient:#D8B4FE:#9333EA>TDC</gradient>", fmt.Sprintf("%v jump onto our discord and type /server verify and enter your code: <hover:show_text:'Click to copy to clipboard!'><color:gold><click:copy_to_clipboard:'%v'>%v</click></color></hover> to verify this account with your discord", player.Username, ID, ID))
		}

		body, err := json.Marshal(Message{
			Username: "TDC Bot",
			Avatar:   "https://cdn.discordapp.com/avatars/769753889960361994/a4876fb3b263409750a0b93feb619386.webp?size=128",
			Embeds: &[]discord.Embed{
				{
					Title:     fmt.Sprintf("%v joined the server!", player.Username),
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
