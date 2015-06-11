package main

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
		cooldownDuration:   2 * second,
		disableTypes: []disableType{
			disableTypeStun,
		},
		perform: func(performer, receiver *unit) {
			receiver.attachOperator(newModifier(
				receiver,
				unitModification{
					magicResistance: -15,
				},
				q.name,
				1,
				12*second,
			))
			// TODO handle the error
			newPhysicalDamage(performer, receiver, 110).perform(performer.game)
			receiver.attachOperator(newDoT(
				newPhysicalDamage(performer, receiver, 25),
				q,
				4*second,
			))
		},
	}
	// Magic damage / Silence
	w = &ability{
		name:               "W",
		targetType:         targetTypeEnemy,
		healthCost:         0,
		manaCost:           40,
		activationDuration: 2 * second,
		cooldownDuration:   8 * second,
		disableTypes: []disableType{
			disableTypeStun,
			disableTypeSilence,
		},
		perform: func(performer, receiver *unit) {
			newMagicDamage(performer, receiver, 220).perform(performer.game)
			receiver.attachOperator(newDisable(
				receiver,
				disableTypeSilence,
				500*millisecond,
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
		cooldownDuration:   20 * second,
		disableTypes: []disableType{
			disableTypeStun,
		},
		perform: func(performer, receiver *unit) {
			newPhysicalDamage(performer, receiver, 280).perform(performer.game)
			receiver.attachOperator(newDisable(
				receiver,
				disableTypeStun,
				2*second,
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
		cooldownDuration:   60 * second,
		disableTypes: []disableType{
			disableTypeStun,
		},
		perform: func(performer, receiver *unit) {
			for _, friend := range performer.game.friends(performer) {
				friend.attachOperator(newModifier(
					friend,
					unitModification{
						criticalStrikeChance: 0.2,
						criticalStrikeFactor: 0.5,
					},
					r.name,
					1,
					10*second,
				))
			}
		},
	}
	class.abilities = []*ability{q, w, e, r}
	return class
}
