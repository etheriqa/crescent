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
		perform: func(subject, object *unit) {
			object.AttachHandler(NewModifier(
				object,
				unitModification{
					magicResistance: -15,
				},
				q.name,
				1,
				12*second,
			))
			// TODO handle the error
			newPhysicalDamage(subject, object, 110).perform(subject.game)
			object.AttachHandler(NewDoT(
				newPhysicalDamage(subject, object, 25),
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
		perform: func(subject, object *unit) {
			newMagicDamage(subject, object, 220).perform(subject.game)
			object.AttachHandler(NewDisable(
				object,
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
		perform: func(subject, object *unit) {
			newPhysicalDamage(subject, object, 280).perform(subject.game)
			object.AttachHandler(NewDisable(
				object,
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
		perform: func(subject, object *unit) {
			for _, friend := range subject.game.friends(subject) {
				friend.AttachHandler(NewModifier(
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
