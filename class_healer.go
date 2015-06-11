package main

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
		activationDuration: 2 * second,
		cooldownDuration:   2 * second,
		disableTypes: []disableType{
			disableTypeSilence,
			disableTypeStun,
		},
		perform: func(performer, receiver *unit) {
			// TODO handle the error
			before, after, _, _ := newMagicDamage(performer, receiver, 100).perform(performer.game)
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
		activationDuration: 2 * second,
		cooldownDuration:   4 * second,
		disableTypes: []disableType{
			disableTypeSilence,
			disableTypeStun,
		},
		perform: func(performer, receiver *unit) {
			receiver.attachOperator(newHoT(
				newHealing(performer, receiver, 20),
				w,
				12*second,
			))
		},
	}
	// Healing
	e = &ability{
		name:               "Healer E",
		targetType:         targetTypeFriend,
		healthCost:         0,
		manaCost:           80,
		activationDuration: 2 * second,
		cooldownDuration:   8 * second,
		disableTypes: []disableType{
			disableTypeSilence,
			disableTypeStun,
		},
		perform: func(performer, receiver *unit) {
			newHealing(performer, receiver, 400).perform(performer.game)
		},
	}
	// HoT / Increasing critical strike chance and critical strike factor
	r = &ability{
		name:               "Healer R",
		targetType:         targetTypeNone,
		healthCost:         0,
		manaCost:           200,
		activationDuration: 0,
		cooldownDuration:   60 * second,
		disableTypes: []disableType{
			disableTypeSilence,
			disableTypeStun,
		},
		perform: func(performer, receiver *unit) {
			performer.attachOperator(newModifier(
				performer,
				unitModification{
					criticalStrikeChance: 0.5,
					criticalStrikeFactor: 1.5,
				},
				r,
				1,
				6*second,
			))
			for _, friend := range performer.game.friends(performer) {
				friend.attachOperator(newHoT(
					newHealing(performer, friend, 20),
					r,
					6*second,
				))
			}
		},
	}
	class.abilities = []*ability{q, w, e, r}
	return class
}
