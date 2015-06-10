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
		perform: func(performer, receiver *unit) {
			receiver.attachOperator(newModifier(
				receiver,
				unitModification{
					armor: -10,
				},
				q,
				1,
				8*second,
			))
			// TODO handle the error
			newMagicDamage(performer, receiver, 120, q.name).perform(performer.game)
			if rand.Float64() > 0.1 {
				for o := range performer.operators {
					switch o := o.(type) {
					case *cooldown:
						if o.ability == w {
							performer.detachOperator(o)
						}
					}
				}
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
		perform: func(performer, receiver *unit) {
			// TODO handle the error
			receiver.attachOperator(newDoT(
				newMagicDamage(performer, receiver, 30, w.name),
				10*second,
			))
			if rand.Float64() > 0.2 {
				for o := range performer.operators {
					switch o := o.(type) {
					case *cooldown:
						if o.ability == e {
							performer.detachOperator(o)
						}
					}
				}
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
		perform: func(performer, receiver *unit) {
			// TODO handle the error
			newMagicDamage(performer, receiver, 400, e.name).perform(performer.game)
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
		perform: func(performer, receiver *unit) {
			for _, enemy := range performer.game.enemies(performer) {
				newMagicDamage(performer, enemy, 400, r.name).perform(performer.game)
				enemy.attachOperator(newDoT(
					newMagicDamage(performer, enemy, 40, r.name),
					10*second,
				))
				enemy.attachOperator(newDisable(
					receiver,
					disableTypeStun,
					500*millisecond,
				))
			}
		},
	}
	return class
}
