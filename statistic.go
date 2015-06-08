package main

import (
	"math/rand"
)

type statistic float64
type damage statistic
type healing statistic

// reductionFactor calculates a reduction factor
func reductionFactor(reduction statistic) statistic {
	return statistic(1 / (1 + float64(reduction)/100))
}

// diceCritical dices whether critical strike is happening or not
func diceCritical(performer *unit) bool {
	return rand.Float64() < float64(performer.criticalStrikeChance())
}

// dicePhysicalDamage calculates amount of physical damage
func dicePhysicalDamage(performer, receiver *unit, baseDamage damage) damage {
	d := diceTrueDamage(performer, receiver, baseDamage)
	d *= damage(reductionFactor(receiver.armor()))
	return d
}

// diceMagicDamage calculates amount of magic damage
func diceMagicDamage(performer, receiver *unit, baseDamage damage) damage {
	d := diceTrueDamage(performer, receiver, baseDamage)
	d *= damage(reductionFactor(receiver.armor()))
	return d
}

// diceTrueDamage calculates amount of true damage
func diceTrueDamage(performer, receiver *unit, baseDamage damage) damage {
	if diceCritical(performer) {
		return baseDamage * damage(1+performer.criticalStrikeFactor())
	} else {
		return baseDamage
	}
}

// diceHealing calculates amount of healing
func diceHealing(performer, receiver *unit, baseHealing healing) healing {
	if diceCritical(performer) {
		return baseHealing * healing(1+performer.criticalStrikeFactor())
	} else {
		return baseHealing
	}
}
