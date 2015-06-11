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
			receiver.attachHandler(NewModifier(
				receiver,
				unitModification{
					armor: -10,
				},
				q.name,
				1,
				8*second,
			))
			// TODO handle the error
			newMagicDamage(performer, receiver, 120).perform(performer.game)
			if rand.Float64() > 0.1 {
				for ha := range performer.handlers {
					switch ha := ha.(type) {
					case *Cooldown:
						if ha.ability == w {
							performer.detachHandler(ha)
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
			receiver.attachHandler(NewDoT(
				newMagicDamage(performer, receiver, 30),
				w,
				10*second,
			))
			if rand.Float64() > 0.2 {
				for ha := range performer.handlers {
					switch ha := ha.(type) {
					case *Cooldown:
						if ha.ability == e {
							performer.detachHandler(ha)
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
			newMagicDamage(performer, receiver, 400).perform(performer.game)
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
				newMagicDamage(performer, enemy, 400).perform(performer.game)
				enemy.attachHandler(NewDoT(
					newMagicDamage(performer, enemy, 40),
					r,
					10*second,
				))
				enemy.attachHandler(NewDisable(
					receiver,
					disableTypeStun,
					500*millisecond,
				))
			}
		},
	}
	class.abilities = []*ability{q, w, e, r}
	return class
}
