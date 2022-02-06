package main

// go build --buildmode=plugin -gcflags="all=-N -l" ./plugin_beta/plugin/test_plugin.go
import (
	"fmt"
	"phoenixbuilder/minecraft/protocol/packet"
	conn "phoenixbuilder/plugin_beta"
)

type SingleQABot struct {
	message *packet.Text
	user    string
}

func (bot SingleQABot) Rule(pk packet.Packet) bool {
	switch pk.(type) {
	case *packet.Text:
		fmt.Println("[plugin]packet coming!")
		return true
	default:
		return false
	}
}

func (bot *SingleQABot) Init(conn *conn.PluginManager) {
	fmt.Println("init!")
	conn.RegisterPlugin(bot, true, true, 5, "SingleQABot")
}

func (bot *SingleQABot) Handler(conn *conn.PluginManager, pk packet.Packet) {
	// conn.Method.GetClientData
	fmt.Println("I receive :", pk)
	conn.Logger.Println("Handler!")
	bot.message = pk.(*packet.Text)

	if bot.message.Message != "留言" {
		return
	}
	bot.user = bot.message.SourceName
	fmt.Println("send chat: what?")
	conn.Method.SendChat("您的留言内容?")
	for {
		imsg := conn.ReadPacketFor(bot)
		fmt.Println(imsg)
		message := imsg.(*packet.Text)
		// fmt.Println(conn.Expand.GetClientData(), conn.Expand.GetGameData())
		bot.message = message
		if bot.message.SourceName != bot.user {
			continue
		} else {
			conn.Method.SendChat(bot.message.Message)
			conn.Logger.Println(bot.message)
			return
		}
	}
}

var Plugin SingleQABot
