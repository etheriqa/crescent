package main

func NewClassHealer() *Class {
	var q, w, e, r *Ability
	class := &Class{
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
		Perform: func(op Operator, s Subject, o *Unit) {
			before, after, _, err := op.MagicDamage(s, o, 100).Perform()
			if err != nil {
				log.Fatal(err)
			}
			s.Subject().ModifyMana(op.Writer(), (before-after)*0.1)
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
		Perform: func(op Operator, s Subject, o *Unit) {
			op.HoT(op.Healing(s, o, 10), 12*Second, w.Name)
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
		Perform: func(op Operator, s Subject, o *Unit) {
			_, _, _, err := op.Healing(s, o, 400).Perform()
			if err != nil {
				log.Fatal(err)
			}
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
		Perform: func(op Operator, s Subject, o *Unit) {
			c := UnitCorrection{

				CriticalStrikeChance: 0.5,
				CriticalStrikeFactor: 1.5,
			}
			op.Correction(s.Subject(), c, 1, 6*Second, r.Name)
			op.Units().EachFriend(s.Subject(), func(friend *Unit) {
				op.HoT(op.Healing(s.Subject(), friend, 20), 12*Second, r.Name)
			})
		},
	}
	return class
}
