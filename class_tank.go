package main

func NewClassTank() *Class {
	var q, w, e, r Ability
	class := &Class{
		Name:                 "Tank",
		Health:               1740,
		HealthRegeneration:   41,
		Mana:                 340,
		ManaRegeneration:     16,
		Armor:                DefaultArmor,
		MagicResistance:      DefaultMagicResistance,
		CriticalStrikeChance: DefaultCriticalStrikeChance,
		CriticalStrikeFactor: DefaultCriticalStrikeFactor,
		CooldownReduction:    DefaultCooldownReduction,
		DamageThreatFactor:   DefaultDamageThreatFactor,
		HealingThreatFactor:  DefaultHealingThreatFactor,
		Abilities:            []*Ability{&q, &w, &e, &r},
	}
	q = Ability{
		Name:               "Primary Gun",
		Description:        "Deals true damage / Increases threat factor for 8 seconds",
		TargetType:         TargetTypeEnemy,
		HealthCost:         0,
		ManaCost:           0,
		ActivationDuration: 0,
		CooldownDuration:   2 * Second,
		DisableTypes: []DisableType{
			DisableTypeStun,
		},
		Perform: func(g Game, s Subject, o *Unit) {
			c := UnitCorrection{
				DamageThreatFactor: 0.4,
			}
			g.Correction(s.Subject(), c, 5, 8*Second, q.Name)
			_, _, _, err := g.TrueDamage(s, o, 120).Perform()
			if err != nil {
				log.Fatal(err)
			}
		},
	}
	w = Ability{
		Name:               "Alertness",
		Description:        "Deals physical damage / Increases armor and magic resistance for 4 seconds",
		TargetType:         TargetTypeEnemy,
		HealthCost:         0,
		ManaCost:           9,
		ActivationDuration: 1 * Second,
		CooldownDuration:   10 * Second,
		DisableTypes: []DisableType{
			DisableTypeStun,
		},
		Perform: func(g Game, s Subject, o *Unit) {
			c := UnitCorrection{
				Armor:           50,
				MagicResistance: 50,
			}
			g.Correction(s.Subject(), c, 1, 4*Second, w.Name)
			_, _, _, err := g.PhysicalDamage(s, o, 200).Perform()
			if err != nil {
				log.Fatal(err)
			}
		},
	}
	e = Ability{
		Name:               "Blood Sword",
		Description:        "Deals physical damage / Drains health",
		TargetType:         TargetTypeEnemy,
		HealthCost:         0,
		ManaCost:           33,
		ActivationDuration: 0,
		CooldownDuration:   15 * Second,
		DisableTypes: []DisableType{
			DisableTypeStun,
		},
		Perform: func(g Game, s Subject, o *Unit) {
			before, after, _, err := g.PhysicalDamage(s, o, 345).Perform()
			if err != nil {
				log.Fatal(err)
			}
			s.Subject().ModifyHealth(g.Writer(), before-after)
		},
	}
	r = Ability{
		Name:               "Equilibrium",
		Description:        "Increases armor and magic resistance for 5 seconds",
		TargetType:         TargetTypeNone,
		HealthCost:         0,
		ManaCost:           61,
		ActivationDuration: 1 * Second,
		CooldownDuration:   60 * Second,
		DisableTypes: []DisableType{
			DisableTypeStun,
			DisableTypeSilence,
		},
		Perform: func(g Game, s Subject, o *Unit) {
			c := UnitCorrection{
				Armor:           150,
				MagicResistance: 150,
			}
			g.Correction(s.Subject(), c, 1, 5*Second, r.Name)
		},
	}
	return class
}
