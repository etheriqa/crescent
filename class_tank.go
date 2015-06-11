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
		Perform: func(subject, object *Unit) {
			subject.AttachHandler(NewCorrector(
				subject,
				UnitCorrection{
					DamageThreatFactor: 1,
				},
				q.name,
				5,
				10*Second,
			))
			// TODO handle the error
			NewTrueDamage(subject, object, 120).Perform()
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
		Perform: func(subject, object *Unit) {
			subject.AttachHandler(NewCorrector(
				subject,
				UnitCorrection{
					Armor:           50,
					MagicResistance: 50,
				},
				w.name,
				1,
				2*Second,
			))
			// TODO handle the error
			NewPhysicalDamage(subject, object, 200).Perform()
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
		Perform: func(subject, object *Unit) {
			// TODO handle the error
			before, after, _, _ := NewPhysicalDamage(subject, object, 300).Perform()
			// TODO handle the error
			NewPureHealing(subject, object, (before-after)*0.6).Perform()
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
		Perform: func(subject, object *Unit) {
			subject.AttachHandler(NewCorrector(
				subject,
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
	class.abilities = []*ability{q, w, e, r}
	return class
}
