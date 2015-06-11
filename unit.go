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

type Unit struct {
	id           unitID
	playerName   string
	unitName     string
	group        uint8
	seat         uint8
	class        *class
	resource     unitResource
	modification unitModification
	*Game
	*EventDispatcher
}

type unitResource struct {
	health Statistic
	mana   Statistic
}

type unitModification struct {
	armor                Statistic
	magicResistance      Statistic
	criticalStrikeChance Statistic
	criticalStrikeFactor Statistic
	cooldownReduction    Statistic
	damageThreatFactor   Statistic
	healingThreatFactor  Statistic
}

// newUnit initializes a unit
func NewUnit(g *Game, c *class) *Unit {
	return &Unit{
		class:           c,
		resource:        unitResource{},
		modification:    unitModification{},
		Game:            g,
		EventDispatcher: NewEventDispatcher(),
	}
}

func (u *Unit) isAlive() bool {
	return u.resource.health > 0
}

func (u *Unit) isDead() bool {
	return u.resource.health <= 0
}

func (u *Unit) health() Statistic {
	return u.resource.health
}

func (u *Unit) healthMax() Statistic {
	return u.class.health
}

func (u *Unit) healthRegeneration() Statistic {
	return u.class.healthRegeneration
}

func (u *Unit) mana() Statistic {
	return u.resource.mana
}

func (u *Unit) manaMax() Statistic {
	return u.class.mana
}

func (u *Unit) manaRegeneration() Statistic {
	return u.class.manaRegeneration
}

func (u *Unit) armor() Statistic {
	return u.class.armor + u.modification.armor
}

func (u *Unit) magicResistance() Statistic {
	return u.class.magicResistance + u.modification.magicResistance
}

func (u *Unit) physicalDamageReductionFactor() Statistic {
	return damageReductionFactor(u.armor())
}

func (u *Unit) magicDamageReductionFactor() Statistic {
	return damageReductionFactor(u.magicResistance())
}

func (u *Unit) criticalStrikeChance() Statistic {
	return u.class.criticalStrikeChance + u.modification.criticalStrikeChance
}

func (u *Unit) criticalStrikeFactor() Statistic {
	return u.class.criticalStrikeFactor + u.modification.criticalStrikeFactor
}

func (u *Unit) cooldownReduction() Statistic {
	return u.class.cooldownReduction + u.modification.cooldownReduction
}

func (u *Unit) damageThreatFactor() Statistic {
	return u.class.damageThreatFactor + u.modification.damageThreatFactor
}

func (u *Unit) healingThreatFactor() Statistic {
	return u.class.healingThreatFactor + u.modification.healingThreatFactor
}

// Friends returns the friend units
func (u *Unit) Friends() []*Unit {
	return u.Game.Friends(u)
}

// Enemies returns the enemy units
func (u *Unit) Enemies() []*Unit {
	return u.Game.Enemies(u)
}

// modifyHealth modifies the unit health and returns before/after health
func (u *Unit) modifyHealth(delta Statistic) (before, after Statistic, err error) {
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
			u.TriggerEvent(EventResourceDecreased)
		case u.isDead():
			u.TriggerEvent(EventDead)
		}
	}
	return
}

// modifyMana modifies the unit mana and returns before/after mana
func (u *Unit) modifyMana(delta Statistic) (before, after Statistic, err error) {
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
		u.TriggerEvent(EventResourceDecreased)
	}
	return
}

// GameTick triggers onComplete iff the handler is completed
func (u *Unit) GameTick() {
	if u.isDead() {
		return
	}
	u.TriggerEvent(EventGameTick)
}

// XoTTick performs regeneration and triggers evnentXoT
func (u *Unit) XoTTick() {
	if u.isDead() {
		return
	}
	u.performHealthRegeneration()
	u.performManaRegeneration()
	u.TriggerEvent(EventXoT)
}

// performHealthRegeneration performs health regeneration
func (u *Unit) performHealthRegeneration() {
	_, _, _, err := NewPureHealing(nil, u, u.healthRegeneration()).Perform()
	if err != nil {
		log.WithFields(logrus.Fields{
			"err": err,
		}).Fatal("Failed unit.performHealthRegeneration")
	}
}

// performManaRegeneration performs mana regeneration
func (u *Unit) performManaRegeneration() {
	err := u.performManaModification(u.manaRegeneration())
	if err != nil {
		log.WithFields(logrus.Fields{
			"err": err,
		}).Fatal("Failed unit.performManaRegeneration")
	}
}

// performManaModification performs mana modification
func (u *Unit) performManaModification(delta Statistic) error {
	_, _, err := u.modifyMana(delta)
	if err != nil {
		return err
	}
	u.Publish(message{
	// TODO pack message
	})
	return nil
}

// updateModification updates the unitModification
func (u *Unit) updateModification() {
	u.modification = unitModification{}
	u.ForSubjectHandler(u, func(ha Handler) {
		switch ha := ha.(type) {
		case *Modifier:
			u.modification.armor += ha.armor
			u.modification.magicResistance += ha.magicResistance
			u.modification.criticalStrikeChance += ha.criticalStrikeChance
			u.modification.criticalStrikeFactor += ha.criticalStrikeFactor
			u.modification.cooldownReduction += ha.cooldownReduction
			u.modification.damageThreatFactor += ha.damageThreatFactor
			u.modification.healingThreatFactor += ha.healingThreatFactor
		}
	})
	u.Publish(message{
	// TODO pack message
	})
}
