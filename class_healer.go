package main

func newClassHealer() *Class {
	var q, w, e, r *Ability
	Class := &Class{
		Name: "Healer",
		// TODO stats
		Health:               700,
		HealthRegeneration:   2,
		Mana:                 400,
		ManaRegeneration:     6,
		Armor:                DefaultArmor,
		MagicResistance:      DefaultMagicResistance,
		CriticalStrikeChance: DefaultCriticalStrikeChance,
		CriticalStrikeFactor: DefaultCriticalStrikeFactor,
		CooldownReduction:    DefaultCooldownReduction,
		DamageThreatFactor:   DefaultDamageThreatFactor,
		HealingThreatFactor:  DefaultHealingThreatFactor,
		Abilities:            []*Ability{q, w, e, r},
	}
	// Magic damage / Mana restore
	q = &Ability{
		Name:               "Healer Q",
		TargetType:         TargetTypeEnemy,
		HealthCost:         0,
		ManaCost:           0,
		ActivationDuration: 2 * Second,
		CooldownDuration:   2 * Second,
		DisableTypes: []DisableType{
			DisableTypeSilence,
			DisableTypeStun,
		},
		Perform: func(up UnitPair) {
			// TODO handle the error
			before, after, _, _ := NewMagicDamage(up, 100).Perform()
			up.Subject().ModifyMana((before - after) * 0.1)
		},
	}
	// HoT
	w = &Ability{
		Name:               "Healer W",
		TargetType:         TargetTypeFriend,
		HealthCost:         0,
		ManaCost:           40,
		ActivationDuration: 2 * Second,
		CooldownDuration:   4 * Second,
		DisableTypes: []DisableType{
			DisableTypeSilence,
			DisableTypeStun,
		},
		Perform: func(up UnitPair) {
			up.AttachHandler(NewTicker(
				NewHealing(up, 20),
				w,
				12*Second,
			))
		},
	}
	// Healing
	e = &Ability{
		Name:               "Healer E",
		TargetType:         TargetTypeFriend,
		HealthCost:         0,
		ManaCost:           80,
		ActivationDuration: 2 * Second,
		CooldownDuration:   8 * Second,
		DisableTypes: []DisableType{
			DisableTypeSilence,
			DisableTypeStun,
		},
		Perform: func(up UnitPair) {
			NewHealing(up, 400).Perform()
		},
	}
	// HoT / Increasing critical strike chance and critical strike factor
	r = &Ability{
		Name:               "Healer R",
		TargetType:         TargetTypeNone,
		HealthCost:         0,
		ManaCost:           200,
		ActivationDuration: 0,
		CooldownDuration:   60 * Second,
		DisableTypes: []DisableType{
			DisableTypeSilence,
			DisableTypeStun,
		},
		Perform: func(up UnitPair) {
			up.AttachHandler(NewCorrector(
				up.Subject(),
				UnitCorrection{
					CriticalStrikeChance: 0.5,
					CriticalStrikeFactor: 1.5,
				},
				r.Name,
				1,
				6*Second,
			))
			for _, friend := range up.Subject().Friends() {
				up.AttachHandler(NewTicker(
					NewHealing(MakeUnitPair(up.Subject(), friend), 20),
					r,
					6*Second,
				))
			}
		},
	}
	return Class
}
