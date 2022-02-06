package main

import (
	"github.com/gorilla/websocket"
	"net/http"
	"phoenixbuilder/minecraft"
	"phoenixbuilder/minecraft/protocol/packet"
)

var CQMessages chan IMessage

var MCMessages chan *packet.Text
var Has_Connected bool
var Conn *minecraft.Conn
var ServerID string

// 将http升级为websocket
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
