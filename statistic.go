package main

import (
	"math/rand"
)

type statistic float64

// reductionFactor calculates a reduction factor
func reductionFactor(reduction statistic) statistic {
	return statistic(1 / (1 + float64(reduction)/100))
}

// diceCritical dices whether critical strike is happening or not
func diceCritical(performer *unit) bool {
	return rand.Float64() < float64(performer.criticalStrikeChance())
}

// dicePhysicalDamage calculates amount of physical damage
func dicePhysicalDamage(performer, receiver *unit, baseDamage statistic) (damage statistic) {
	damage = diceTrueDamage(performer, receiver, baseDamage)
	damage *= reductionFactor(receiver.armor())
	return
}

// diceMagicDamage calculates amount of magic damage
func diceMagicDamage(performer, receiver *unit, baseDamage statistic) (damage statistic) {
	damage = diceTrueDamage(performer, receiver, baseDamage)
	damage *= reductionFactor(receiver.armor())
	return
}

// diceTrueDamage calculates amount of true damage
func diceTrueDamage(performer, receiver *unit, baseDamage statistic) statistic {
	damage := baseDamage
	if diceCritical(performer) {
		damage += damage * performer.criticalStrikeFactor()
	}
	return damage
}

// diceHealing calculates amount of healing
func diceHealing(performer, receiver *unit, baseHealing statistic) statistic {
	healing := baseHealing
	if diceCritical(performer) {
		healing += healing * performer.criticalStrikeFactor()
	}
	return healing
}
