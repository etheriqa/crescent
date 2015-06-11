package main

import (
	"math/rand"
)

type Statistic float64

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
