package service

import (
	"fmt"
	"time"

	"github.com/jiny3/BanPickSys/model"
	"github.com/sirupsen/logrus"
)

var Games = make(map[int64]*model.Game)

var gameCreator = map[string]func(game *model.Game) error{
	"豹豹碰碰大作战": baobaoStages,
}

var baobaoStages = func(game *model.Game) error {
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
		delete(Games, game.ID)
	}
	// blue ban 1
	cur := model.NewStage(game.NewStageId(), "蓝方禁用1", "蓝方禁用1次", &game.Players[0], banDuration, blueBanHandler)
	LinkStages(cur, game.Stage0)
	// red ban 1
	cur, pre := model.NewStage(game.NewStageId(), "红方禁用1", "红方禁用1次", &game.Players[1], banDuration, redBanHandler), cur
	LinkStages(cur, pre)
	// blue pick 1
	cur, pre = model.NewStage(game.NewStageId(), "蓝方选人1", "蓝方选人1次", &game.Players[0], pickDuration, bluePickHandler), cur
	LinkStages(cur, pre)
	// red pick 2
	cur, pre = model.NewStage(game.NewStageId(), "红方选人2-1", "红方选人2次(第1次)", &game.Players[1], pickDuration, redPickHandler), cur
	LinkStages(cur, pre)
	cur, pre = model.NewStage(game.NewStageId(), "红方选人2-2", "红方选人2次(第2次)", &game.Players[1], pickDuration, redPickHandler), cur
	LinkStages(cur, pre)
	// blue pick 2
	cur, pre = model.NewStage(game.NewStageId(), "蓝方选人2-1", "蓝方选人2次(第1次)", &game.Players[0], pickDuration, bluePickHandler), cur
	LinkStages(cur, pre)
	cur, pre = model.NewStage(game.NewStageId(), "蓝方选人2-2", "蓝方选人2次(第2次)", &game.Players[0], pickDuration, bluePickHandler), cur
	LinkStages(cur, pre)
	// red pick 1
	cur, pre = model.NewStage(game.NewStageId(), "红方选人1", "红方选人1次", &game.Players[1], pickDuration, redPickHandler), cur
	LinkStages(cur, pre)
	// end
	cur, pre = model.NewStage(game.NewStageId(), "游戏结束", "游戏结束", &model.Player{ID: 0}, 0, endHandler), cur
	LinkStages(cur, pre)
	return nil
}

func NewGame(name string, players []model.Player) (int64, error) {
	if _, ok := gameCreator[name]; !ok {
		return -1, fmt.Errorf("unsupported game: %s", name)
	}
	game := model.NewGame(name, GetEntries(name), players)
	h := gameCreator[name]
	err := h(&game)
	if err != nil {
		return -1, fmt.Errorf("game handler error: %w", err)
	}
	Games[game.ID] = &game
	go RunState(game.Stage0, &game)
	return game.ID, nil
}

func GetGame(id int64) (*model.Game, error) {
	game, ok := Games[id]
	if !ok {
		return nil, fmt.Errorf("game not found")
	}
	return game, nil
}

func GetResult(id int64) ([]model.Player, error) {
	game, err := GetGame(id)
	if err != nil {
		return nil, err
	}
	if game.Stage0.ID == 0 {
		return nil, fmt.Errorf("game[%d] not started", id)
	}
	if game.Stage0.AvailablePlayer.ID != 0 {
		return nil, fmt.Errorf("game[%d] not finished", id)
	}
	return game.Result(), nil
}

func SendEvent(gameID, playerID, entryID int64) error {
	game, err := GetGame(gameID)
	if err != nil {
		return err
	}
	if game.Stage0.ID == 0 {
		return fmt.Errorf("game[%d] not started", gameID)
	}
	if game.Stage0.AvailablePlayer.ID != playerID {
		return fmt.Errorf("player[%d] not available", playerID)
	}
	for i := range game.Entries {
		if game.Entries[i].ID == entryID {
			if game.Entries[i].Banned || game.Entries[i].Picked {
				return fmt.Errorf("entry[%d] already used", entryID)
			}
			game.Stage0.Recv(&game.Entries[i])
			return nil
		}
	}
	return fmt.Errorf("entry[%d] not found", entryID)
}
