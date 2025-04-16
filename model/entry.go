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

var InitEntries = []Entry{
	{
		ID:   1,
		Name: "绷带海豹",
		Args: []Arg{
			{
				Key:   "attack",
				Value: "5",
			},
			{
				Key:   "life",
				Value: "35",
			},
			{
				Key:   "skill",
				Value: "生命值低于25时，攻击力+10",
			},
			{
				Key:   "evolve",
				Value: "50%概率免疫致命伤害, 并将生命值回复至24",
			},
		},
	},
	{
		ID:   2,
		Name: "天使海豹",
		Args: []Arg{
			{
				Key:   "attack",
				Value: "6",
			},
			{
				Key:   "life",
				Value: "45",
			},
			{
				Key:   "skill",
				Value: "行动中碰撞到的己方海豹回复天使海豹攻击力相应的生命值",
			},
			{
				Key:   "evolve",
				Value: "行动中碰撞到的首个小海豹受到的回复效果翻倍",
			},
		},
	},
	{
		ID:   3,
		Name: "刺头海豹",
		Args: []Arg{
			{
				Key:   "attack",
				Value: "6",
			},
			{
				Key:   "life",
				Value: "35",
			},
			{
				Key:   "skill",
				Value: "发射速度降低40%, 被碰撞时对敌方造成刺头海豹攻击力相应的反击伤害",
			},
			{
				Key:   "evolve",
				Value: "反击伤害提升至攻击力的200%",
			},
		},
	},
	{
		ID:   4,
		Name: "橡胶海豹",
		Args: []Arg{
			{
				Key:   "attack",
				Value: "5",
			},
			{
				Key:   "life",
				Value: "40",
			},
			{
				Key:   "skill",
				Value: "行动中, 每次碰撞后攻击力+2, 持续至本回合结束",
			},
			{
				Key:   "evolve",
				Value: "发射速度提升60%",
			},
		},
	},
	{
		ID:   5,
		Name: "大力海豹",
		Args: []Arg{
			{
				Key:   "attack",
				Value: "6",
			},
			{
				Key:   "life",
				Value: "35",
			},
			{
				Key:   "skill",
				Value: "行动中碰撞到的己方海豹攻击力+3, 持续至本回合结束",
			},
			{
				Key:   "evolve",
				Value: "行动中碰撞到的首个小海豹获得额外行动",
			},
		},
	},
	{
		ID:   6,
		Name: "泡泡",
		Args: []Arg{
			{
				Key:   "attack",
				Value: "6",
			},
			{
				Key:   "life",
				Value: "40",
			},
			{
				Key:   "skill",
				Value: "行动中, 击败敌方海豹可获得额外行动, 该效果无法连续触发",
			},
			{
				Key:   "evolve",
				Value: "在额外行动中, 攻击力+5",
			},
		},
	},
	{
		ID:   7,
		Name: "红温海豹",
		Args: []Arg{
			{
				Key:   "attack",
				Value: "8",
			},
			{
				Key:   "life",
				Value: "35",
			},
			{
				Key:   "skill",
				Value: "行动中碰撞到首个敌方海豹时造成200%攻击力的伤害, 并立刻停下",
			},
			{
				Key:   "evolve",
				Value: "行动中首次造成的碰撞伤害提升至攻击力的300%",
			},
		},
	},
	{
		ID:   8,
		Name: "普通灰豹",
		Args: []Arg{
			{
				Key:   "attack",
				Value: "8",
			},
			{
				Key:   "life",
				Value: "40",
			},
			{
				Key:   "skill",
				Value: "它很普通",
			},
			{
				Key:   "evolve",
				Value: "攻击力+4",
			},
		},
	},
	{
		ID:   9,
		Name: "嘭嘭海豹",
		Args: []Arg{
			{
				Key:   "attack",
				Value: "6",
			},
			{
				Key:   "life",
				Value: "40",
			},
			{
				Key:   "skill",
				Value: "行动中, 发射2个水弹攻击随机敌方海豹, 每个造成4点伤害",
			},
			{
				Key:   "evolve",
				Value: "水弹伤害提升至10点",
			},
		},
	},
	{
		ID:   10,
		Name: "幽灵海豹",
		Args: []Arg{
			{
				Key:   "attack",
				Value: "4",
			},
			{
				Key:   "life",
				Value: "40",
			},
			{
				Key:   "skill",
				Value: "每轮开始时在场地中生成3个鬼火, 每次吃掉鬼火可增加3点攻击力",
			},
			{
				Key:   "evolve",
				Value: "每次吃掉鬼火可增加的攻击力提升至6点",
			},
		},
	},
	{
		ID:   11,
		Name: "黑化海豹",
		Args: []Arg{
			{
				Key:   "attack",
				Value: "1",
			},
			{
				Key:   "life",
				Value: "1",
			},
			{
				Key:   "skill",
				Value: "比赛开始时吞噬己方海豹, 获得其150%的攻击力与70%的生命值. 黑化海豹被击败后将吐出友方海豹, 且无法复活",
			},
			{
				Key:   "evolve",
				Value: "吐出的友方海豹攻击力+3",
			},
		},
	},
	{
		ID:   12,
		Name: "双色海豹",
		Args: []Arg{
			{
				Key:   "attack",
				Value: "6",
			},
			{
				Key:   "life",
				Value: "35",
			},
			{
				Key:   "skill",
				Value: "红色部分造成200%攻击力的伤害, 蓝色部分受到50%伤害",
			},
			{
				Key:   "evolve",
				Value: "红色部分造成的伤害提升至300%",
			},
		},
	},
}

func GetEntries(name string) []Entry {
	// TODO: 数据库读取entries
	_entries := make([]Entry, len(InitEntries))
	copy(_entries, InitEntries)
	return _entries
}
