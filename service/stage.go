package service

import (
	"context"

	"github.com/jiny3/BanPickSys/model"
)

func LinkStages(cur *model.Stage, prevs ...*model.Stage) {
	cur.Link(prevs...)
}

func RunState(s *model.Stage, game *model.BP) {
	// wait for previous stages
	s.Waiting()

	// startup with timer
	game.SetStage0(s)
	ctx, cancel := context.WithTimeout(context.Background(), s.Duration)
	defer cancel()

	// send stage change signal to all channels
	for _, s := range game.Send {
		s <- game
	}
	s.Handle(ctx, cancel)

	// wakeup next stages
	for _, next := range s.Nexts {
		if next.WaitingCounter == 0 {
			go RunState(next, game)
		}
		next.Wakeup()
	}
}
