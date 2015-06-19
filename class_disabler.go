package main

func NewClassDisabler() *Class {
	var q, w, e, r Ability
	class := &Class{
		Name:                 "Disabler",
		Health:               800,
		HealthRegeneration:   25,
		Mana:                 300,
		ManaRegeneration:     19,
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
		Name:               "Disintegrate",
		Description:        "Deals physical damage / Grants a physical damage over time effect for 4 seconds to target / Reduces 15 magic resistance for 12 seconds to target",
		TargetType:         TargetTypeEnemy,
		HealthCost:         0,
		ManaCost:           0,
		ActivationDuration: 2 * Second,
		CooldownDuration:   0,
		DisableTypes: []DisableType{
			DisableTypeStun,
		},
		Perform: func(op Operator, s Subject, o *Unit) {
			c := UnitCorrection{
				MagicResistance: -15,
			}
			op.Correction(o, c, 1, 12*Second, q.Name)
			_, _, _, err := op.PhysicalDamage(s, o, 160).Perform()
			if err != nil {
				log.Fatal(err)
			}
			if o.IsDead() {
				return
			}
			op.DoT(op.PhysicalDamage(s, o, 20), 4*Second, q.Name)
		},
	}
	w = Ability{
		Name:               "Suffocation",
		Description:        "Deals physical damage / Silences target for 0.5 seconds",
		TargetType:         TargetTypeEnemy,
		HealthCost:         0,
		ManaCost:           40,
		ActivationDuration: 0,
		CooldownDuration:   10 * Second,
		DisableTypes: []DisableType{
			DisableTypeStun,
			DisableTypeSilence,
		},
		Perform: func(op Operator, s Subject, o *Unit) {
			_, _, _, err := op.PhysicalDamage(s, o, 275).Perform()
			if err != nil {
				log.Fatal(err)
			}
			if o.IsDead() {
				return
			}
			op.Disable(o, DisableTypeSilence, Second/2)
		},
	}
	e = Ability{
		Name:               "Call Void",
		Description:        "Deals magic damage / Stuns target for 2 seconds",
		TargetType:         TargetTypeEnemy,
		HealthCost:         0,
		ManaCost:           60,
		ActivationDuration: 2 * Second,
		CooldownDuration:   18 * Second,
		DisableTypes: []DisableType{
			DisableTypeStun,
		},
		Perform: func(op Operator, s Subject, o *Unit) {
			_, _, _, err := op.MagicDamage(s, o, 430).Perform()
			if err != nil {
				log.Fatal(err)
			}
			if o.IsDead() {
				return
			}
			op.Disable(o, DisableTypeStun, 2*Second)
		},
	}
	r = Ability{
		Name:               "Wind of Mistral",
		Description:        "Increases critical strike chance and critical strike factor for 10 seconds to all party members",
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
				if u.IsDead() {
					return
				}
				op.Correction(u, c, 1, 10*Second, r.Name)
			})
		},
	}
	return class
}
