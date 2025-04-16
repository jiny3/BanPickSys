package model

import (
	"context"
	"time"
)

const (
	// 角色
	START = "start" // 初始化阶段
	END   = "end"   // 结束阶段
)

type Stage struct {
	ID          int64         `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Role        string        `json:"role"`
	Duration    time.Duration `json:"duration"`
	Start       time.Time     `json:"start"`

	Nexts          []*Stage `json:"-"`
	WaitingCounter int64    `json:"-"`
	handler        StageHandler
	recv           chan any
	previous       []*Stage
	wakeup         chan struct{}
}

// 阶段处理函数, ctx用于超时退出等, recv用于接收消息
type StageHandler = func(ctx context.Context, recv chan any)

func NewStage(id int64, name, description string, role string, duration time.Duration, handler StageHandler) *Stage {
	return &Stage{
		ID:          id,
		Name:        name,
		Description: description,
		Role:        role,
		Duration:    duration,

		Nexts:    []*Stage{},
		handler:  handler,
		recv:     make(chan any),
		previous: []*Stage{},
		wakeup:   make(chan struct{}, 10),
	}
}

func (s *Stage) Link(prevs ...*Stage) {
	s.previous = prevs
	for _, prev := range prevs {
		prev.Nexts = append(prev.Nexts, s)
	}
}

func (s *Stage) Waiting() {
	for s.WaitingCounter < int64(len(s.previous)) {
		<-s.wakeup
		s.WaitingCounter++
	}
}

func (s *Stage) Handle(ctx context.Context, cancel context.CancelFunc) {
	defer cancel()
	s.Start = time.Now()
	s.handler(ctx, s.recv)
}

func (s *Stage) Recv(en *Entry) {
	s.recv <- en
}

func (s *Stage) Wakeup() {
	s.wakeup <- struct{}{}
}
