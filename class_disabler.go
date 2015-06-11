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
		armor:                DefaultArmor,
		magicResistance:      DefaultMagicResistance,
		criticalStrikeChance: DefaultCriticalStrikeChance,
		criticalStrikeFactor: DefaultCriticalStrikeFactor,
		cooldownReduction:    DefaultCooldownReduction,
		damageThreatFactor:   DefaultDamageThreatFactor,
		healingThreatFactor:  DefaultHealingThreatFactor,
	}
	// Physical damage / DoT / Magic resistance reduction
	q = &ability{
		name:               "Q",
		TargetType:         TargetTypeEnemy,
		healthCost:         0,
		manaCost:           0,
		activationDuration: 0,
		cooldownDuration:   2 * Second,
		disableTypes: []DisableType{
			DisableTypeStun,
		},
		perform: func(subject, object *Unit) {
			object.AttachHandler(NewModifier(
				object,
				unitModification{
					magicResistance: -15,
				},
				q.name,
				1,
				12*Second,
			))
			// TODO handle the error
			NewPhysicalDamage(subject, object, 110).Perform()
			object.AttachHandler(NewDoT(
				NewPhysicalDamage(subject, object, 25),
				q,
				4*Second,
			))
		},
	}
	// Magic damage / Silence
	w = &ability{
		name:               "W",
		TargetType:         TargetTypeEnemy,
		healthCost:         0,
		manaCost:           40,
		activationDuration: 2 * Second,
		cooldownDuration:   8 * Second,
		disableTypes: []DisableType{
			DisableTypeStun,
			DisableTypeSilence,
		},
		perform: func(subject, object *Unit) {
			NewMagicDamage(subject, object, 220).Perform()
			object.AttachHandler(NewDisable(
				object,
				DisableTypeSilence,
				500*Millisecond,
			))
		},
	}
	// Physical damage / Stun
	e = &ability{
		name:               "E",
		TargetType:         TargetTypeEnemy,
		healthCost:         0,
		manaCost:           60,
		activationDuration: 0,
		cooldownDuration:   20 * Second,
		disableTypes: []DisableType{
			DisableTypeStun,
		},
		perform: func(subject, object *Unit) {
			NewPhysicalDamage(subject, object, 280).Perform()
			object.AttachHandler(NewDisable(
				object,
				DisableTypeStun,
				2*Second,
			))
		},
	}
	// Increasing critical / All
	r = &ability{
		name:               "R",
		TargetType:         TargetTypeNone,
		healthCost:         0,
		manaCost:           120,
		activationDuration: 0,
		cooldownDuration:   60 * Second,
		disableTypes: []DisableType{
			DisableTypeStun,
		},
		perform: func(subject, object *Unit) {
			for _, friend := range subject.Friends() {
				friend.AttachHandler(NewModifier(
					friend,
					unitModification{
						criticalStrikeChance: 0.2,
						criticalStrikeFactor: 0.5,
					},
					r.name,
					1,
					10*Second,
				))
			}
		},
	}
	class.abilities = []*ability{q, w, e, r}
	return class
}
