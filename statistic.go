package main

import (
	"math/rand"
)

type Statistic float64

type Damage struct {
	subject              *Unit
	object               *Unit
	amount               Statistic
	criticalStrikeChance Statistic
	criticalStrikeFactor Statistic
}

type Healing struct {
	subject              *Unit
	object               *Unit
	amount               Statistic
	criticalStrikeChance Statistic
	criticalStrikeFactor Statistic
}

// damageReductionFactor calculates a damage reduction factor on armor or magic resistance
func damageReductionFactor(damageReduction Statistic) Statistic {
	return Statistic(1 / (1 + float64(damageReduction)/100))
}

// applyCriticalStrike judges whether critical strike happens or not and returns amount of damage / healing that affected by critical strike
func applyCriticalStrike(base, chance, factor Statistic) (amount Statistic, critical bool) {
	amount = base
	critical = rand.Float64() < float64(chance)
	if critical {
		amount += base * factor
	}
	return
}

// NewPhysicalDamage returns a damage affected by armor of the object
func NewPhysicalDamage(subject, object *Unit, baseDamage Statistic) *Damage {
	return NewTrueDamage(
		subject,
		object,
		baseDamage*object.physicalDamageReductionFactor(),
	)
}

// NewMagicDamage returns a damage affected by magic resistance of the object
func NewMagicDamage(subject, object *Unit, baseDamage Statistic) *Damage {
	return NewTrueDamage(
		subject,
		object,
		baseDamage*object.magicDamageReductionFactor(),
	)
}

// NewTrueDamage returns a damage that ignores damage reduction
func NewTrueDamage(subject, object *Unit, baseDamage Statistic) *Damage {
	return &Damage{
		subject:              subject,
		object:               object,
		amount:               baseDamage,
		criticalStrikeChance: subject.criticalStrikeChance(),
		criticalStrikeFactor: subject.criticalStrikeFactor(),
	}
}

// NewPureDamage returns a damage that ignores both damage reduction and critical strike
func NewPureDamage(subject, object *Unit, baseDamage Statistic) *Damage {
	return &Damage{
		subject:              subject,
		object:               object,
		amount:               baseDamage,
		criticalStrikeChance: 0,
		criticalStrikeFactor: 0,
	}
}

// Perform subtracts amount the damage from the object and attaches a threat handler to the subject and publishes a message
func (d *Damage) Perform() (before, after Statistic, crit bool, err error) {
	amount, crit := applyCriticalStrike(
		d.amount,
		d.criticalStrikeChance,
		d.criticalStrikeFactor,
	)
	after, before, err = d.object.modifyHealth(-amount)
	if err != nil {
		return
	}
	if d.subject != nil {
		d.object.AttachHandler(newDamageThreat(d.subject, d.object, d.amount))
	}
	d.object.Publish(message{
	// TODO pack message
	})
	return
}

// NewHealing returns a healing
func NewHealing(subject, object *Unit, baseHealing Statistic) *Healing {
	return &Healing{
		subject:              subject,
		object:               object,
		amount:               baseHealing,
		criticalStrikeChance: subject.criticalStrikeChance(),
		criticalStrikeFactor: subject.criticalStrikeFactor(),
	}
}

// NewPureHealing returns a healing that ignores critical strike
func NewPureHealing(subject, object *Unit, baseHealing Statistic) *Healing {
	return &Healing{
		subject:              subject,
		object:               object,
		amount:               baseHealing,
		criticalStrikeChance: 0,
		criticalStrikeFactor: 0,
	}
}

// Perform adds amount of healing to the object and attaches a threat handler to the enemies and publish a message
func (h *Healing) Perform() (after, before Statistic, crit bool, err error) {
	amount, crit := applyCriticalStrike(
		h.amount,
		h.criticalStrikeChance,
		h.criticalStrikeFactor,
	)
	after, before, err = h.object.modifyHealth(amount)
	if err != nil {
		return
	}
	if h.subject != nil {
		for _, enemy := range h.subject.Enemies() {
			enemy.AttachHandler(NewHealingThreat(h.subject, enemy, h.amount))
		}
	}
	h.object.Publish(message{
	// TODO pack message
	})
	return
}
