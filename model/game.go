package model

import (
	"fmt"
	"sync"
	"time"

	"slices"

	"github.com/jiny3/BanPickSys/pkg"
	"github.com/sirupsen/logrus"
)

type Game struct {
	ID        int64    `json:"id"`
	Stage0    *Stage   `json:"stage0"`
	Entries   []Entry  `json:"entries"`
	Players   []Player `json:"players"`
	PlayerCap int      `json:"player_cap"`

	Send       []chan *Game `json:"-"`
	stageIndex int64
	mu         sync.Mutex
}

func NewGame(name string, entries []Entry) Game {
	gid := pkg.GenerateUUID(name)
	return Game{
		ID: gid,
		Stage0: NewStage(0, "初始化", name, &Player{ID: -1}, time.Second*3, func(chan struct{}, chan *Entry) {
			logrus.Debug("游戏初始化")
		}),
		Entries: entries,
		Players: []Player{},
		Send:    []chan *Game{},
	}
}

func (g *Game) NewStageId() int64 {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.stageIndex++
	return g.stageIndex
}

func (g *Game) SetStage0(s *Stage) {
	g.Stage0 = s
}

func (g *Game) Join(p Player) error {
	if len(g.Players) >= g.PlayerCap {
		return fmt.Errorf("玩家人数已满")
	}
	g.Players = append(g.Players, p)
	return nil
}

func (g *Game) Leave(p Player) error {
	for i, player := range g.Players {
		if player.ID == p.ID {
			g.Players = slices.Delete(g.Players, i, i+1)
			return nil
		}
	}
	return fmt.Errorf("玩家不存在")
}

func (g *Game) Result() []Player {
	return g.Players
}
