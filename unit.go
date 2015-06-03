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
	us             *unitStatistics
	um             *unitModification
	operators      map[operator]interface{}
	statsSubject   *subject
	disableSubject *subject
}

// newUnit initializes a unit
func newUnit() *unit {
	return &unit{
		us:             &unitStatistics{},
		um:             &unitModification{},
		operators:      make(map[operator]interface{}),
		statsSubject:   newSubject(),
		disableSubject: newSubject(),
	}
}

func (u *unit) isAlive() bool {
	return u.health() > 0
}

func (u *unit) isDead() bool {
	return u.health() <= 0
}

func (u *unit) health() int32 {
	return u.us.health
}

func (u *unit) healthRegeneration() int32 {
	return u.us.healthRegeneration
}

func (u *unit) mana() int32 {
	return u.us.mana
}

func (u *unit) manaRegeneration() int32 {
	return u.us.manaRegeneration
}

func (u *unit) armor() int32 {
	return u.us.armor + u.um.armor
}

func (u *unit) magicResistance() int32 {
	return u.us.magicResistance + u.um.magicResistance
}

func (u *unit) criticalStrikeChance() int32 {
	return u.us.criticalStrikeChance + u.um.criticalStrikeChance
}

func (u *unit) criticalStrikeDamage() int32 {
	return u.us.criticalStrikeDamage + u.um.criticalStrikeDamage
}

func (u *unit) cooldownReduction() int32 {
	return u.us.cooldownReduction + u.um.cooldownReduction
}

func (u *unit) threatFactor() int32 {
	return u.us.threatFactor + u.um.threatFactor
}

// attachOperator adds the operator
func (u *unit) attachOperator(o operator) {
	u.operators[o] = nil
	o.onAttach(u)
}

// detachOperator removes the operator
func (u *unit) detachOperator(o operator) {
	delete(u.operators, o)
	o.onDetach(u)
}

func (u *unit) attachStatsObserver(o observer) { u.statsSubject.attach(o) }
func (u *unit) detachStatsObserver(o observer) { u.statsSubject.detach(o) }
func (u *unit) notifyStats()                   { u.statsSubject.notify() }

func (u *unit) attachDisableObserver(o observer) { u.disableSubject.attach(o) }
func (u *unit) detachDisableObserver(o observer) { u.disableSubject.detach(o) }
func (u *unit) notifyDisable()                   { u.disableSubject.notify() }

// progress triggers onComplete iff the operator is completed
func (u *unit) progress(out chan message) {
	if u.isDead() {
		return
	}
	for o := range u.operators {
		if !o.isComplete(u) {
			continue
		}
		o.onComplete(u)
		u.detachOperator(o)
	}
}

// tick triggers onTick and performs regeneration
func (u *unit) tick(out chan message) {
	if u.isDead() {
		return
	}
	// todo perform health regeneration
	out <- message{
		// todo pack message
		t: outHealthReg,
	}
	// todo perform mana regeneration
	out <- message{
		// todo pack message
		t: outManaReg,
	}
	for o := range u.operators {
		o.onTick(u)
	}
}

// updateModification updates the unitModification
func (u *unit) updateModification() {
	u.um = &unitModification{}
	for o := range u.operators {
		if m, ok := o.(*modifier); ok {
			u.um.add(m.um)
		}
	}
	u.notifyStats()
}
