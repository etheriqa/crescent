package main

import (
	"time"
)

func newClassDisabler() *class {
	var q, w, e, r *ability
	class := &class{
		name: "Disabler",
		// TODO stats
		health:               800,
		healthRegeneration:   2,
		mana:                 300,
		manaRegeneration:     4,
		armor:                defaultArmor,
		magicResistance:      defaultMagicResistance,
		criticalStrikeChance: defaultCriticalStrikeChance,
		criticalStrikeFactor: defaultCriticalStrikeFactor,
		cooldownReduction:    defaultCooldownReduction,
		damageThreatFactor:   defaultDamageThreatFactor,
		healingThreatFactor:  defaultHealingThreatFactor,
	}
	// Physical damage / DoT / Magic resistance reduction
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
			receiver.attachOperator(newModifier(
				receiver,
				12*time.Second,
				unitModification{
					magicResistance: -15,
				},
				q,
				1,
			))
			// TODO handle the error
			newPhysicalDamage(performer, receiver, 110, q.name).perform(performer.game)
			receiver.attachOperator(newDoT(
				newPhysicalDamage(performer, receiver, 25, q.name),
				// TODO converter
				gameTime(4*time.Second),
			))
		},
	}
	// Magic damage / Silence
	w = &ability{
		name:               "W",
		targetType:         targetTypeEnemy,
		healthCost:         0,
		manaCost:           40,
		activationDuration: 2 * time.Second,
		cooldownDuration:   8 * time.Second,
		disableTypes: []disableType{
			disableTypeStun,
			disableTypeSilence,
		},
		perform: func(performer, receiver *unit) {
			newMagicDamage(performer, receiver, 220, w.name).perform(performer.game)
			receiver.attachOperator(newDisable(
				performer,
				receiver,
				disableTypeSilence,
				// TODO converter
				gameTime(500*time.Millisecond),
			))
		},
	}
	// Physical damage / Stun
	e = &ability{
		name:               "E",
		targetType:         targetTypeEnemy,
		healthCost:         0,
		manaCost:           60,
		activationDuration: 0,
		cooldownDuration:   20 * time.Second,
		disableTypes: []disableType{
			disableTypeStun,
		},
		perform: func(performer, receiver *unit) {
			newPhysicalDamage(performer, receiver, 280, e.name).perform(performer.game)
			receiver.attachOperator(newDisable(
				performer,
				receiver,
				disableTypeStun,
				// TODO converter
				gameTime(2*time.Second),
			))
		},
	}
	// Increasing critical / All
	r = &ability{
		name:               "R",
		targetType:         targetTypeNone,
		healthCost:         0,
		manaCost:           120,
		activationDuration: 0,
		cooldownDuration:   60 * time.Second,
		disableTypes: []disableType{
			disableTypeStun,
		},
		perform: func(performer, receiver *unit) {
			for _, friend := range performer.game.friends(performer) {
				friend.attachOperator(newModifier(
					friend,
					10*time.Second,
					unitModification{
						criticalStrikeChance: 0.2,
						criticalStrikeFactor: 0.5,
					},
					r,
					1,
				))
			}
		},
	}
	class.abilities = []*ability{q, w, e, r}
	return class
}
