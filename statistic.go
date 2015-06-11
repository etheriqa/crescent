package main

import (
	"math/rand"
)

type statistic float64

type damage struct {
	performer            *unit
	receiver             *unit
	amount               statistic
	criticalStrikeChance statistic
	criticalStrikeFactor statistic
}

type healing struct {
	performer            *unit
	receiver             *unit
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
func diceCritical(performer *unit) bool {
	return rand.Float64() < float64(performer.criticalStrikeChance())
}

// newPhysicalDamage returns a damage affected by armor of the receiver
func newPhysicalDamage(performer, receiver *unit, baseDamage statistic) *damage {
	return newTrueDamage(
		performer,
		receiver,
		baseDamage*receiver.physicalDamageReductionFactor(),
	)
}

// newMagicDamage returns a damage affected by magic resistance of the receiver
func newMagicDamage(performer, receiver *unit, baseDamage statistic) *damage {
	return newTrueDamage(
		performer,
		receiver,
		baseDamage*receiver.magicDamageReductionFactor(),
	)
}

// newTrueDamage returns a damage that ignores damage reduction
func newTrueDamage(performer, receiver *unit, baseDamage statistic) *damage {
	return &damage{
		performer:            performer,
		receiver:             receiver,
		amount:               baseDamage,
		criticalStrikeChance: performer.criticalStrikeChance(),
		criticalStrikeFactor: performer.criticalStrikeFactor(),
	}
}

// newPureDamage returns a damage that ignores both damage reduction and critical strike
func newPureDamage(performer, receiver *unit, baseDamage statistic) *damage {
	return &damage{
		performer:            performer,
		receiver:             receiver,
		amount:               baseDamage,
		criticalStrikeChance: 0,
		criticalStrikeFactor: 0,
	}
}

// perform subtracts amount the damage from the receiver and attaches a threat handler to the performer and publishes a message
func (d damage) perform(g *game) (before, after statistic, crit bool, err error) {
	amount, crit := applyCriticalStrike(
		d.amount,
		d.criticalStrikeChance,
		d.criticalStrikeFactor,
	)
	after, before, err = d.receiver.modifyHealth(-amount)
	if err != nil {
		return
	}
	if d.performer != nil {
		d.receiver.attachHandler(newDamageThreat(d.performer, d.receiver, d.amount))
	}
	g.publish(message{
	// TODO pack message
	})
	return
}

// newHealing returns a healing
func newHealing(performer, receiver *unit, baseHealing statistic) *healing {
	return &healing{
		performer:            performer,
		receiver:             receiver,
		amount:               baseHealing,
		criticalStrikeChance: performer.criticalStrikeChance(),
		criticalStrikeFactor: performer.criticalStrikeFactor(),
	}
}

// newPureHealing returns a healing that ignores critical strike
func newPureHealing(performer, receiver *unit, baseHealing statistic) *healing {
	return &healing{
		performer:            performer,
		receiver:             receiver,
		amount:               baseHealing,
		criticalStrikeChance: 0,
		criticalStrikeFactor: 0,
	}
}

// perform adds amount of healing to the receiver and attaches a threat handler to the enemies and publish a message
func (h healing) perform(g *game) (after, before statistic, crit bool, err error) {
	amount, crit := applyCriticalStrike(
		h.amount,
		h.criticalStrikeChance,
		h.criticalStrikeFactor,
	)
	after, before, err = h.receiver.modifyHealth(amount)
	if err != nil {
		return
	}
	if h.performer != nil {
		for _, enemy := range g.enemies(h.performer) {
			enemy.attachHandler(newHealingThreat(h.performer, enemy, h.amount))
		}
	}
	g.publish(message{
	// TODO pack message
	})
	return
}
