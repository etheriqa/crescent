package main

func newClassTank() *Class {
	var q, w, e, r *Ability
	Class := &Class{
		Name: "Tank",
		// TODO stats
		Health:               1000,
		HealthRegeneration:   5,
		Mana:                 200,
		ManaRegeneration:     3,
		Armor:                DefaultArmor,
		MagicResistance:      DefaultMagicResistance,
		CriticalStrikeChance: DefaultCriticalStrikeChance,
		CriticalStrikeFactor: DefaultCriticalStrikeFactor,
		CooldownReduction:    DefaultCooldownReduction,
		DamageThreatFactor:   DefaultDamageThreatFactor,
		HealingThreatFactor:  DefaultHealingThreatFactor,
		Abilities:            []*Ability{q, w, e, r},
	}
	// True damage / Increasing threat factor
	q = &Ability{
		Name:               "Tank Q",
		TargetType:         TargetTypeEnemy,
		HealthCost:         0,
		ManaCost:           0,
		ActivationDuration: 0,
		CooldownDuration:   2 * Second,
		DisableTypes: []DisableType{
			DisableTypeStun,
		},
		Perform: func(up UnitPair) {
			up.AttachHandler(NewCorrector(
				up.Subject(),
				UnitCorrection{
					DamageThreatFactor: 1,
				},
				q.Name,
				5,
				10*Second,
			))
			// TODO handle the error
			NewTrueDamage(up, 120).Perform()
		},
	}
	// Physical damage / Increasing AR & MR
	w = &Ability{
		Name:               "Tank W",
		TargetType:         TargetTypeEnemy,
		HealthCost:         0,
		ManaCost:           15,
		ActivationDuration: 2 * Second,
		CooldownDuration:   8 * Second,
		DisableTypes: []DisableType{
			DisableTypeStun,
		},
		Perform: func(up UnitPair) {
			up.AttachHandler(NewCorrector(
				up.Subject(),
				UnitCorrection{
					Armor:           50,
					MagicResistance: 50,
				},
				w.Name,
				1,
				2*Second,
			))
			// TODO handle the error
			NewPhysicalDamage(up, 200).Perform()
		},
	}
	// Physical damage / Life steal
	e = &Ability{
		Name:               "Tank E",
		TargetType:         TargetTypeEnemy,
		HealthCost:         0,
		ManaCost:           50,
		ActivationDuration: 0,
		CooldownDuration:   15 * Second,
		DisableTypes: []DisableType{
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
	r = &Ability{
		Name:               "Tank R",
		TargetType:         TargetTypeNone,
		HealthCost:         0,
		ManaCost:           120,
		ActivationDuration: 4,
		CooldownDuration:   60,
		DisableTypes: []DisableType{
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
				r.Name,
				1,
				8*Second,
			))
		},
	}
	return Class
}
