package model

import (
	"sync"
	"time"

	"github.com/jiny3/BanPickSys/pkg"
	"github.com/sirupsen/logrus"
)

type Game struct {
	ID         int64    `json:"id"`
	Stage0     *Stage   `json:"stage0"`
	Entries    []Entry  `json:"entries"`
	Players    []Player `json:"players"`
	stageIndex int64
	mu         sync.Mutex `json:"-"`
}

func NewGame(name string, entries []Entry, players []Player) Game {
	gid := pkg.GenerateUUID(name)
	return Game{
		ID: gid,
		Stage0: NewStage(0, "初始化", name, &Player{ID: -1}, time.Second*3, func(chan struct{}, chan *Entry) {
			logrus.Debug("游戏初始化")
		}),
		Entries: entries,
		Players: players,
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

func (g *Game) Result() []Player {
	return g.Players
}
