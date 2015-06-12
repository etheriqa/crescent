package main

func newClassDisabler() *Class {
	var q, w, e, r *Ability
	Class := &Class{
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
		Perform: func(up UnitPair) {
			up.AttachHandler(NewCorrector(
				up.Object(),
				UnitCorrection{
					MagicResistance: -15,
				},
				q.Name,
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
		Perform: func(up UnitPair) {
			for _, friend := range up.Subject().Friends() {
				friend.AttachHandler(NewCorrector(
					friend,
					UnitCorrection{
						CriticalStrikeChance: 0.2,
						CriticalStrikeFactor: 0.5,
					},
					r.Name,
					1,
					10*Second,
				))
			}
		},
	}
	return Class
}
