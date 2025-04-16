package service

import (
	"fmt"

	"github.com/jiny3/BanPickSys/model"
	"github.com/sirupsen/logrus"
)

var BPs = make(map[int64]*model.BP)

type GameInitFunc func(game *model.BP) error

func NewBP(name string, f GameInitFunc) (int64, error) {
	game := model.NewBP(name, model.GetEntries(name))
	if f == nil {
		logrus.Errorf("game handler is nil")
		return -1, fmt.Errorf("miss function")
	}
	err := f(&game)
	if err != nil {
		return -1, fmt.Errorf("game handler error: %w", err)
	}
	BPs[game.ID] = &game
	go RunState(game.Stage0, &game)
	return game.ID, nil
}

func GetBP(id int64) (*model.BP, error) {
	game, ok := BPs[id]
	if !ok {
		return nil, fmt.Errorf("game not found")
	}
	return game, nil
}

func GetEntries(id int64) ([]model.Entry, error) {
	game, err := GetBP(id)
	if err != nil {
		return nil, err
	}
	return game.Entries, nil
}

func GetResult(id int64) (map[string]model.Player, error) {
	game, err := GetBP(id)
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

func Join(gameID, playerID int64, role string) error {
	game, err := GetBP(gameID)
	if err != nil {
		return err
	}
	for i := range game.Players {
		if game.Players[i].ID == playerID {
			return fmt.Errorf("player[%d] already joined", playerID)
		}
	}
	err = game.Join(Players[playerID], role)
	if err != nil {
		return err
	}
	return nil
}

func Leave(gameID, playerID int64) error {
	game, err := GetBP(gameID)
	if err != nil {
		return err
	}
	err = game.Leave(Players[playerID])
	if err != nil {
		return err
	}
	return nil
}

func SendEvent(gameID, playerID, entryID int64) error {
	game, err := GetBP(gameID)
	if err != nil {
		return err
	}
	if game.Stage0.ID == 0 {
		return fmt.Errorf("game[%d] not started", gameID)
	}
	if game.Players[game.Stage0.Role] != Players[playerID] {
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
