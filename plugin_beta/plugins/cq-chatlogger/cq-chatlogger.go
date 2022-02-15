package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path"

	// "phoenixbuilder/minecraft"

	"phoenixbuilder/minecraft/protocol/packet"
	conn "phoenixbuilder/plugin_beta"
	"phoenixbuilder/plugin_beta/plugins/cq-chatlogger/comm/recv"
	"phoenixbuilder/plugin_beta/plugins/cq-chatlogger/comm/send"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type cqchat struct {
	conf ChatSettings
	conn *conn.PluginManager
}

func (cq *cqchat) Rule(pk packet.Packet) bool {
	// _, ok := pk.(*packet.Text)
	return true

}

func (cq *cqchat) Init(conn *conn.PluginManager) {
	// fmt.Println("cq init start")
	cq.conn = conn
	conn.RegisterPlugin(cq, true, false, 99, "CQChat")
	confdir, err := os.UserHomeDir()
	if err != nil {
		conn.Logger.Println(err)
	}
	confdir = path.Join(confdir, ".config/fastbuilder/plugins_beta")
	cq.conf, err = ReadSettings(confdir)
	if err != nil {
		fmt.Println("Failed to interpret config.")
	}
	// fmt.Println("cq init finish")
}

func (cq *cqchat) Handler(conn *conn.PluginManager, pk packet.Packet) {
	conn.Logger.Println("Start")
	// fmt.Println("cq handle start")
	// fmt.Println(cq.conf.Port)
	http.HandleFunc("/fastbuilder/cqchat", cq.handleFunc)
	err := http.ListenAndServe(cq.conf.Port, nil)
	if err != nil {
		cq.conn.Logger.Panicln("cqchat监听异常, 将重新尝试.")
		cq.Handler(conn, pk)
	}

}

func (cq *cqchat) handleFunc(w http.ResponseWriter, r *http.Request) {
	// fmt.Println("handle!")
	cqconn, err := upgrader.Upgrade(w, r, nil) //websocket链接通道
	// fmt.Println("here!")
	if err != nil {
		fmt.Println("err here:", err)
	}
	// 这个conn是和gocq连接的
	// 所有和游戏内的conn连接(发包)的操作全部在main.go里
	go cq.receiveMessage(cq.conn, cqconn)
	go cq.sendMessage(cq.conn, cqconn)
}

func (cq *cqchat) receiveMessage(conn *conn.PluginManager, cqconn *websocket.Conn) {
	for {
		msgType, data, err := cqconn.ReadMessage()
		if err != nil {
			fmt.Printf("ERR: %s", err)
			conn.Logger.Printf("ERR: %s", err)
		}
		if msgType != 1 {
			// fmt.Println(msgType)
			continue
		}
		post, err := recv.ParseMetaPost(data)
		if post.PostType == "meta_event" && post.MetaEventType == "lifecycle" {
			conn.Logger.Println("QQ Connection succeed.")
		}
		if post.PostType == "message" && err == nil {
			// fmt.Println(post)
			action, err := recv.ParseMessageData(data, cq.conf.GameMessageFormat, cq.conf.GroupNickname)
			if err != nil {
				continue
			}
			conn.Method.SendCommand(TellrawCommand(action.GameRawText, cq.conf.FilteredPlayerTag))
			// conn.Method.SendChat(action.GameRawText)
		}
	}
}

func (cq *cqchat) sendMessage(conn *conn.PluginManager, cqconn *websocket.Conn) {
	for {
		pk := conn.ReadPacketFor(cq)
		if pkmsg, ok := pk.(*packet.Text); ok {
			groupID := send.ParseGroupIDFrom(pkmsg.Message, cq.conf.GroupNickname)
			if groupID == -1 {
				groupID = cq.conf.DefaultGroupID
			}
			// don't send bot's message
			for _, filtered := range cq.conf.FilteredUsers {
				if filtered == pkmsg.SourceName {
					continue
				}
			}

			msg := send.FormatGameMsg(pkmsg.Message,
				cq.conf.QQMessageFormat,
				pkmsg.SourceName,
				// conn.Expand.GetGameData().WorldName,
				"",
			)

			action := send.GroupMessage(groupID, msg)
			data, err := json.Marshal(action)
			if err != nil {
				fmt.Println(err)
				continue
			}
			cqconn.WriteMessage(1, data)
		}
	}
}

var Plugin cqchat
