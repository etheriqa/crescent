package main

const (
	groupPlayer = iota
	groupEnemy
)

type uidType uint64

type unit struct {
	id         uidType
	playerName string
	unitName   string
	groupType  uint8
	seat       uint8
	abilities  map[string]*ability
	base       stats
	state      stats
	operators  map[*operator]interface{}
}

type stats struct {
	health               int32
	healthRegeneration   int32
	mana                 int32
	manaRegeneration     int32
	armor                int32
	magicResistance      int32
	criticalStrikeChance int32
	criticalStrikeDamage int32
	cooldownReduction    int32
	threatFactor         int32
}

func newPlayerUnit(id uidType, playerName, unitName string, seat uint8, s stats) *unit {
	return &unit{
		id:         id,
		playerName: playerName,
		unitName:   unitName,
		groupType:  groupPlayer,
		seat:       seat,
		abilities:  map[string]*ability{},
		base:       s,
		state:      s,
		operators:  make(map[*operator]interface{}),
	}
}
