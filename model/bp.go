package model

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/jiny3/BanPickSys/pkg"
	"github.com/sirupsen/logrus"
)

type BP struct {
	ID        int64              `json:"id"`
	Stage0    *Stage             `json:"stage0"`
	Entries   []Entry            `json:"entries"`
	Players   map[string]*Player `json:"players"`
	Results   map[string]Player  `json:"results"`
	PlayerCap int                `json:"player_cap"`

	Send       []chan any `json:"-"`
	stageIndex int64
	mu         sync.Mutex
}

func NewBP(name string, entries []Entry) BP {
	gid := pkg.GenerateUUID(name)
	return BP{
		ID: gid,
		Stage0: NewStage(0, "初始化", name, START, time.Second*3, func(context.Context, chan any) {
			logrus.Debug("游戏初始化")
		}),
		Entries: entries,
		Players: make(map[string]*Player),
		Results: make(map[string]Player),
		Send:    []chan any{},
	}
}

func (g *BP) NewStageId() int64 {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.stageIndex++
	return g.stageIndex
}

func (g *BP) SetStage0(s *Stage) {
	g.Stage0 = s
}

func (g *BP) Join(player *Player, role string) error {
	if player == nil {
		return fmt.Errorf("非法id")
	}
	if len(g.Players) >= g.PlayerCap {
		return fmt.Errorf("玩家人数已满")
	}
	if p, ok := g.Players[role]; ok {
		if p.ID != player.ID {
			return fmt.Errorf("%s已被占用", role)
		}
		return nil
	}
	g.Players[role] = player
	if p, ok := g.Results[role]; ok {
		player.Picked = p.Picked
		player.Banned = p.Banned
		g.Results[role] = *player
	}
	return nil
}

func (g *BP) Leave(player *Player) error {
	if player == nil {
		return fmt.Errorf("非法id")
	}
	for role, p := range g.Players {
		if p.ID == player.ID {
			delete(g.Players, role)
			logrus.Debugf("玩家[%s]离开BP", player.Name)
			return nil
		}
	}
	return nil
}

func (g *BP) Result() map[string]Player {
	return g.Results
}
