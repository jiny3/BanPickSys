package service

import (
	"fmt"

	"github.com/jiny3/BanPickSys/model"
)

var Games = make(map[int64]*model.Game)

type GameInitFunc func(game *model.Game) error

func NewGame(name string, f GameInitFunc) (int64, error) {
	game := model.NewGame(name, GetEntries(name))
	err := f(&game)
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

func GetResult(id int64) (map[string]*model.Player, error) {
	game, err := GetGame(id)
	if err != nil {
		return nil, err
	}
	if game.Stage0.ID == 0 {
		return nil, fmt.Errorf("game[%d] not started", id)
	}
	if game.Stage0.Role != model.END {
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
	// TODO: 这里需要判断是否是当前玩家
	if game.Stage0.Role == "" {
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
