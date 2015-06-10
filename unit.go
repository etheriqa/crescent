package main

import (
	"errors"

	"github.com/Sirupsen/logrus"
)

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

func (u *unit) healthMax() statistic {
	return u.class.health
}

func (u *unit) healthRegeneration() statistic {
	return u.class.healthRegeneration
}

func (u *unit) mana() statistic {
	return u.resource.mana
}

func (u *unit) manaMax() statistic {
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

func (u *unit) physicalDamageReductionFactor() statistic {
	return damageReductionFactor(u.armor())
}

func (u *unit) magicDamageReductionFactor() statistic {
	return damageReductionFactor(u.magicResistance())
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

// modifyHealth modifies the unit health and returns before/after health
func (u *unit) modifyHealth(delta statistic) (before, after statistic, err error) {
	if u.isDead() {
		return u.health(), u.health(), errors.New("Cannot modify the health of dead unit")
	}
	before = u.health()
	after = u.health() + delta
	if after < 0 {
		after = 0
	}
	if after > u.healthMax() {
		after = u.healthMax()
	}
	u.resource.health = after
	if delta < 0 {
		switch {
		case u.isAlive():
			u.triggerEvent(eventResourceDecreased)
		case u.isDead():
			u.triggerEvent(eventDead)
		}
	}
	return
}

// modifyMana modifies the unit mana and returns before/after mana
func (u *unit) modifyMana(delta statistic) (before, after statistic, err error) {
	if u.isDead() {
		return u.health(), u.health(), errors.New("Cannot modify the mana of dead unit")
	}
	before = u.mana()
	after = u.mana() + delta
	if after < 0 {
		after = 0
	}
	if after > u.manaMax() {
		after = u.manaMax()
	}
	u.resource.mana = after
	if delta < 0 {
		u.triggerEvent(eventResourceDecreased)
	}
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
	_, _, _, err := newPureHealing(nil, u, u.healthRegeneration(), "").perform(u.game)
	if err != nil {
		log.WithFields(logrus.Fields{
			"err": err,
		}).Fatal("Failed unit.performHealthRegeneration")
	}
}

// performManaRegeneration performs mana regeneration
func (u *unit) performManaRegeneration() {
	err := u.performManaModification(u.manaRegeneration())
	if err != nil {
		log.WithFields(logrus.Fields{
			"err": err,
		}).Fatal("Failed unit.performManaRegeneration")
	}
}

// performManaModification performs mana modification
func (u *unit) performManaModification(delta statistic) error {
	_, _, err := u.modifyMana(delta)
	if err != nil {
		return err
	}
	u.game.publish(message{
	// TODO pack message
	})
	return nil
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
