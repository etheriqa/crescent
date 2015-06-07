package main

const (
	groupPlayer = iota
	groupEnemy
)

type unitID uint64

type unit struct {
	id           unitID
	playerName   string
	unitName     string
	group        uint8
	seat         uint8
	class        *class
	resource     unitResource
	modification unitModification
	operators    map[operator]bool
	dispatcher   *eventDispatcher
	game         *game
}

type unitResource struct {
	health statistic
	mana   statistic
}

type unitModification struct {
	armor                statistic
	magicResistance      statistic
	criticalStrikeChance statistic
	criticalStrikeFactor statistic
	cooldownReduction    statistic
	damageThreatFactor   statistic
	healingThreatFactor  statistic
}

// newUnit initializes a unit
func newUnit(g *game, c *class) *unit {
	return &unit{
		class:        c,
		resource:     unitResource{},
		modification: unitModification{},
		operators:    make(map[operator]bool),
		dispatcher:   newEventDispatcher(),
		game:         g,
	}
}

func (u *unit) now() gameTime {
	return u.game.now()
}

func (u *unit) publish(m message) {
	u.game.publish(m)
}

func (u *unit) isAlive() bool {
	return u.resource.health > 0
}

func (u *unit) isDead() bool {
	return u.resource.health <= 0
}

func (u *unit) health() statistic {
	return u.resource.health
}

func (u *unit) maxHealth() statistic {
	return u.class.health
}

func (u *unit) healthRegeneration() statistic {
	return u.class.healthRegeneration
}

func (u *unit) mana() statistic {
	return u.resource.mana
}

func (u *unit) maxMana() statistic {
	return u.class.mana
}

func (u *unit) manaRegeneration() statistic {
	return u.class.manaRegeneration
}

func (u *unit) armor() statistic {
	return u.class.armor + u.modification.armor
}

func (u *unit) magicResistance() statistic {
	return u.class.magicResistance + u.modification.magicResistance
}

func (u *unit) criticalStrikeChance() statistic {
	return u.class.criticalStrikeChance + u.modification.criticalStrikeChance
}

func (u *unit) criticalStrikeFactor() statistic {
	return u.class.criticalStrikeFactor + u.modification.criticalStrikeFactor
}

func (u *unit) cooldownReduction() statistic {
	return u.class.cooldownReduction + u.modification.cooldownReduction
}

func (u *unit) damageThreatFactor() statistic {
	return u.class.damageThreatFactor + u.modification.damageThreatFactor
}

func (u *unit) healingThreatFactor() statistic {
	return u.class.healingThreatFactor + u.modification.healingThreatFactor
}

// addHealth adds health and returns before/after health
func (u *unit) addHealth(delta statistic) (before, after statistic) {
	before = u.health()
	after = u.health() + delta
	if after < 0 {
		after = 0
	}
	if after > u.maxHealth() {
		after = u.maxHealth()
	}
	u.resource.health = after
	return
}

// addMana adds mana and returns before/after mana
func (u *unit) addMana(delta statistic) (before, after statistic) {
	before = u.mana()
	after = u.mana() + delta
	if after < 0 {
		after = 0
	}
	if after > u.maxMana() {
		after = u.maxMana()
	}
	u.resource.mana = after
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
	u.modification = unitModification{}
	for o := range u.operators {
		switch o := o.(type) {
		case *modifier:
			u.modification.armor += o.armor
			u.modification.magicResistance += o.magicResistance
			u.modification.criticalStrikeChance += o.criticalStrikeChance
			u.modification.criticalStrikeFactor += o.criticalStrikeFactor
			u.modification.cooldownReduction += o.cooldownReduction
			u.modification.damageThreatFactor += o.damageThreatFactor
			u.modification.healingThreatFactor += o.healingThreatFactor
		}
	}
}
