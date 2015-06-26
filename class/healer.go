package class

import (
	. "github.com/etheriqa/crescent"
)

func NewClassHealer() *Class {
	var q, w, e, r Ability
	class := &Class{
		Name:                 "Healer",
		Health:               1030,
		HealthRegeneration:   24,
		Mana:                 550,
		ManaRegeneration:     31,
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
		Perform: func(g Game, s Subject, o *Unit) {
			before, after, _, err := MakeMagicDamage(s, o, 175).Perform(g)
			if err != nil {
				Logger().Fatal(err)
			}
			s.Subject().ModifyMana(g.Writer(), (before-after)*0.05)
		},
	}
	w = Ability{
		Name:               "Cure",
		Description:        "Restores target's health",
		TargetType:         TargetTypeFriend,
		HealthCost:         0,
		ManaCost:           25,
		ActivationDuration: 2 * Second,
		CooldownDuration:   1 * Second,
		DisableTypes: []DisableType{
			DisableTypeSilence,
			DisableTypeStun,
		},
		Perform: func(g Game, s Subject, o *Unit) {
			_, _, _, err := MakeHealing(s, o, 620).Perform(g)
			if err != nil {
				Logger().Fatal(err)
			}
		},
	}
	e = Ability{
		Name:               "Embrace",
		Description:        "Grants a healing over time effect for 8 seconds to target",
		TargetType:         TargetTypeFriend,
		HealthCost:         0,
		ManaCost:           46,
		ActivationDuration: 3 * Second,
		CooldownDuration:   7 * Second,
		DisableTypes: []DisableType{
			DisableTypeSilence,
			DisableTypeStun,
		},
		Perform: func(g Game, s Subject, o *Unit) {
			g.HoT(MakeHealing(s.Subject(), o, 55), e.Name, 12*Second)
		},
	}
	r = Ability{
		Name:               "Ascension",
		Description:        "Restores health to all party members / Grunt healing over time effects for 8 seconds to all party members / Increases critical strike chance and critical strike factor for 8 seconds",
		TargetType:         TargetTypeNone,
		HealthCost:         0,
		ManaCost:           168,
		ActivationDuration: 4 * Second,
		CooldownDuration:   60 * Second,
		DisableTypes: []DisableType{
			DisableTypeSilence,
			DisableTypeStun,
		},
		Perform: func(g Game, s Subject, o *Unit) {
			c := UnitCorrection{
				CriticalStrikeChance: 0.5,
				CriticalStrikeFactor: 1.5,
			}
			g.Correction(s.Subject(), c, r.Name, 1, 8*Second)
			g.UnitQuery().EachFriend(s.Subject(), func(friend *Unit) {
				if friend.IsDead() {
					return
				}
				_, _, _, err := MakeHealing(s.Subject(), friend, 425).Perform(g)
				if err != nil {
					Logger().Fatal(err)
				}
				g.HoT(MakeHealing(s.Subject(), friend, 25), r.Name, 8*Second)
			})
		},
	}
	return class
}
