package main

import (
	"time"
)

func newClassAssassin() *class {
	type stack struct {
		partialOperator
		stack int
	}
	var q, w, e, r *ability
	class := &class{
		name: "Assassin",
		// TODO stats
		health:               600,
		healthRegeneration:   2,
		mana:                 200,
		manaRegeneration:     3,
		armor:                defaultArmor,
		magicResistance:      defaultMagicResistance,
		criticalStrikeChance: defaultCriticalStrikeChance + 0.05,
		criticalStrikeFactor: defaultCriticalStrikeFactor + 0.5,
		cooldownReduction:    defaultCooldownReduction,
		damageThreatFactor:   defaultDamageThreatFactor,
		healingThreatFactor:  defaultHealingThreatFactor,
	}
	// Physical damage
	q = &ability{
		name:               "Q",
		targetType:         targetTypeEnemy,
		healthCost:         0,
		manaCost:           0,
		activationDuration: 0,
		cooldownDuration:   2 * time.Second,
		disableTypes: []disableType{
			disableTypeStun,
		},
		perform: func(performer, receiver *unit) {
			newPhysicalDamage(performer, receiver, 140, q.name).perform(performer.game)
		},
	}
	// Physical damage / DoT / Increasing stacks
	w = &ability{
		name:               "W",
		targetType:         targetTypeEnemy,
		healthCost:         0,
		manaCost:           20,
		activationDuration: 0,
		cooldownDuration:   8 * time.Second,
		disableTypes: []disableType{
			disableTypeStun,
		},
		perform: func(performer, receiver *unit) {
			// TODO check crits and increase stack
			newPhysicalDamage(performer, receiver, 80, w.name).perform(performer.game)
			receiver.attachOperator(newDoT(
				newPhysicalDamage(performer, receiver, 20, w.name),
				// TODO converter
				gameTime(10*time.Second),
			))
		},
	}
	// Increasing stacks / Decreasing armor and magic resistance
	e = &ability{
		name:               "E",
		targetType:         targetTypeSelf,
		healthCost:         0,
		manaCost:           40,
		activationDuration: 0,
		cooldownDuration:   20 * time.Second,
		disableTypes: []disableType{
			disableTypeStun,
		},
		perform: func(performer, receiver *unit) {
			// TODO increase stack
			performer.attachOperator(newModifier(
				performer,
				8*time.Second,
				unitModification{
					armor:           -25,
					magicResistance: -25,
				},
				e,
				1,
			))
		},
	}
	// Physical
	r = &ability{
		name:               "R",
		targetType:         targetTypeEnemy,
		healthCost:         0,
		manaCost:           120,
		activationDuration: 0,
		cooldownDuration:   60 * time.Second,
		disableTypes: []disableType{
			disableTypeStun,
		},
		perform: func(performer, receiver *unit) {
			newPhysicalDamage(performer, receiver, 600, r.name).perform(performer.game)
			// TODO consume all stacks
		},
	}
	class.abilities = []*ability{q, w, e, r}
	return class
}
