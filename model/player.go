package model

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/jiny3/BanPickSys/pkg"
	"github.com/sirupsen/logrus"
)

type Player struct {
	ID     int64   `json:"id"`
	Name   string  `json:"name"`
	Banned []Entry `json:"banned"`
	Picked []Entry `json:"picked"`
	ws     *websocket.Conn
}

func NewPlayer(ws *websocket.Conn) Player {
	return Player{
		ID: pkg.GeneratePlayerUUID(),
		ws: ws,
	}
}

func (p *Player) Listen(ch chan *Game) {
	for {
		evt := <-ch
		// TODO: 处理来自game send的消息
		if err := p.ws.WriteJSON(gin.H{"bp": evt}); err != nil {
			logrus.WithError(err).Error("write json error")
			return
		}
	}
}

func (p *Player) WS() (*websocket.Conn, error) {
	if p.ws == nil {
		return nil, fmt.Errorf("wsconn is nil")
	}
	return p.ws, nil
}

func (p *Player) Ban(entry *Entry) {
	entry.Ban()
	p.Banned = append(p.Banned, *entry)
}

func (p *Player) Pick(entry *Entry) {
	entry.Pick()
	p.Picked = append(p.Picked, *entry)
}
