package main

import (
	"time"
)

func newClassTank() *class {
	class := &class{
		name: "Tank",
		// TODO stats
		health:               1000,
		healthRegeneration:   5,
		mana:                 200,
		manaRegeneration:     3,
		armor:                defaultArmor,
		magicResistance:      defaultMagicResistance,
		criticalStrikeChance: defaultCriticalStrikeChance,
		criticalStrikeFactor: defaultCriticalStrikeFactor,
		cooldownReduction:    defaultCooldownReduction,
		damageThreatFactor:   defaultDamageThreatFactor,
		healingThreatFactor:  defaultHealingThreatFactor,
	}
	// Physical / Increasing threat factor
	q := &ability{
		name:               "Mortal Breath",
		targetType:         targetTypeEnemy,
		healthCost:         0,
		manaCost:           0,
		activationDuration: 0,
		cooldownDuration:   2 * time.Second,
		disableTypes: []disableType{
			disableTypeStun,
		},
		perform: func(performer, receiver *unit) {
			// TODO remove duplicate modifiers
			performer.attachOperator(newModifier(performer, 10*time.Second, unitModification{
				damageThreatFactor: 4,
			}))
			damage := diceTrueDamage(performer, receiver, 120)
			receiver.takeDamage(damage)
			receiver.attachOperator(newDamageThreat(performer, receiver, damage))
			// TODO refactor
			if receiver.isAlive() {
				receiver.triggerEvent(eventResourceDecreased)
			} else {
				receiver.triggerEvent(eventDead)
			}
		},
	}
	// Physical / Increasing AR & MR
	w := &ability{
		name:               "Type Unsafe",
		targetType:         targetTypeEnemy,
		healthCost:         0,
		manaCost:           15,
		activationDuration: 2 * time.Second,
		cooldownDuration:   8 * time.Second,
		disableTypes: []disableType{
			disableTypeStun,
		},
		perform: func(performer, receiver *unit) {
			// TODO remove duplicate modifiers
			performer.attachOperator(newModifier(performer, 2*time.Second, unitModification{
				armor:           50,
				magicResistance: 50,
			}))
			damage := dicePhysicalDamage(performer, receiver, 200)
			receiver.takeDamage(damage)
			receiver.attachOperator(newDamageThreat(performer, receiver, damage))
			// TODO refactor
			if receiver.isAlive() {
				receiver.triggerEvent(eventResourceDecreased)
			} else {
				receiver.triggerEvent(eventDead)
			}
		},
	}
	// Physical / Life steal
	e := &ability{
		name:               "Bloody Mary",
		targetType:         targetTypeEnemy,
		healthCost:         0,
		manaCost:           50,
		activationDuration: 0,
		cooldownDuration:   15 * time.Second,
		disableTypes: []disableType{
			disableTypeStun,
		},
		perform: func(performer, receiver *unit) {
			// TODO remove duplicate modifiers
			damage := dicePhysicalDamage(performer, receiver, 300)
			before, after := receiver.takeDamage(damage)
			receiver.attachOperator(newDamageThreat(performer, receiver, damage))
			performer.takeHealing(healing(before - after))
			// TODO refactor
			if receiver.isAlive() {
				receiver.triggerEvent(eventResourceDecreased)
			} else {
				receiver.triggerEvent(eventDead)
			}
		},
	}
	// Increasing AR & MR
	r := &ability{
		name:               "Tetragrammaton",
		targetType:         targetTypeOneself,
		healthCost:         0,
		manaCost:           120,
		activationDuration: 4,
		cooldownDuration:   60,
		disableTypes: []disableType{
			disableTypeStun,
			disableTypeSilence,
		},
		perform: func(performer, receiver *unit) {
			performer.attachOperator(newModifier(performer, 8*time.Second, unitModification{
				armor:           150,
				magicResistance: 150,
			}))
		},
	}
	class.abilities = []*ability{q, w, e, r}
	return class
}
