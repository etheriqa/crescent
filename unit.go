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
	us             unitStatistics
	um             unitModification
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

type unitModification struct {
	armor                int32
	magicResistance      int32
	criticalStrikeChance int32
	criticalStrikeDamage int32
	cooldownReduction    int32
	threatFactor         int32
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

func (u *unit) attachOperator(o operator) {
	u.operators[o] = nil
	o.onAttach(u)
}

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

// updateModification updates the unitModification by iterating over operators
func (u *unit) updateModification() {
	for o := range u.operators {
		if m, ok := o.(modifier); ok {
			// todo sum up
			log.Debug(m)
		}
	}
	// todo update u.um
	u.notifyStats()
}
