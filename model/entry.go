package model

type Entry struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Banned bool   `json:"banned"`
	Picked bool   `json:"picked"`
	Args   []Arg  `json:"args"`
}

type Arg struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (e *Entry) Ban() {
	e.Banned = true
}

func (e *Entry) Pick() {
	e.Picked = true
}
