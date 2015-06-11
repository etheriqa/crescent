package main

import (
	"math/rand"
)

func newClassMage() *class {
	var q, w, e, r *ability
	class := &class{
		name: "Mage",
		// TODO stats
		health:               600,
		healthRegeneration:   2,
		mana:                 400,
		manaRegeneration:     6,
		armor:                DefaultArmor,
		magicResistance:      DefaultMagicResistance,
		criticalStrikeChance: DefaultCriticalStrikeChance,
		criticalStrikeFactor: DefaultCriticalStrikeFactor,
		cooldownReduction:    DefaultCooldownReduction,
		damageThreatFactor:   DefaultDamageThreatFactor,
		healingThreatFactor:  DefaultHealingThreatFactor,
		abilities:            []*ability{q, w, e, r},
	}
	// Magic damage / Armor reduction / Proc 10% W
	q = &ability{
		name:               "Q",
		TargetType:         TargetTypeEnemy,
		healthCost:         0,
		manaCost:           0,
		activationDuration: 2 * Second,
		cooldownDuration:   0,
		disableTypes: []DisableType{
			DisableTypeSilence,
			DisableTypeStun,
		},
		Perform: func(up UnitPair) {
			up.AttachHandler(NewCorrector(
				up.Object(),
				UnitCorrection{
					Armor: -10,
				},
				q.name,
				1,
				8*Second,
			))
			// TODO handle the error
			NewMagicDamage(up, 120).Perform()
			if rand.Float64() > 0.1 {
				up.ForSubjectHandler(up.Subject(), func(ha Handler) {
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
	w = &ability{
		name:               "W",
		TargetType:         TargetTypeEnemy,
		healthCost:         0,
		manaCost:           20,
		activationDuration: 2 * Second,
		cooldownDuration:   8 * Second,
		disableTypes: []DisableType{
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
				up.ForSubjectHandler(up.Subject(), func(ha Handler) {
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
	e = &ability{
		name:               "E",
		TargetType:         TargetTypeEnemy,
		healthCost:         0,
		manaCost:           60,
		activationDuration: 2 * Second,
		cooldownDuration:   18 * Second,
		disableTypes: []DisableType{
			DisableTypeSilence,
			DisableTypeStun,
		},
		Perform: func(up UnitPair) {
			// TODO handle the error
			NewMagicDamage(up, 400).Perform()
		},
	}
	// Magic damage / All / DoT / Stun
	r = &ability{
		name:               "R",
		TargetType:         TargetTypeNone,
		healthCost:         0,
		manaCost:           200,
		activationDuration: 2 * Second,
		cooldownDuration:   60 * Second,
		disableTypes: []DisableType{
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
	return class
}
