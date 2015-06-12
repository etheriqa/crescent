package main

import (
	"math/rand"
)

func newClassMage() *Class {
	var q, w, e, r *Ability
	Class := &Class{
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
		Perform: func(up UnitPair) {
			up.AttachHandler(NewCorrector(
				up.Object(),
				UnitCorrection{
					Armor: -10,
				},
				q.Name,
				1,
				8*Second,
			))
			// TODO handle the error
			NewMagicDamage(up, 120).Perform()
			if rand.Float64() > 0.1 {
				up.ForSubjectHandler(func(ha Handler) {
					switch ha := ha.(type) {
					case *Cooldown:
						if ha.ability == w {
							up.DetachHandler(ha)
						}
					}
				})
			}
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
		Perform: func(up UnitPair) {
			// TODO handle the error
			up.AttachHandler(NewTicker(
				NewMagicDamage(up, 30),
				w,
				10*Second,
			))
			if rand.Float64() > 0.2 {
				up.ForSubjectHandler(func(ha Handler) {
					switch ha := ha.(type) {
					case *Cooldown:
						if ha.ability == e {
							up.DetachHandler(ha)
						}
					}
				})
			}
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
		Perform: func(up UnitPair) {
			// TODO handle the error
			NewMagicDamage(up, 400).Perform()
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
		Perform: func(up UnitPair) {
			for _, enemy := range up.Subject().Enemies() {
				NewMagicDamage(MakeUnitPair(up.Subject(), enemy), 400).Perform()
				up.AttachHandler(NewTicker(
					NewMagicDamage(MakeUnitPair(up.Subject(), enemy), 40),
					r,
					10*Second,
				))
				up.AttachHandler(NewDisable(
					enemy,
					DisableTypeStun,
					500*Millisecond,
				))
			}
		},
	}
	return Class
}
