package main

type targetType uint8

const (
	_ targetType = iota
	targetTypeOneself
	targetTypeFriend
	targetTypeAllFriends
	targetTypeEnemy
	targetTypeAllEnemies
)

type damageType uint8

const (
	_ damageType = iota
	damageTypePhysical
	damageTypeMagic
	damageTypeTrue
)

type ability struct {
	name               string
	targetType         targetType
	damageType         damageType
	healthCost         statistic
	manaCost           statistic
	activationDuration gameDuration
	cooldownDuration   gameDuration
	disableTypes       []disableType
	perform            func(performer, receiver *unit)
}

// satisfiedRequirements returns true iff the ability satisfy activation requirements
func (a *ability) satisfiedRequirements(performer *unit) bool {
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
