package main

import (
	"math/rand"
)

type statistic float64

type damage struct {
	subject              *unit
	object               *unit
	amount               statistic
	criticalStrikeChance statistic
	criticalStrikeFactor statistic
}

type healing struct {
	subject              *unit
	object               *unit
	amount               statistic
	criticalStrikeChance statistic
	criticalStrikeFactor statistic
}

// damageReductionFactor calculates a damage reduction factor on armor or magic resistance
func damageReductionFactor(damageReduction statistic) statistic {
	return statistic(1 / (1 + float64(damageReduction)/100))
}

// applyCriticalStrike judges whether critical strike happens or not and returns amount of damage / healing that affected by critical strike
func applyCriticalStrike(base, chance, factor statistic) (amount statistic, critical bool) {
	amount = base
	critical = rand.Float64() < float64(chance)
	if critical {
		amount += base * factor
	}
	return
}

// diceCritical dices whether critical strike is happening or not
func diceCritical(subject *unit) bool {
	return rand.Float64() < float64(subject.criticalStrikeChance())
}

// newPhysicalDamage returns a damage affected by armor of the object
func newPhysicalDamage(subject, object *unit, baseDamage statistic) *damage {
	return newTrueDamage(
		subject,
		object,
		baseDamage*object.physicalDamageReductionFactor(),
	)
}

// newMagicDamage returns a damage affected by magic resistance of the object
func newMagicDamage(subject, object *unit, baseDamage statistic) *damage {
	return newTrueDamage(
		subject,
		object,
		baseDamage*object.magicDamageReductionFactor(),
	)
}

// newTrueDamage returns a damage that ignores damage reduction
func newTrueDamage(subject, object *unit, baseDamage statistic) *damage {
	return &damage{
		subject:              subject,
		object:               object,
		amount:               baseDamage,
		criticalStrikeChance: subject.criticalStrikeChance(),
		criticalStrikeFactor: subject.criticalStrikeFactor(),
	}
}

// newPureDamage returns a damage that ignores both damage reduction and critical strike
func newPureDamage(subject, object *unit, baseDamage statistic) *damage {
	return &damage{
		subject:              subject,
		object:               object,
		amount:               baseDamage,
		criticalStrikeChance: 0,
		criticalStrikeFactor: 0,
	}
}

// perform subtracts amount the damage from the object and attaches a threat handler to the subject and publishes a message
func (d damage) perform(g *game) (before, after statistic, crit bool, err error) {
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
	g.publish(message{
	// TODO pack message
	})
	return
}

// newHealing returns a healing
func newHealing(subject, object *unit, baseHealing statistic) *healing {
	return &healing{
		subject:              subject,
		object:               object,
		amount:               baseHealing,
		criticalStrikeChance: subject.criticalStrikeChance(),
		criticalStrikeFactor: subject.criticalStrikeFactor(),
	}
}

// newPureHealing returns a healing that ignores critical strike
func newPureHealing(subject, object *unit, baseHealing statistic) *healing {
	return &healing{
		subject:              subject,
		object:               object,
		amount:               baseHealing,
		criticalStrikeChance: 0,
		criticalStrikeFactor: 0,
	}
}

// perform adds amount of healing to the object and attaches a threat handler to the enemies and publish a message
func (h healing) perform(g *game) (after, before statistic, crit bool, err error) {
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
		for _, enemy := range g.enemies(h.subject) {
			enemy.AttachHandler(newHealingThreat(h.subject, enemy, h.amount))
		}
	}
	g.publish(message{
	// TODO pack message
	})
	return
}
