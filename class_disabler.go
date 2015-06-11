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
		abilities:            []*ability{q, w, e, r},
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
		Perform: func(up UnitPair) {
			up.AttachHandler(NewCorrector(
				up.Object(),
				UnitCorrection{
					MagicResistance: -15,
				},
				q.name,
				1,
				12*Second,
			))
			// TODO handle the error
			NewPhysicalDamage(up, 110).Perform()
			up.AttachHandler(NewTicker(
				NewPhysicalDamage(up, 25),
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
		Perform: func(up UnitPair) {
			NewMagicDamage(up, 220).Perform()
			up.AttachHandler(NewDisable(
				up.Object(),
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
		Perform: func(up UnitPair) {
			NewPhysicalDamage(up, 280).Perform()
			up.AttachHandler(NewDisable(
				up.Object(),
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
		Perform: func(up UnitPair) {
			for _, friend := range up.Subject().Friends() {
				friend.AttachHandler(NewCorrector(
					friend,
					UnitCorrection{
						CriticalStrikeChance: 0.2,
						CriticalStrikeFactor: 0.5,
					},
					r.name,
					1,
					10*Second,
				))
			}
		},
	}
	return class
}
