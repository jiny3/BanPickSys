package api

import (
	"time"

	"github.com/jiny3/BanPickSys/model"
	"github.com/jiny3/BanPickSys/service"
	"github.com/sirupsen/logrus"
)

var GameHandler = map[string]service.GameInitFunc{
	"豹豹碰碰大作战": service.GameInitFunc(baobaoStages),
}

// 豹豹碰碰大作战 stages 编排
func baobaoStages(game *model.Game) error {
	// 确定 player 人数
	game.PlayerCap = 2

	banDuration := time.Minute * 1
	pickDuration := time.Minute * 1
	blueBanHandler := func(done chan struct{}, ch chan *model.Entry) {
		logrus.Debug("蓝方禁用开始")
		select {
		case <-done:
			logrus.Debug("蓝方禁用结束")
		case en := <-ch:
			game.Players[0].Ban(en)
			logrus.Debug("蓝方禁用", *en)
			done <- struct{}{}
		}
	}
	redBanHandler := func(done chan struct{}, ch chan *model.Entry) {
		logrus.Debug("红方禁用开始")
		select {
		case <-done:
			logrus.Debug("红方禁用结束")
		case en := <-ch:
			game.Players[1].Ban(en)
			logrus.Debug("红方禁用", *en)
			done <- struct{}{}
		}
	}
	bluePickHandler := func(done chan struct{}, ch chan *model.Entry) {
		logrus.Debug("蓝方选人开始")
		select {
		case <-done:
			for i := range game.Entries {
				if !game.Entries[i].Banned && !game.Entries[i].Picked {
					game.Players[0].Pick(&game.Entries[i])
					break
				}
			}
		case en := <-ch:
			game.Players[0].Pick(en)
			logrus.Debug("蓝方选人", *en)
			done <- struct{}{}
		}
	}
	redPickHandler := func(done chan struct{}, ch chan *model.Entry) {
		logrus.Debug("红方选人开始")
		select {
		case <-done:
			for i := range game.Entries {
				if !game.Entries[i].Banned && !game.Entries[i].Picked {
					game.Players[1].Pick(&game.Entries[i])
					break
				}
			}
		case en := <-ch:
			game.Players[1].Pick(en)
			logrus.Debug("红方选人", *en)
			done <- struct{}{}
		}
	}
	endHandler := func(done chan struct{}, ch chan *model.Entry) {
		logrus.Debug("游戏结束")
		time.Sleep(time.Minute * 10)
		delete(service.Games, game.ID)
	}
	// blue ban 1
	cur := model.NewStage(game.NewStageId(), "蓝方禁用1", "蓝方禁用1次", &game.Players[0], banDuration, blueBanHandler)
	service.LinkStages(cur, game.Stage0)
	// red ban 1
	cur, pre := model.NewStage(game.NewStageId(), "红方禁用1", "红方禁用1次", &game.Players[1], banDuration, redBanHandler), cur
	service.LinkStages(cur, pre)
	// blue pick 1
	cur, pre = model.NewStage(game.NewStageId(), "蓝方选人1", "蓝方选人1次", &game.Players[0], pickDuration, bluePickHandler), cur
	service.LinkStages(cur, pre)
	// red pick 2
	cur, pre = model.NewStage(game.NewStageId(), "红方选人2-1", "红方选人2次(第1次)", &game.Players[1], pickDuration, redPickHandler), cur
	service.LinkStages(cur, pre)
	cur, pre = model.NewStage(game.NewStageId(), "红方选人2-2", "红方选人2次(第2次)", &game.Players[1], pickDuration, redPickHandler), cur
	service.LinkStages(cur, pre)
	// blue pick 2
	cur, pre = model.NewStage(game.NewStageId(), "蓝方选人2-1", "蓝方选人2次(第1次)", &game.Players[0], pickDuration, bluePickHandler), cur
	service.LinkStages(cur, pre)
	cur, pre = model.NewStage(game.NewStageId(), "蓝方选人2-2", "蓝方选人2次(第2次)", &game.Players[0], pickDuration, bluePickHandler), cur
	service.LinkStages(cur, pre)
	// red pick 1
	cur, pre = model.NewStage(game.NewStageId(), "红方选人1", "红方选人1次", &game.Players[1], pickDuration, redPickHandler), cur
	service.LinkStages(cur, pre)
	// end
	cur, pre = model.NewStage(game.NewStageId(), "游戏结束", "游戏结束", &model.Player{ID: 0}, 0, endHandler), cur
	service.LinkStages(cur, pre)
	return nil
}
