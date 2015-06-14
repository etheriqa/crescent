package main

func NewClassTank() *Class {
	var q, w, e, r Ability
	class := &Class{
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
		Abilities:            []*Ability{&q, &w, &e, &r},
	}
	// True damage / Increasing threat factor
	q = Ability{
		Name:               "Tank Q",
		TargetType:         TargetTypeEnemy,
		HealthCost:         0,
		ManaCost:           0,
		ActivationDuration: 0,
		CooldownDuration:   2 * Second,
		DisableTypes: []DisableType{
			DisableTypeStun,
		},
		Perform: func(op Operator, s Subject, o *Unit) {
			c := UnitCorrection{
				DamageThreatFactor: 1,
			}
			op.Correction(s.Subject(), c, 5, 10*Second, q.Name)
			_, _, _, err := op.TrueDamage(s, o, 120).Perform()
			if err != nil {
				log.Fatal(err)
			}
		},
	}
	// Physical damage / Increasing AR & MR
	w = Ability{
		Name:               "Tank W",
		TargetType:         TargetTypeEnemy,
		HealthCost:         0,
		ManaCost:           15,
		ActivationDuration: 2 * Second,
		CooldownDuration:   8 * Second,
		DisableTypes: []DisableType{
			DisableTypeStun,
		},
		Perform: func(op Operator, s Subject, o *Unit) {
			c := UnitCorrection{
				Armor:           50,
				MagicResistance: 50,
			}
			op.Correction(s.Subject(), c, 1, 2*Second, w.Name)
			_, _, _, err := op.PhysicalDamage(s, o, 200).Perform()
			if err != nil {
				log.Fatal(err)
			}
		},
	}
	// Physical damage / Life steal
	e = Ability{
		Name:               "Tank E",
		TargetType:         TargetTypeEnemy,
		HealthCost:         0,
		ManaCost:           50,
		ActivationDuration: 0,
		CooldownDuration:   15 * Second,
		DisableTypes: []DisableType{
			DisableTypeStun,
		},
		Perform: func(op Operator, s Subject, o *Unit) {
			before, after, _, err := op.PhysicalDamage(s, o, 300).Perform()
			if err != nil {
				log.Fatal(err)
			}
			s.Subject().ModifyHealth(op.Writer(), (before-after)*0.6)
		},
	}
	// Increasing AR & MR
	r = Ability{
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
		Perform: func(op Operator, s Subject, o *Unit) {
			c := UnitCorrection{
				Armor:           150,
				MagicResistance: 150,
			}
			op.Correction(s.Subject(), c, 1, 8*Second, r.Name)
		},
	}
	return class
}
