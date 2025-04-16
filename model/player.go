package model

import "github.com/jiny3/BanPickSys/pkg"

type Player struct {
	ID     int64   `json:"id"`
	Name   string  `json:"name"`
	Banned []Entry `json:"banned"`
	Picked []Entry `json:"picked"`
}

func NewPlayer(name string) *Player {
	return &Player{
		ID:     pkg.GenerateUUID(name),
		Name:   name,
		Banned: []Entry{},
		Picked: []Entry{},
	}
}

func (p *Player) Ban(entry *Entry) {
	entry.Ban()
	p.Banned = append(p.Banned, *entry)
}

func (p *Player) Pick(entry *Entry) {
	entry.Pick()
	p.Picked = append(p.Picked, *entry)
}
