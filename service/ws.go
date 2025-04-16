package service

import (
	"github.com/gorilla/websocket"
	"github.com/jiny3/BanPickSys/model"
)

func WsHandler(wsConn *websocket.Conn, bpID int64) error {
	// 获取游戏实例
	game, err := GetGame(bpID)
	if err != nil {
		return err
	}

	// 创建玩家实例
	player := model.NewPlayer(wsConn)

	// 发送初始数据
	err = wsConn.WriteJSON(player)
	if err != nil {
		return err
	}
	recv := make(chan *model.Game)
	game.Send = append(game.Send, recv)

	player.Listen(recv)
	return nil
}
