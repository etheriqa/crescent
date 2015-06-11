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
		armor:                defaultArmor,
		magicResistance:      defaultMagicResistance,
		criticalStrikeChance: defaultCriticalStrikeChance,
		criticalStrikeFactor: defaultCriticalStrikeFactor,
		cooldownReduction:    defaultCooldownReduction,
		damageThreatFactor:   defaultDamageThreatFactor,
		healingThreatFactor:  defaultHealingThreatFactor,
	}
	// Magic damage / Armor reduction / Proc 10% W
	q = &ability{
		name:               "Q",
		targetType:         targetTypeEnemy,
		healthCost:         0,
		manaCost:           0,
		activationDuration: 2 * second,
		cooldownDuration:   0,
		disableTypes: []disableType{
			disableTypeSilence,
			disableTypeStun,
		},
		perform: func(subject, object *unit) {
			object.AttachHandler(NewModifier(
				object,
				unitModification{
					armor: -10,
				},
				q.name,
				1,
				8*second,
			))
			// TODO handle the error
			newMagicDamage(subject, object, 120).perform(subject.game)
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
		targetType:         targetTypeEnemy,
		healthCost:         0,
		manaCost:           20,
		activationDuration: 2 * second,
		cooldownDuration:   8 * second,
		disableTypes: []disableType{
			disableTypeSilence,
			disableTypeStun,
		},
		perform: func(subject, object *unit) {
			// TODO handle the error
			object.AttachHandler(NewDoT(
				newMagicDamage(subject, object, 30),
				w,
				10*second,
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
		targetType:         targetTypeEnemy,
		healthCost:         0,
		manaCost:           60,
		activationDuration: 2 * second,
		cooldownDuration:   18 * second,
		disableTypes: []disableType{
			disableTypeSilence,
			disableTypeStun,
		},
		perform: func(subject, object *unit) {
			// TODO handle the error
			newMagicDamage(subject, object, 400).perform(subject.game)
		},
	}
	// Magic damage / All / DoT / Stun
	r = &ability{
		name:               "R",
		targetType:         targetTypeNone,
		healthCost:         0,
		manaCost:           200,
		activationDuration: 2 * second,
		cooldownDuration:   60 * second,
		disableTypes: []disableType{
			disableTypeSilence,
			disableTypeStun,
		},
		perform: func(subject, object *unit) {
			for _, enemy := range subject.game.enemies(subject) {
				newMagicDamage(subject, enemy, 400).perform(subject.game)
				enemy.AttachHandler(NewDoT(
					newMagicDamage(subject, enemy, 40),
					r,
					10*second,
				))
				enemy.AttachHandler(NewDisable(
					object,
					disableTypeStun,
					500*millisecond,
				))
			}
		},
	}
	class.abilities = []*ability{q, w, e, r}
	return class
}
