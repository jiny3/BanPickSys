package game

import (
	"context"
	"time"

	"github.com/jiny3/BanPickSys/model"
	"github.com/jiny3/BanPickSys/service"
	"github.com/sirupsen/logrus"
)

// Handlers 注册
var Handlers = map[string]service.GameInitFunc{
	"豹豹碰碰大作战": service.GameInitFunc(baobaoStages),
}

// 豹豹碰碰大作战 stages 编排
func baobaoStages(game *model.BP) error {
	// 确定 player 人数
	game.PlayerCap = 2

	// 确定每个阶段的时间
	banDuration := time.Minute * 1
	pickDuration := time.Minute * 1

	// 确定每个阶段的处理函数
	blueBanHandler := func(ctx context.Context, ch chan any) {
		logrus.Debug("蓝方禁用开始")
		select {
		case <-ctx.Done():
			logrus.Debug("蓝方禁用结束")
		case en := <-ch:
			game.Players["blue"].Ban(en.(*model.Entry))
			logrus.Debug("蓝方禁用", *(en.(*model.Entry)))
		}
	}
	redBanHandler := func(ctx context.Context, ch chan any) {
		logrus.Debug("红方禁用开始")
		select {
		case <-ctx.Done():
			logrus.Debug("红方禁用结束")
		case en := <-ch:
			game.Players["red"].Ban(en.(*model.Entry))
			logrus.Debug("红方禁用", *(en.(*model.Entry)))
		}
	}
	bluePickHandler := func(ctx context.Context, ch chan any) {
		logrus.Debug("蓝方选人开始")
		select {
		case <-ctx.Done():
			for i := range game.Entries {
				if !game.Entries[i].Banned && !game.Entries[i].Picked {
					game.Players["blue"].Pick(&game.Entries[i])
					break
				}
			}
		case en := <-ch:
			game.Players["blue"].Pick(en.(*model.Entry))
			logrus.Debug("蓝方选人", *en.(*model.Entry))
		}
	}
	redPickHandler := func(ctx context.Context, ch chan any) {
		logrus.Debug("红方选人开始")
		select {
		case <-ctx.Done():
			for i := range game.Entries {
				if !game.Entries[i].Banned && !game.Entries[i].Picked {
					game.Players["red"].Pick(&game.Entries[i])
					break
				}
			}
		case en := <-ch:
			game.Players["red"].Pick(en.(*model.Entry))
			logrus.Debug("红方选人", *(en.(*model.Entry)))
		}
	}
	endHandler := func(ctx context.Context, ch chan any) {
		logrus.Debug("游戏结束")
		time.Sleep(time.Minute * 10)
		delete(service.BPs, game.ID)
	}
	// blue ban 1
	cur := model.NewStage(game.NewStageId(), "蓝方禁用1", "蓝方禁用1次", "blue", banDuration, blueBanHandler)
	service.LinkStages(cur, game.Stage0)
	// red ban 1
	cur, pre := model.NewStage(game.NewStageId(), "红方禁用1", "红方禁用1次", "red", banDuration, redBanHandler), cur
	service.LinkStages(cur, pre)
	// blue pick 1
	cur, pre = model.NewStage(game.NewStageId(), "蓝方选人1", "蓝方选人1次", "blue", pickDuration, bluePickHandler), cur
	service.LinkStages(cur, pre)
	// red pick 2
	cur, pre = model.NewStage(game.NewStageId(), "红方选人2-1", "红方选人2次(第1次)", "red", pickDuration, redPickHandler), cur
	service.LinkStages(cur, pre)
	cur, pre = model.NewStage(game.NewStageId(), "红方选人2-2", "红方选人2次(第2次)", "red", pickDuration, redPickHandler), cur
	service.LinkStages(cur, pre)
	// blue pick 2
	cur, pre = model.NewStage(game.NewStageId(), "蓝方选人2-1", "蓝方选人2次(第1次)", "blue", pickDuration, bluePickHandler), cur
	service.LinkStages(cur, pre)
	cur, pre = model.NewStage(game.NewStageId(), "蓝方选人2-2", "蓝方选人2次(第2次)", "blue", pickDuration, bluePickHandler), cur
	service.LinkStages(cur, pre)
	// red pick 1
	cur, pre = model.NewStage(game.NewStageId(), "红方选人1", "红方选人1次", "red", pickDuration, redPickHandler), cur
	service.LinkStages(cur, pre)
	// end
	cur, pre = model.NewStage(game.NewStageId(), "游戏结束", "游戏结束", model.END, 0, endHandler), cur
	service.LinkStages(cur, pre)
	return nil
}
