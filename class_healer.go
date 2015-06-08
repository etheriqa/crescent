package main

import (
	"time"
)

func newClassHealer() *class {
	var q, w, e, r *ability
	class := &class{
		name: "Healer",
		// TODO stats
		health:               700,
		healthRegeneration:   2,
		mana:                 400,
		manaRegeneration:     6,
		armor:                defaultArmor,
		magicResistance:      defaultMagicResistance,
		criticalStrikeChance: defaultCriticalStrikeChance,
		criticalStrikeFactor: defaultCriticalStrikeFactor,
		cooldownReduction:    defaultCooldownReduction,
		damageThreatFactor:   defaultDamageThreatFactor,
		healingThreatFactor:  defaultHealingThreatFactor,
	}
	// Magic damage / Mana restore
	q = &ability{
		name:               "Healer Q",
		targetType:         targetTypeEnemy,
		healthCost:         0,
		manaCost:           0,
		activationDuration: 2 * time.Second,
		cooldownDuration:   2 * time.Second,
		disableTypes: []disableType{
			disableTypeSilence,
			disableTypeStun,
		},
		perform: func(performer, receiver *unit) {
			// TODO handle the error
			before, after, _ := newMagicDamage(performer, receiver, 100, q.name).perform(performer.game)
			// TODO send a message including the ability name
			performer.performManaModification((before - after) * 0.1)
		},
	}
	// HoT
	w = &ability{
		name:               "Healer W",
		targetType:         targetTypeFriend,
		healthCost:         0,
		manaCost:           40,
		activationDuration: 2 * time.Second,
		cooldownDuration:   4 * time.Second,
		disableTypes: []disableType{
			disableTypeSilence,
			disableTypeStun,
		},
		perform: func(performer, receiver *unit) {
			receiver.attachOperator(newHoT(
				newHealing(performer, receiver, 20, w.name),
				// TODO converter
				gameTime(12*time.Second),
			))
		},
	}
	// Healing
	e = &ability{
		name:               "Healer E",
		targetType:         targetTypeFriend,
		healthCost:         0,
		manaCost:           80,
		activationDuration: 2 * time.Second,
		cooldownDuration:   8 * time.Second,
		disableTypes: []disableType{
			disableTypeSilence,
			disableTypeStun,
		},
		perform: func(performer, receiver *unit) {
			newHealing(performer, receiver, 400, e.name).perform(performer.game)
		},
	}
	// HoT / Increasing critical strike chance and critical strike factor
	r = &ability{
		name:               "Healer R",
		targetType:         targetTypeAllFriends,
		healthCost:         0,
		manaCost:           200,
		activationDuration: 0,
		cooldownDuration:   60 * time.Second,
		disableTypes: []disableType{
			disableTypeSilence,
			disableTypeStun,
		},
		perform: func(performer, receiver *unit) {
			performer.attachOperator(newModifier(performer, 6*time.Second, unitModification{
				criticalStrikeChance: 0.5,
				criticalStrikeFactor: 1.5,
			}))
			for _, friend := range performer.game.friends(performer) {
				friend.attachOperator(newHoT(
					newHealing(performer, friend, 20, r.name),
					// TODO converter
					gameTime(6*time.Second),
				))
			}
		},
	}
	return class
}
