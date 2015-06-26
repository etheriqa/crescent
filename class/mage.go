package class

import (
	"math/rand"

	. "github.com/etheriqa/crescent"
)

func NewClassMage() *Class {
	var q, w, e, r Ability
	class := &Class{
		Name:                 "Mage",
		Health:               930,
		HealthRegeneration:   23,
		Mana:                 580,
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
		Name:               "Frost Bolt",
		Description:        "Deals magic damage / Reduces 10 armor for 8 seconds to target / 20% chance to reset cooldown for Icicle",
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
			c := UnitCorrection{
				Armor: -10,
			}
			g.Correction(o, c, q.Name, 1, 8*Second)
			_, _, _, err := MakeMagicDamage(s, o, 155).Perform(g)
			if err != nil {
				Logger().Fatal(err)
			}
			if rand.Float64() > 0.1 {
				return
			}
			g.ResetCooldown(s.Subject(), &w)
		},
	}
	w = Ability{
		Name:               "Icicle",
		Description:        "Grants a magic damage over time effect for 10 seconds to target / 20% chance to reset cooldown for Absolute Zero",
		TargetType:         TargetTypeEnemy,
		HealthCost:         0,
		ManaCost:           31,
		ActivationDuration: 2 * Second,
		CooldownDuration:   8 * Second,
		DisableTypes: []DisableType{
			DisableTypeSilence,
			DisableTypeStun,
		},
		Perform: func(g Game, s Subject, o *Unit) {
			g.DoT(MakeMagicDamage(s, o, 15), w.Name, 10*Second)
			if rand.Float64() > 0.2 {
				return
			}
			g.ResetCooldown(s.Subject(), &e)
		},
	}
	e = Ability{
		Name:               "Absolute Zero",
		Description:        "Deals magic damage",
		TargetType:         TargetTypeEnemy,
		HealthCost:         0,
		ManaCost:           49,
		ActivationDuration: 2 * Second,
		CooldownDuration:   18 * Second,
		DisableTypes: []DisableType{
			DisableTypeSilence,
			DisableTypeStun,
		},
		Perform: func(g Game, s Subject, o *Unit) {
			_, _, _, err := MakeMagicDamage(s, o, 420).Perform(g)
			if err != nil {
				Logger().Fatal(err)
			}
		},
	}
	r = Ability{
		Name:               "Blizzard",
		Description:        "Deals magic damage to all enemies / Grants magic damage over time effects to all enemies",
		TargetType:         TargetTypeNone,
		HealthCost:         0,
		ManaCost:           252,
		ActivationDuration: 2 * Second,
		CooldownDuration:   60 * Second,
		DisableTypes: []DisableType{
			DisableTypeSilence,
			DisableTypeStun,
		},
		Perform: func(g Game, s Subject, o *Unit) {
			g.UnitQuery().EachEnemy(s.Subject(), func(enemy *Unit) {
				_, _, _, err := MakeMagicDamage(s, enemy, 800).Perform(g)
				if err != nil {
					Logger().Fatal(err)
				}
				if enemy.IsDead() {
					return
				}
				g.DoT(MakeMagicDamage(s, enemy, 10), r.Name, 10*Second)
			})
		},
	}
	return class
}
