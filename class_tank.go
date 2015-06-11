package main

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
		cooldownDuration:   2 * second,
		disableTypes: []disableType{
			disableTypeStun,
		},
		perform: func(performer, receiver *unit) {
			performer.attachOperator(newModifier(
				performer,
				unitModification{
					damageThreatFactor: 1,
				},
				q.name,
				5,
				10*second,
			))
			// TODO handle the error
			newTrueDamage(performer, receiver, 120).perform(performer.game)
		},
	}
	// Physical damage / Increasing AR & MR
	w = &ability{
		name:               "Tank W",
		targetType:         targetTypeEnemy,
		healthCost:         0,
		manaCost:           15,
		activationDuration: 2 * second,
		cooldownDuration:   8 * second,
		disableTypes: []disableType{
			disableTypeStun,
		},
		perform: func(performer, receiver *unit) {
			performer.attachOperator(newModifier(
				performer,
				unitModification{
					armor:           50,
					magicResistance: 50,
				},
				w.name,
				1,
				2*second,
			))
			// TODO handle the error
			newPhysicalDamage(performer, receiver, 200).perform(performer.game)
		},
	}
	// Physical damage / Life steal
	e = &ability{
		name:               "Tank E",
		targetType:         targetTypeEnemy,
		healthCost:         0,
		manaCost:           50,
		activationDuration: 0,
		cooldownDuration:   15 * second,
		disableTypes: []disableType{
			disableTypeStun,
		},
		perform: func(performer, receiver *unit) {
			// TODO handle the error
			before, after, _, _ := newPhysicalDamage(performer, receiver, 300).perform(performer.game)
			// TODO handle the error
			newPureHealing(performer, receiver, (before-after)*0.6).perform(performer.game)
		},
	}
	// Increasing AR & MR
	r = &ability{
		name:               "Tank R",
		targetType:         targetTypeNone,
		healthCost:         0,
		manaCost:           120,
		activationDuration: 4,
		cooldownDuration:   60,
		disableTypes: []disableType{
			disableTypeStun,
			disableTypeSilence,
		},
		perform: func(performer, receiver *unit) {
			performer.attachOperator(newModifier(
				performer,
				unitModification{
					armor:           150,
					magicResistance: 150,
				},
				r.name,
				1,
				8*second,
			))
		},
	}
	class.abilities = []*ability{q, w, e, r}
	return class
}
