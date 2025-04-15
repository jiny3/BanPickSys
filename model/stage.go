package model

import "time"

type Stage struct {
	ID              int64   `json:"id"`
	Name            string  `json:"name"`
	Description     string  `json:"description"`
	AvailablePlayer *Player `json:"available_player"`
	handler         func(chan struct{}, chan *Entry)
	recv            chan *Entry
	Duration        time.Duration `json:"duration"`
	StartTime       time.Time     `json:"start"`
	previous        []*Stage
	Nexts           []*Stage `json:"-"`
	wakeup          chan struct{}
	WaitingCounter  int64 `json:"-"`
}

func NewStage(id int64, name, description string, availablePlayer *Player, duration time.Duration, handler func(chan struct{}, chan *Entry)) *Stage {
	return &Stage{
		ID:              id,
		Name:            name,
		Description:     description,
		AvailablePlayer: availablePlayer,
		handler:         handler,
		recv:            make(chan *Entry),
		Duration:        duration,
		previous:        []*Stage{},
		wakeup:          make(chan struct{}, 10),
		Nexts:           []*Stage{},
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
	s.StartTime = time.Now()
	go s.handler(done, s.recv)
}

func (s *Stage) Recv(en *Entry) {
	s.recv <- en
}

func (s *Stage) Wakeup() {
	s.wakeup <- struct{}{}
}

func (s Stage) GetUntilTime() time.Duration {
	return time.Until(s.StartTime.Add(s.Duration))
}
