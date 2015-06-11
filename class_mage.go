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
		Perform: func(subject, object *Unit) {
			object.AttachHandler(NewCorrector(
				object,
				UnitCorrection{
					Armor: -10,
				},
				q.name,
				1,
				8*Second,
			))
			// TODO handle the error
			NewMagicDamage(subject, object, 120).Perform()
			if rand.Float64() > 0.1 {
				subject.ForSubjectHandler(subject, func(ha Handler) {
					switch ha := ha.(type) {
					case *Cooldown:
						if ha.ability == w {
							subject.DetachHandler(ha)
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
		Perform: func(subject, object *Unit) {
			// TODO handle the error
			object.AttachHandler(NewDoT(
				NewMagicDamage(subject, object, 30),
				w,
				10*Second,
			))
			if rand.Float64() > 0.2 {
				subject.ForSubjectHandler(subject, func(ha Handler) {
					switch ha := ha.(type) {
					case *Cooldown:
						if ha.ability == e {
							subject.DetachHandler(ha)
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
		Perform: func(subject, object *Unit) {
			// TODO handle the error
			NewMagicDamage(subject, object, 400).Perform()
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
		Perform: func(subject, object *Unit) {
			for _, enemy := range subject.Enemies() {
				NewMagicDamage(subject, enemy, 400).Perform()
				enemy.AttachHandler(NewDoT(
					NewMagicDamage(subject, enemy, 40),
					r,
					10*Second,
				))
				enemy.AttachHandler(NewDisable(
					object,
					DisableTypeStun,
					500*Millisecond,
				))
			}
		},
	}
	class.abilities = []*ability{q, w, e, r}
	return class
}
