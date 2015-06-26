package class

import (
	. "github.com/etheriqa/crescent"
)

func NewClassDisabler() *Class {
	var q, w, e, r Ability
	class := &Class{
		Name:                 "Disabler",
		Health:               1190,
		HealthRegeneration:   28,
		Mana:                 410,
		ManaRegeneration:     21,
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
		Perform: func(g Game, s Subject, o *Unit) {
			c := UnitCorrection{
				MagicResistance: -15,
			}
			g.Correction(o, c, q.Name, 1, 12*Second)
			_, _, _, err := MakePhysicalDamage(s, o, 160).Perform(g)
			if err != nil {
				Logger().Fatal(err)
			}
			if o.IsDead() {
				return
			}
			g.DoT(MakePhysicalDamage(s, o, 20), q.Name, 4*Second)
		},
	}
	w = Ability{
		Name:               "Suffocation",
		Description:        "Deals physical damage / Silences target for 0.5 seconds",
		TargetType:         TargetTypeEnemy,
		HealthCost:         0,
		ManaCost:           25,
		ActivationDuration: 0,
		CooldownDuration:   10 * Second,
		DisableTypes: []DisableType{
			DisableTypeStun,
			DisableTypeSilence,
		},
		Perform: func(g Game, s Subject, o *Unit) {
			_, _, _, err := MakePhysicalDamage(s, o, 275).Perform(g)
			if err != nil {
				Logger().Fatal(err)
			}
			if o.IsDead() {
				return
			}
			g.Disable(o, DisableTypeSilence, Second/2)
		},
	}
	e = Ability{
		Name:               "Call Void",
		Description:        "Deals magic damage / Stuns target for 2 seconds",
		TargetType:         TargetTypeEnemy,
		HealthCost:         0,
		ManaCost:           45,
		ActivationDuration: 2 * Second,
		CooldownDuration:   18 * Second,
		DisableTypes: []DisableType{
			DisableTypeStun,
		},
		Perform: func(g Game, s Subject, o *Unit) {
			_, _, _, err := MakeMagicDamage(s, o, 430).Perform(g)
			if err != nil {
				Logger().Fatal(err)
			}
			if o.IsDead() {
				return
			}
			g.Disable(o, DisableTypeStun, 2*Second)
		},
	}
	r = Ability{
		Name:               "Wind of Mistral",
		Description:        "Increases critical strike chance and critical strike factor for 10 seconds to all party members",
		TargetType:         TargetTypeNone,
		HealthCost:         0,
		ManaCost:           88,
		ActivationDuration: 0,
		CooldownDuration:   60 * Second,
		DisableTypes: []DisableType{
			DisableTypeStun,
		},
		Perform: func(g Game, s Subject, o *Unit) {
			c := UnitCorrection{
				CriticalStrikeChance: 0.2,
				CriticalStrikeFactor: 0.5,
			}
			g.UnitQuery().EachFriend(s.Subject(), func(u *Unit) {
				if u.IsDead() {
					return
				}
				g.Correction(u, c, r.Name, 1, 10*Second)
			})
		},
	}
	return class
}
