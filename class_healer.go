package main

func NewClassHealer() *Class {
	var q, w, e, r Ability
	class := &Class{
		Name:                 "Healer",
		Health:               700,
		HealthRegeneration:   15,
		Mana:                 400,
		ManaRegeneration:     33,
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
		Name:               "Conversion",
		Description:        "Deals magic damage / Restores mana",
		TargetType:         TargetTypeEnemy,
		HealthCost:         0,
		ManaCost:           0,
		ActivationDuration: 2 * Second,
		CooldownDuration:   0,
		DisableTypes: []DisableType{
			DisableTypeSilence,
			DisableTypeStun,
		},
		Perform: func(op Operator, s Subject, o *Unit) {
			before, after, _, err := op.MagicDamage(s, o, 175).Perform()
			if err != nil {
				log.Fatal(err)
			}
			s.Subject().ModifyMana(op.Writer(), (before-after)*0.1)
		},
	}
	w = Ability{
		Name:               "Cure",
		Description:        "Restores target's health",
		TargetType:         TargetTypeFriend,
		HealthCost:         0,
		ManaCost:           40,
		ActivationDuration: 2 * Second,
		CooldownDuration:   1 * Second,
		DisableTypes: []DisableType{
			DisableTypeSilence,
			DisableTypeStun,
		},
		Perform: func(op Operator, s Subject, o *Unit) {
			_, _, _, err := op.Healing(s, o, 620).Perform()
			if err != nil {
				log.Fatal(err)
			}
		},
	}
	e = Ability{
		Name:               "Embrace",
		Description:        "Grants a healing over time effect for 8 seconds to target",
		TargetType:         TargetTypeFriend,
		HealthCost:         0,
		ManaCost:           80,
		ActivationDuration: 3 * Second,
		CooldownDuration:   7 * Second,
		DisableTypes: []DisableType{
			DisableTypeSilence,
			DisableTypeStun,
		},
		Perform: func(op Operator, s Subject, o *Unit) {
			op.HoT(op.Healing(s.Subject(), o, 65), 8*Second, e.Name)
		},
	}
	r = Ability{
		Name:               "Ascension",
		Description:        "Restores health to all party members / Grunt healing over time effects for 8 seconds to all party members / Increases critical strike chance and critical strike factor for 8 seconds",
		TargetType:         TargetTypeNone,
		HealthCost:         0,
		ManaCost:           200,
		ActivationDuration: 4 * Second,
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
			op.Correction(s.Subject(), c, 1, 8*Second, r.Name)
			op.Units().EachFriend(s.Subject(), func(friend *Unit) {
				if friend.IsDead() {
					return
				}
				_, _, _, err := op.Healing(s.Subject(), friend, 425).Perform()
				if err != nil {
					log.Fatal(err)
				}
				op.HoT(op.Healing(s.Subject(), friend, 25), 8*Second, r.Name)
			})
		},
	}
	return class
}
