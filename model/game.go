package model

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/jiny3/BanPickSys/pkg"
	"github.com/sirupsen/logrus"
)

type Game struct {
	ID        int64              `json:"id"`
	Stage0    *Stage             `json:"stage0"`
	Entries   []Entry            `json:"entries"`
	Players   map[string]*Player `json:"players"`
	PlayerCap int                `json:"player_cap"`

	Send       []chan *Game `json:"-"`
	stageIndex int64
	mu         sync.Mutex
}

func NewGame(name string, entries []Entry) Game {
	gid := pkg.GenerateUUID(name)
	return Game{
		ID: gid,
		Stage0: NewStage(0, "初始化", name, START, time.Second*3, func(context.Context, chan any) {
			logrus.Debug("游戏初始化")
		}),
		Entries: entries,
		Players: make(map[string]*Player),
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

func (g *Game) Join(player Player, role string) error {
	if len(g.Players) >= g.PlayerCap {
		return fmt.Errorf("玩家人数已满")
	}
	if p, ok := g.Players[role]; ok {
		return fmt.Errorf("角色已被占用")
	} else if p.ID == player.ID {
		return fmt.Errorf("玩家已在游戏中")
	}
	g.Players[role] = &player
	return nil
}

func (g *Game) Leave(p Player) error {
	for role, player := range g.Players {
		if player.ID == p.ID {
			delete(g.Players, role)
			return nil
		}
	}
	return fmt.Errorf("玩家不存在")
}

func (g *Game) Result() map[string]*Player {
	return g.Players
}
