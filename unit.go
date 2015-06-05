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
	group      uint8
	seat       uint8
	hp         int32
	mp         int32
	us         *unitStatistics
	um         *unitModification
	operators  map[operator]bool
	dispatcher *eventDispatcher
}

// newUnit initializes a unit
func newUnit() *unit {
	return &unit{
		us:         &unitStatistics{},
		um:         &unitModification{},
		operators:  make(map[operator]bool),
		dispatcher: newEventDispatcher(),
	}
}

func (u *unit) isAlive() bool {
	return u.hp > 0
}

func (u *unit) isDead() bool {
	return u.hp <= 0
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
	u.operators[o] = true
	o.onAttach(u)
}

// detachOperator removes the operator
func (u *unit) detachOperator(o operator) {
	delete(u.operators, o)
	o.onDetach(u)
}

// addEventHandler registers the eventHandler
func (u *unit) addEventHandler(h eventHandler, e event) {
	u.dispatcher.addEventHandler(h, e)
}

// addEventHandler unregisters the eventHandler
func (u *unit) removeEventHandler(h eventHandler, e event) {
	u.dispatcher.removeEventHandler(h, e)
}

// triggerEvent triggers the event
func (u *unit) triggerEvent(e event) {
	u.dispatcher.triggerEvent(e)
}

// gameTick triggers onComplete iff the operator is completed
func (u *unit) gameTick(out chan message) {
	if u.isDead() {
		return
	}
	u.triggerEvent(eventGameTick)
}

// statsTick performs regeneration and triggers statsTick
func (u *unit) statsTick(out chan message) {
	if u.isDead() {
		return
	}
	u.performHealthRegeneration(out)
	u.performManaRegeneration(out)
	u.triggerEvent(eventStatsTick)
}

// performHealthRegeneration performs health regeneration
func (u *unit) performHealthRegeneration(out chan message) {
	reg := u.healthRegeneration() + u.hp - u.health()
	if reg < 0 {
		return
	}
	u.hp += reg
	out <- message{
		// todo pack message
		t: outHealthReg,
	}
}

// performManaRegeneration performs mana regeneration
func (u *unit) performManaRegeneration(out chan message) {
	reg := u.manaRegeneration() + u.mp - u.mana()
	if reg < 0 {
		return
	}
	u.mp += reg
	out <- message{
		// todo pack message
		t: outManaReg,
	}
}

// updateModification updates the unitModification
func (u *unit) updateModification() {
	u.um = &unitModification{}
	for o := range u.operators {
		if _, ok := o.(*modifier); ok {
			u.um.add(o.(*modifier).um)
		}
	}
	u.triggerEvent(eventStats)
}
