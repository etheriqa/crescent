package main

import (
	"time"
)

type targetType uint8

const (
	_ targetType = iota
	targetTypeOneself
	targetTypeFriend
	targetTypeAllFriends
	targetTypeEnemy
	targetTypeAllEnemies
)

type ability struct {
	name               string
	targetType         targetType
	healthCost         statistic
	manaCost           statistic
	activationDuration time.Duration
	cooldownDuration   time.Duration
	disableTypes       []disableType
	perform            func(performer, receiver *unit)
}

// satisfiedRequirements returns true iff the ability satisfy activation requirements
func (a *ability) satisfiedRequirements(performer *unit) bool {
	// TODO check receiver
	if performer.health() < a.healthCost {
		return false
	}
	if performer.mana() < a.manaCost {
		return false
	}
	for o := range performer.operators {
		switch o := o.(type) {
		case *cooldown:
			if a == o.ability {
				return false
			}
		case *disable:
			for d := range a.disableTypes {
				if disableType(d) == o.disableType {
					return false
				}
			}
		}
	}
	return true
}
