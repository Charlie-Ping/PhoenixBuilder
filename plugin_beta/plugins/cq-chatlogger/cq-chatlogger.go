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
}

func (cq *cqchat) Rule(pk packet.Packet) bool {
	// _, ok := pk.(*packet.Text)
	return true

}

func (cq *cqchat) Init(conn *conn.PluginManager) {
	conn.RegisterPlugin(cq, true, false, 99, "CQChat")
	confdir, err := os.UserHomeDir()
	if err != nil {
		conn.Logger.Println(err)
	}
	confdir = path.Join(confdir, ".config/fastbuilder/plugins_beta")
	cq.conf, err = ReadSettings(confdir)
}

func (cq *cqchat) Handler(conn *conn.PluginManager, pk packet.Packet) {

}

func (cq *cqchat) receiveMessage(conn *conn.PluginManager, cqconn *websocket.Conn) {
	for {
		msgType, data, err := cqconn.ReadMessage()
		if err != nil {
			conn.Logger.Printf("ERR: ", err)
		}
		if msgType != 0 {
			continue
		}
		post, err := recv.ParseMetaPost(data)
		if post.PostType == "meta_event" && post.MetaEventType == "lifecycle" {
			conn.Logger.Println("QQ Connection succeed.")
		}
		if post.PostType == "message" && err == nil {
			action, err := recv.ParseMessageData(data, cq.conf.GameMessageFormat, cq.conf.GroupNickname)
			if err != nil {
				continue
			}
			conn.Method.SendCommand(TellrawCommand(action.GameRawText, cq.conf.FilteredPlayerTag))
		}
	}
}

func (cq *cqchat) sendMessage(conn *conn.PluginManager, cqconn *websocket.Conn) {
	for {
		pk := conn.ReadPacketFor(cq)
		if pk, ok := pk.(*packet.Text); ok {
			groupID := send.ParseGroupIDFrom(pk.Message, cq.conf.GroupNickname)
			if groupID == -1 {
				groupID = cq.conf.DefaultGroupID
			}

			msg := send.FormatGameMsg(pk.Message,
				cq.conf.QQMessageFormat,
				pk.SourceName,
				conn.Expand.GetGameData().WorldName,
			)
			action := send.GroupMessage(groupID, msg)
			data, err := json.Marshal(action)
			if err != nil {
				fmt.Println(err)
				continue
			}
			cqconn.WriteMessage(0, data)
		}
	}
}
