package pkg

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var WsUpgrade = websocket.Upgrader{
	HandshakeTimeout: time.Second * 3,
	ReadBufferSize:   1024,
	WriteBufferSize:  1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WsEvent struct {
}
