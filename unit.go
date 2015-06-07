package main

const (
	groupPlayer = iota
	groupEnemy
)

type unitID uint64

type unit struct {
	id         unitID
	playerName string
	unitName   string
	group      uint8
	seat       uint8
	ur         *unitResource
	us         *unitStatistics
	um         *unitModification
	operators  map[operator]bool
	dispatcher *eventDispatcher
	game       *game
}

type unitResource struct {
	health int32
	mana   int32
}

// newUnit initializes a unit
func newUnit(g *game) *unit {
	return &unit{
		us:         &unitStatistics{},
		um:         &unitModification{},
		operators:  make(map[operator]bool),
		dispatcher: newEventDispatcher(),
		game:       g,
	}
}

func (u *unit) now() gameTime {
	return u.game.now()
}

func (u *unit) publish(m message) {
	u.game.publish(m)
}

func (u *unit) isAlive() bool {
	return u.ur.health > 0
}

func (u *unit) isDead() bool {
	return u.ur.health <= 0
}

func (u *unit) health() int32 {
	return u.ur.health
}

func (u *unit) maxHealth() int32 {
	return u.us.health
}

func (u *unit) healthRegeneration() int32 {
	return u.us.healthRegeneration
}

func (u *unit) mana() int32 {
	return u.ur.mana
}

func (u *unit) maxMana() int32 {
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

func (u *unit) criticalStrikeFactor() int32 {
	return u.us.criticalStrikeFactor + u.um.criticalStrikeFactor
}

func (u *unit) cooldownReduction() int32 {
	return u.us.cooldownReduction + u.um.cooldownReduction
}

func (u *unit) damageThreatFactor() int32 {
	return u.us.damageThreatFactor + u.um.damageThreatFactor
}

func (u *unit) healingThreatFactor() int32 {
	return u.us.healingThreatFactor + u.um.healingThreatFactor
}

// addHealth adds health and returns before/after health
func (u *unit) addHealth(d int32) (before, after int32) {
	before = u.health()
	after = u.health() + d
	if after < 0 {
		after = 0
	}
	if after > u.maxHealth() {
		after = u.maxHealth()
	}
	u.ur.health = after
	return
}

// addMana adds mana and returns before/after mana
func (u *unit) addMana(d int32) (before, after int32) {
	before = u.mana()
	after = u.mana() + d
	if after < 0 {
		after = 0
	}
	if after > u.maxMana() {
		after = u.maxMana()
	}
	u.ur.mana = after
	return
}

// attachOperator adds the operator
func (u *unit) attachOperator(o operator) {
	u.operators[o] = true
	o.onAttach()
}

// detachOperator removes the operator
func (u *unit) detachOperator(o operator) {
	delete(u.operators, o)
	o.onDetach()
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
func (u *unit) gameTick() {
	if u.isDead() {
		return
	}
	u.triggerEvent(eventGameTick)
}

// xotTick performs regeneration and triggers evnentXoT
func (u *unit) xotTick() {
	if u.isDead() {
		return
	}
	u.performHealthRegeneration()
	u.performManaRegeneration()
	u.triggerEvent(eventXoT)
}

// performHealthRegeneration performs health regeneration
func (u *unit) performHealthRegeneration() {
	u.addHealth(u.healthRegeneration())
	u.publish(message{
		// TODO pack message
		t: outHealthReg,
	})
}

// performManaRegeneration performs mana regeneration
func (u *unit) performManaRegeneration() {
	u.addMana(u.manaRegeneration())
	u.publish(message{
		// TODO pack message
		t: outManaReg,
	})
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
