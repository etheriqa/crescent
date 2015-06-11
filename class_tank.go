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
		armor:                DefaultArmor,
		magicResistance:      DefaultMagicResistance,
		criticalStrikeChance: DefaultCriticalStrikeChance,
		criticalStrikeFactor: DefaultCriticalStrikeFactor,
		cooldownReduction:    DefaultCooldownReduction,
		damageThreatFactor:   DefaultDamageThreatFactor,
		healingThreatFactor:  DefaultHealingThreatFactor,
		abilities:            []*ability{q, w, e, r},
	}
	// True damage / Increasing threat factor
	q = &ability{
		name:               "Tank Q",
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
				up.Subject(),
				UnitCorrection{
					DamageThreatFactor: 1,
				},
				q.name,
				5,
				10*Second,
			))
			// TODO handle the error
			NewTrueDamage(up, 120).Perform()
		},
	}
	// Physical damage / Increasing AR & MR
	w = &ability{
		name:               "Tank W",
		TargetType:         TargetTypeEnemy,
		healthCost:         0,
		manaCost:           15,
		activationDuration: 2 * Second,
		cooldownDuration:   8 * Second,
		disableTypes: []DisableType{
			DisableTypeStun,
		},
		Perform: func(up UnitPair) {
			up.AttachHandler(NewCorrector(
				up.Subject(),
				UnitCorrection{
					Armor:           50,
					MagicResistance: 50,
				},
				w.name,
				1,
				2*Second,
			))
			// TODO handle the error
			NewPhysicalDamage(up, 200).Perform()
		},
	}
	// Physical damage / Life steal
	e = &ability{
		name:               "Tank E",
		TargetType:         TargetTypeEnemy,
		healthCost:         0,
		manaCost:           50,
		activationDuration: 0,
		cooldownDuration:   15 * Second,
		disableTypes: []DisableType{
			DisableTypeStun,
		},
		Perform: func(up UnitPair) {
			// TODO handle the error
			before, after, _, _ := NewPhysicalDamage(up, 300).Perform()
			// TODO handle the error
			NewPureHealing(MakeUnitPair(up.Subject(), up.Subject()), (before-after)*0.6).Perform()
		},
	}
	// Increasing AR & MR
	r = &ability{
		name:               "Tank R",
		TargetType:         TargetTypeNone,
		healthCost:         0,
		manaCost:           120,
		activationDuration: 4,
		cooldownDuration:   60,
		disableTypes: []DisableType{
			DisableTypeStun,
			DisableTypeSilence,
		},
		Perform: func(up UnitPair) {
			up.AttachHandler(NewCorrector(
				up.Subject(),
				UnitCorrection{
					Armor:           150,
					MagicResistance: 150,
				},
				r.name,
				1,
				8*Second,
			))
		},
	}
	return class
}
