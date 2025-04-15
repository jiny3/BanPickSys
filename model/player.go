package model

type Player struct {
	ID     int64   `json:"id"`
	Name   string  `json:"name"`
	Banned []Entry `json:"banned"`
	Picked []Entry `json:"picked"`
}

func (p *Player) Ban(entry *Entry) {
	entry.Ban()
	p.Banned = append(p.Banned, *entry)
}

func (p *Player) Pick(entry *Entry) {
	entry.Pick()
	p.Picked = append(p.Picked, *entry)
}
