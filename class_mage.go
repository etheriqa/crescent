package main

import (
	"math/rand"
)

func NewClassMage() *Class {
	var q, w, e, r Ability
	class := &Class{
		Name:                 "Mage",
		Health:               600,
		HealthRegeneration:   14,
		Mana:                 400,
		ManaRegeneration:     28,
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
		Perform: func(op Operator, s Subject, o *Unit) {
			c := UnitCorrection{
				Armor: -10,
			}
			op.Correction(o, c, 1, 8*Second, q.Name)
			_, _, _, err := op.MagicDamage(s, o, 155).Perform()
			if err != nil {
				log.Fatal(err)
			}
			if rand.Float64() > 0.1 {
				return
			}
			op.Handlers().BindObject(s.Subject()).Each(func(h Handler) {
				switch h := h.(type) {
				case *Cooldown:
					if h.Ability() == &w {
						op.Handlers().Detach(h)
					}
				}
			})
		},
	}
	w = Ability{
		Name:               "Icicle",
		Description:        "Grants a magic damage over time effect for 10 seconds to target / 20% chance to reset cooldown for Absolute Zero",
		TargetType:         TargetTypeEnemy,
		HealthCost:         0,
		ManaCost:           20,
		ActivationDuration: 2 * Second,
		CooldownDuration:   8 * Second,
		DisableTypes: []DisableType{
			DisableTypeSilence,
			DisableTypeStun,
		},
		Perform: func(op Operator, s Subject, o *Unit) {
			op.DoT(op.MagicDamage(s, o, 15), 10*Second, w.Name)
			if rand.Float64() > 0.2 {
				return
			}
			op.Handlers().BindObject(s.Subject()).Each(func(h Handler) {
				switch h := h.(type) {
				case *Cooldown:
					if h.Ability() == &e {
						op.Handlers().Detach(h)
					}
				}
			})
		},
	}
	e = Ability{
		Name:               "Absolute Zero",
		Description:        "Deals magic damage",
		TargetType:         TargetTypeEnemy,
		HealthCost:         0,
		ManaCost:           60,
		ActivationDuration: 2 * Second,
		CooldownDuration:   18 * Second,
		DisableTypes: []DisableType{
			DisableTypeSilence,
			DisableTypeStun,
		},
		Perform: func(op Operator, s Subject, o *Unit) {
			_, _, _, err := op.MagicDamage(s, o, 420).Perform()
			if err != nil {
				log.Fatal(err)
			}
		},
	}
	r = Ability{
		Name:               "Bllizard",
		Description:        "Deals magic damage to all enemies / Grants magic damage over time effects to all enemies",
		TargetType:         TargetTypeNone,
		HealthCost:         0,
		ManaCost:           200,
		ActivationDuration: 2 * Second,
		CooldownDuration:   60 * Second,
		DisableTypes: []DisableType{
			DisableTypeSilence,
			DisableTypeStun,
		},
		Perform: func(op Operator, s Subject, o *Unit) {
			op.Units().EachEnemy(s.Subject(), func(enemy *Unit) {
				_, _, _, err := op.MagicDamage(s, enemy, 600).Perform()
				if err != nil {
					log.Fatal(err)
				}
				if enemy.IsDead() {
					return
				}
				op.DoT(op.MagicDamage(s, enemy, 10), 10*Second, r.Name)
			})
		},
	}
	return class
}
