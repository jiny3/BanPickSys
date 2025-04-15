package service

import (
	"time"

	"github.com/jiny3/BanPickSys/model"
)

func LinkStages(cur *model.Stage, prevs ...*model.Stage) {
	cur.Link(prevs...)
}

func RunState(s *model.Stage, game *model.Game) {
	// wait for previous stages
	s.Waiting()

	// startup with timer
	game.SetStage0(s)
	done := make(chan struct{})
	defer close(done)

	s.Handle(done)
	select {
	case <-done:
	case <-time.After(time.Until(s.StartTime.Add(s.Duration))):
	}

	// wakeup next stages
	for _, next := range s.Nexts {
		if next.WaitingCounter == 0 {
			go RunState(next, game)
		}
		next.Wakeup()
	}
	done <- struct{}{}
}
