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
	id         unitID
	playerName string
	unitName   string
	group      uint8
	seat       uint8
	class      *class
	resource   unitResource
	correction UnitCorrection
	*Game
	*EventDispatcher
}

type unitResource struct {
	health Statistic
	mana   Statistic
}

type UnitCorrection struct {
	Armor                Statistic
	MagicResistance      Statistic
	CriticalStrikeChance Statistic
	CriticalStrikeFactor Statistic
	CooldownReduction    Statistic
	DamageThreatFactor   Statistic
	HealingThreatFactor  Statistic
}

// newUnit initializes a unit
func NewUnit(g *Game, c *class) *Unit {
	return &Unit{
		class:           c,
		resource:        unitResource{},
		correction:      UnitCorrection{},
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
	return u.class.armor + u.correction.Armor
}

func (u *Unit) magicResistance() Statistic {
	return u.class.magicResistance + u.correction.MagicResistance
}

func (u *Unit) physicalDamageReductionFactor() Statistic {
	return damageReductionFactor(u.armor())
}

func (u *Unit) magicDamageReductionFactor() Statistic {
	return damageReductionFactor(u.magicResistance())
}

func (u *Unit) criticalStrikeChance() Statistic {
	return u.class.criticalStrikeChance + u.correction.CriticalStrikeChance
}

func (u *Unit) criticalStrikeFactor() Statistic {
	return u.class.criticalStrikeFactor + u.correction.CriticalStrikeFactor
}

func (u *Unit) cooldownReduction() Statistic {
	return u.class.cooldownReduction + u.correction.CooldownReduction
}

func (u *Unit) damageThreatFactor() Statistic {
	return u.class.damageThreatFactor + u.correction.DamageThreatFactor
}

func (u *Unit) healingThreatFactor() Statistic {
	return u.class.healingThreatFactor + u.correction.HealingThreatFactor
}

// Friends returns the friend units
func (u *Unit) Friends() []*Unit {
	return u.Game.Friends(u)
}

// Enemies returns the enemy units
func (u *Unit) Enemies() []*Unit {
	return u.Game.Enemies(u)
}

// ForSubjectHandler calls the callback with the handler has this unit as subject
func (u *Unit) ForSubjectHandler(callback func(Handler)) {
	u.Game.ForSubjectHandler(u, callback)
}

// ForObjectHandler calls the callback with the handler has this unit as object
func (u *Unit) ForObjectHandler(callback func(Handler)) {
	u.Game.ForObjectHandler(u, callback)
}

// EverySubjectHandler returns true if all of callback results are true
func (u *Unit) EverySubjectHandler(callback func(Handler) bool) bool {
	return u.Game.EverySubjectHandler(u, callback)
}

// EveryObjectHandler returns true if all of callback results are true
func (u *Unit) EveryObjectHandler(callback func(Handler) bool) bool {
	return u.Game.EveryObjectHandler(u, callback)
}

// SomeSubjectHandler returns true if any of callback results are true
func (u *Unit) SomeSubjectHandler(callback func(Handler) bool) bool {
	return u.Game.SomeSubjectHandler(u, callback)
}

// SomeObjectHandler returns true if any of callback results are true
func (u *Unit) SomeObjectHandler(callback func(Handler) bool) bool {
	return u.Game.SomeObjectHandler(u, callback)
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

// TickerTick performs regeneration and triggers eventTicker
func (u *Unit) TickerTick() {
	if u.isDead() {
		return
	}
	u.performHealthRegeneration()
	u.performManaRegeneration()
	u.TriggerEvent(EventTicker)
}

// performHealthRegeneration performs health regeneration
func (u *Unit) performHealthRegeneration() {
	_, _, _, err := NewPureHealing(MakeObject(u), u.healthRegeneration()).Perform()
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

// ReloadCorrection updates the UnitCorrection
func (u *Unit) ReloadCorrection() {
	u.correction = UnitCorrection{}
	u.ForSubjectHandler(func(ha Handler) {
		switch ha := ha.(type) {
		case *Corrector:
			u.correction.Armor += ha.Armor()
			u.correction.MagicResistance += ha.MagicResistance()
			u.correction.CriticalStrikeChance += ha.CriticalStrikeChance()
			u.correction.CriticalStrikeFactor += ha.CriticalStrikeFactor()
			u.correction.CooldownReduction += ha.CooldownReduction()
			u.correction.DamageThreatFactor += ha.DamageThreatFactor()
			u.correction.HealingThreatFactor += ha.HealingThreatFactor()
		}
	})
	u.Publish(message{
	// TODO pack message
	})
}
