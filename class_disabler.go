package main

func NewClassDisabler() *Class {
	var q, w, e, r *Ability
	class := &Class{
		Name: "Disabler",
		// TODO stats
		Health:               800,
		HealthRegeneration:   2,
		Mana:                 300,
		ManaRegeneration:     4,
		Armor:                DefaultArmor,
		MagicResistance:      DefaultMagicResistance,
		CriticalStrikeChance: DefaultCriticalStrikeChance,
		CriticalStrikeFactor: DefaultCriticalStrikeFactor,
		CooldownReduction:    DefaultCooldownReduction,
		DamageThreatFactor:   DefaultDamageThreatFactor,
		HealingThreatFactor:  DefaultHealingThreatFactor,
		Abilities:            []*Ability{q, w, e, r},
	}
	// Physical damage / DoT / Magic resistance reduction
	q = &Ability{
		Name:               "Q",
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
				MagicResistance: -15,
			}
			op.Correction(o, c, 1, 12*Second, q.Name)
			_, _, _, err := op.PhysicalDamage(s, o, 110).Perform()
			if err != nil {
				log.Fatal(err)
			}
			op.DoT(op.PhysicalDamage(s, o, 12), 4*Second, q.Name)
		},
	}
	// Magic damage / Silence
	w = &Ability{
		Name:               "W",
		TargetType:         TargetTypeEnemy,
		HealthCost:         0,
		ManaCost:           40,
		ActivationDuration: 2 * Second,
		CooldownDuration:   8 * Second,
		DisableTypes: []DisableType{
			DisableTypeStun,
			DisableTypeSilence,
		},
		Perform: func(op Operator, s Subject, o *Unit) {
			_, _, _, err := op.MagicDamage(s, o, 220).Perform()
			if err != nil {
				log.Fatal(err)
			}
			op.Disable(o, DisableTypeSilence, Second)
		},
	}
	// Physical damage / Stun
	e = &Ability{
		Name:               "E",
		TargetType:         TargetTypeEnemy,
		HealthCost:         0,
		ManaCost:           60,
		ActivationDuration: 0,
		CooldownDuration:   20 * Second,
		DisableTypes: []DisableType{
			DisableTypeStun,
		},
		Perform: func(op Operator, s Subject, o *Unit) {
			_, _, _, err := op.PhysicalDamage(s, o, 280).Perform()
			if err != nil {
				log.Fatal(err)
			}
			op.Disable(o, DisableTypeStun, 2*Second)
		},
	}
	// Increasing critical / All
	r = &Ability{
		Name:               "R",
		TargetType:         TargetTypeNone,
		HealthCost:         0,
		ManaCost:           120,
		ActivationDuration: 0,
		CooldownDuration:   60 * Second,
		DisableTypes: []DisableType{
			DisableTypeStun,
		},
		Perform: func(op Operator, s Subject, o *Unit) {
			c := UnitCorrection{
				CriticalStrikeChance: 0.2,
				CriticalStrikeFactor: 0.5,
			}
			op.Units().EachFriend(s.Subject(), func(u *Unit) {
				op.Correction(u, c, 1, 10*Second, r.Name)
			})
		},
	}
	return class
}
