package main

import (
	"math/rand"
)

func newClassMage() *Class {
	var q, w, e, r *Ability
	class := &Class{
		Name: "Mage",
		// TODO stats
		Health:               600,
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
	// Magic damage / Armor reduction / Proc 10% W
	q = &Ability{
		Name:               "Q",
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
			_, _, _, err := op.MagicDamage(s, o, 120).Perform()
			if err != nil {
				log.Fatal(err)
			}
			if rand.Float64() > 0.1 {
				return
			}
			op.Handlers().BindObject(s.Subject()).Each(func(h Handler) {
				switch h := h.(type) {
				case *Cooldown:
					if h.Ability() == w {
						op.Handlers().Detach(h)
					}
				}
			})
		},
	}
	// Magic damage / DoT / Proc 20% E
	w = &Ability{
		Name:               "W",
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
			op.DoT(op.MagicDamage(s, o, 30), 10*Second, w.Name)
			if rand.Float64() > 0.2 {
				return
			}
			op.Handlers().BindObject(s.Subject()).Each(func(h Handler) {
				switch h := h.(type) {
				case *Cooldown:
					if h.Ability() == e {
						op.Handlers().Detach(h)
					}
				}
			})
		},
	}
	// Magic damage
	e = &Ability{
		Name:               "E",
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
			_, _, _, err := op.MagicDamage(s, o, 400).Perform()
			if err != nil {
				log.Fatal(err)
			}
		},
	}
	// Magic damage / All / DoT / Stun
	r = &Ability{
		Name:               "R",
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
				_, _, _, err := op.MagicDamage(s, enemy, 400).Perform()
				if err != nil {
					log.Fatal(err)
				}
				op.DoT(op.MagicDamage(s, enemy, 20), 10*Second, r.Name)
				op.Disable(enemy, DisableTypeStun, Second)
			})
		},
	}
	return class
}
