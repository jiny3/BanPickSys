package service

import (
	"context"
	"fmt"
	"time"

	"slices"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/jiny3/BanPickSys/model"
	"github.com/sirupsen/logrus"
)

var Players = make(map[int64]*model.Player)

func WsHandler(wsConn *websocket.Conn, bpID int64) error {
	game, err := GetBP(bpID)
	if err != nil {
		return err
	}
	player := model.NewPlayer(fmt.Sprintf("玩家%d", len(game.Players)))
	Players[player.ID] = player
	defer func() {
		delete(Players, player.ID)
		game.Leave(player)
	}()
	err = wsConn.WriteJSON(gin.H{
		"event": "init",
		"data": gin.H{
			"player_id": fmt.Sprintf("%d", player.ID),
		},
	})
	if err != nil {
		return err
	}
	recv := make(chan any)
	game.Send = append(game.Send, recv)
	defer func() {
		for i, ch := range game.Send {
			if ch == recv {
				game.Send = slices.Delete(game.Send, i, i+1)
				break
			}
		}
		close(recv)
	}()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		for {
			select {
			case <-ctx.Done():
				logrus.Info("ws connection closed")
				return
			case evt := <-recv:
				// TODO: 处理来自game send的消息
				if err := wsConn.WriteJSON(gin.H{"event": evt}); err != nil {
					logrus.WithError(err).Error("write json error")
					return
				}
			}
		}
	}()
	for {
		err = wsConn.WriteJSON(gin.H{
			"event": "heartbeat",
		})
		if err != nil {
			logrus.WithError(err).Debug("heartbeat send miss")
			return err
		}
		err = wsConn.ReadJSON(&gin.H{
			"event": "heartbeat",
		})
		if err != nil {
			logrus.WithError(err).Debug("heartbeat read miss")
			return err
		}
		<-time.After(3 * time.Second)
	}
}
