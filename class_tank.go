package main

import (
	"time"
)

func newClassTank() *class {
	var q, w, e, r *ability
	class := &class{
		name: "Tank",
		// TODO stats
		health:               1000,
		healthRegeneration:   5,
		mana:                 200,
		manaRegeneration:     3,
		armor:                defaultArmor,
		magicResistance:      defaultMagicResistance,
		criticalStrikeChance: defaultCriticalStrikeChance,
		criticalStrikeFactor: defaultCriticalStrikeFactor,
		cooldownReduction:    defaultCooldownReduction,
		damageThreatFactor:   defaultDamageThreatFactor,
		healingThreatFactor:  defaultHealingThreatFactor,
	}
	// True damage / Increasing threat factor
	q = &ability{
		name:               "Tank Q",
		targetType:         targetTypeEnemy,
		healthCost:         0,
		manaCost:           0,
		activationDuration: 0,
		cooldownDuration:   2 * time.Second,
		disableTypes: []disableType{
			disableTypeStun,
		},
		perform: func(performer, receiver *unit) {
			performer.attachOperator(newModifier(
				performer,
				10*time.Second,
				unitModification{
					damageThreatFactor: 4,
				},
			))
			// TODO handle the error
			newTrueDamage(performer, receiver, 120, q.name).perform(performer.game)
		},
	}
	// Physical damage / Increasing AR & MR
	w = &ability{
		name:               "Tank W",
		targetType:         targetTypeEnemy,
		healthCost:         0,
		manaCost:           15,
		activationDuration: 2 * time.Second,
		cooldownDuration:   8 * time.Second,
		disableTypes: []disableType{
			disableTypeStun,
		},
		perform: func(performer, receiver *unit) {
			performer.attachOperator(newModifier(
				performer,
				2*time.Second,
				unitModification{
					armor:           50,
					magicResistance: 50,
				},
			))
			// TODO handle the error
			newPhysicalDamage(performer, receiver, 200, w.name).perform(performer.game)
		},
	}
	// Physical damage / Life steal
	e = &ability{
		name:               "Tank E",
		targetType:         targetTypeEnemy,
		healthCost:         0,
		manaCost:           50,
		activationDuration: 0,
		cooldownDuration:   15 * time.Second,
		disableTypes: []disableType{
			disableTypeStun,
		},
		perform: func(performer, receiver *unit) {
			// TODO handle the error
			before, after, _ := newPhysicalDamage(performer, receiver, 300, e.name).perform(performer.game)
			// TODO handle the error
			newPureHealing(performer, receiver, (before-after)*0.6, e.name).perform(performer.game)
		},
	}
	// Increasing AR & MR
	r = &ability{
		name:               "Tank R",
		targetType:         targetTypeSelf,
		healthCost:         0,
		manaCost:           120,
		activationDuration: 4,
		cooldownDuration:   60,
		disableTypes: []disableType{
			disableTypeStun,
			disableTypeSilence,
		},
		perform: func(performer, receiver *unit) {
			performer.attachOperator(newModifier(performer, 8*time.Second, unitModification{
				armor:           150,
				magicResistance: 150,
			}))
		},
	}
	class.abilities = []*ability{q, w, e, r}
	return class
}
