package main

const (
	groupPlayer = iota
	groupEnemy
)

type uidType uint64

type unit struct {
	id             uidType
	playerName     string
	unitName       string
	group          uint8
	seat           uint8
	baseStats      unitStatistics
	currentStats   unitStatistics
	operators      map[operator]interface{}
	statsSubject   subject
	disableSubject subject
}

type unitStatistics struct {
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

func (u *unit) attachOperator(o operator) {
	u.operators[o] = nil
	o.onAttach()
}

func (u *unit) detachOperator(o operator) {
	delete(u.operators, o)
	o.onDetach()
}

func (u *unit) attachStatsObserver(o observer) { u.statsSubject.attach(o) }
func (u *unit) detachStatsObserver(o observer) { u.statsSubject.detach(o) }
func (u *unit) notifyStats()                   { u.statsSubject.notify() }

func (u *unit) attachDisableObserver(o observer) { u.disableSubject.attach(o) }
func (u *unit) detachDisableObserver(o observer) { u.disableSubject.detach(o) }
func (u *unit) notifyDisable()                   { u.disableSubject.notify() }
