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
		perform: func(subject, object *unit) {
			subject.AttachHandler(NewModifier(
				subject,
				unitModification{
					damageThreatFactor: 1,
				},
				q.name,
				5,
				10*second,
			))
			// TODO handle the error
			newTrueDamage(subject, object, 120).perform(subject.game)
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
		perform: func(subject, object *unit) {
			subject.AttachHandler(NewModifier(
				subject,
				unitModification{
					armor:           50,
					magicResistance: 50,
				},
				w.name,
				1,
				2*second,
			))
			// TODO handle the error
			newPhysicalDamage(subject, object, 200).perform(subject.game)
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
		perform: func(subject, object *unit) {
			// TODO handle the error
			before, after, _, _ := newPhysicalDamage(subject, object, 300).perform(subject.game)
			// TODO handle the error
			newPureHealing(subject, object, (before-after)*0.6).perform(subject.game)
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
		perform: func(subject, object *unit) {
			subject.AttachHandler(NewModifier(
				subject,
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
