package pkg

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var WsUpgrade = websocket.Upgrader{
	HandshakeTimeout: time.Second * 3,
	ReadBufferSize:   256,
	WriteBufferSize:  256,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WsEvent struct {
	Type string `json:"type"`
	Data any    `json:"data"`
}
