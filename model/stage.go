package model

import "time"

type Stage struct {
	ID          int64         `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Player0     *Player       `json:"available_player"`
	Duration    time.Duration `json:"duration"`
	Start       time.Time     `json:"start"`

	Nexts          []*Stage `json:"-"`
	WaitingCounter int64    `json:"-"`
	handler        func(chan struct{}, chan *Entry)
	recv           chan *Entry
	previous       []*Stage
	wakeup         chan struct{}
}

func NewStage(id int64, name, description string, Player0 *Player, duration time.Duration, handler func(chan struct{}, chan *Entry)) *Stage {
	return &Stage{
		ID:          id,
		Name:        name,
		Description: description,
		Player0:     Player0,
		Duration:    duration,

		Nexts:    []*Stage{},
		handler:  handler,
		recv:     make(chan *Entry),
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

func (s *Stage) Handle(done chan struct{}) {
	s.Start = time.Now()
	go s.handler(done, s.recv)
}

func (s *Stage) Recv(en *Entry) {
	s.recv <- en
}

func (s *Stage) Wakeup() {
	s.wakeup <- struct{}{}
}

func (s Stage) GetUntilTime() time.Duration {
	return time.Until(s.Start.Add(s.Duration))
}
