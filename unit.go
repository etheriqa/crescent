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
	class      *Class
	resource   UnitResource
	correction UnitCorrection
	*Game
	*EventDispatcher
}

type UnitResource struct {
	Health Statistic
	Mana   Statistic
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
func NewUnit(game *Game, class *Class) *Unit {
	return &Unit{
		class:           class,
		resource:        UnitResource{},
		correction:      UnitCorrection{},
		Game:            game,
		EventDispatcher: NewEventDispatcher(),
	}
}

func (u *Unit) IsAlive() bool {
	return u.resource.Health > 0
}

func (u *Unit) IsDead() bool {
	return u.resource.Health <= 0
}

func (u *Unit) Health() Statistic {
	return u.resource.Health
}

func (u *Unit) HealthMax() Statistic {
	return u.class.Health
}

func (u *Unit) HealthRegeneration() Statistic {
	return u.class.HealthRegeneration
}

func (u *Unit) Mana() Statistic {
	return u.resource.Mana
}

func (u *Unit) ManaMax() Statistic {
	return u.class.Mana
}

func (u *Unit) ManaRegeneration() Statistic {
	return u.class.ManaRegeneration
}

func (u *Unit) Armor() Statistic {
	return u.class.Armor + u.correction.Armor
}

func (u *Unit) MagicResistance() Statistic {
	return u.class.MagicResistance + u.correction.MagicResistance
}

func (u *Unit) PhysicalDamageReductionFactor() Statistic {
	return damageReductionFactor(u.Armor())
}

func (u *Unit) MagicDamageReductionFactor() Statistic {
	return damageReductionFactor(u.MagicResistance())
}

func (u *Unit) CriticalStrikeChance() Statistic {
	return u.class.CriticalStrikeChance + u.correction.CriticalStrikeChance
}

func (u *Unit) CriticalStrikeFactor() Statistic {
	return u.class.CriticalStrikeFactor + u.correction.CriticalStrikeFactor
}

func (u *Unit) CooldownReduction() Statistic {
	return u.class.CooldownReduction + u.correction.CooldownReduction
}

func (u *Unit) DamageThreatFactor() Statistic {
	return u.class.DamageThreatFactor + u.correction.DamageThreatFactor
}

func (u *Unit) HealingThreatFactor() Statistic {
	return u.class.HealingThreatFactor + u.correction.HealingThreatFactor
}

// Ability returns the ability
func (u *Unit) Ability(key string) *Ability {
	return u.class.Ability(key)
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

// ModifyHealth modifies the unit health and returns before/after health
func (u *Unit) ModifyHealth(delta Statistic) (before, after Statistic, err error) {
	if u.IsDead() {
		return u.Health(), u.Health(), errors.New("Cannot modify the health of dead unit")
	}
	before = u.Health()
	after = u.Health() + delta
	if after < 0 {
		after = 0
	}
	if after > u.HealthMax() {
		after = u.HealthMax()
	}
	u.resource.Health = after
	u.Publish(message{}) // TODO pack message
	if delta < 0 {
		switch {
		case u.IsAlive():
			u.TriggerEvent(EventResourceDecreased)
		case u.IsDead():
			u.TriggerEvent(EventDead)
		}
	}
	return
}

// ModifyMana modifies the unit mana and returns before/after mana
func (u *Unit) ModifyMana(delta Statistic) (before, after Statistic, err error) {
	if u.IsDead() {
		return u.Health(), u.Health(), errors.New("Cannot modify the mana of dead unit")
	}
	before = u.Mana()
	after = u.Mana() + delta
	if after < 0 {
		after = 0
	}
	if after > u.ManaMax() {
		after = u.ManaMax()
	}
	u.resource.Mana = after
	u.Publish(message{}) // TODO pack message
	if delta < 0 {
		u.TriggerEvent(EventResourceDecreased)
	}
	return
}

// GameTick triggers onComplete iff the handler is completed
func (u *Unit) GameTick() {
	if u.IsDead() {
		return
	}
	u.TriggerEvent(EventGameTick)
}

// TickerTick performs regeneration and triggers eventTicker
func (u *Unit) TickerTick() {
	if u.IsDead() {
		return
	}
	u.performHealthRegeneration()
	u.performManaRegeneration()
	u.TriggerEvent(EventTicker)
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

// performHealthRegeneration performs health regeneration
func (u *Unit) performHealthRegeneration() {
	_, _, err := u.ModifyHealth(u.HealthRegeneration())
	if err != nil {
		log.WithFields(logrus.Fields{
			"err": err,
		}).Panic("Failed unit.performHealthRegeneration")
	}
}

// performManaRegeneration performs mana regeneration
func (u *Unit) performManaRegeneration() {
	_, _, err := u.ModifyMana(u.ManaRegeneration())
	if err != nil {
		log.WithFields(logrus.Fields{
			"err": err,
		}).Panic("Failed unit.performManaRegeneration")
	}
}
