package main

const assassinStack string = "Assassin Stack"

func newAssassinStack(subject *unit) *Modifier {
	return NewModifier(
		subject,
		unitModification{
			criticalStrikeChance: 0.05,
			criticalStrikeFactor: 0.1,
		},
		assassinStack,
		10,
		10*second,
	)
}

func newClassAssassin() *class {
	var q, w, e, r *ability
	class := &class{
		name: "Assassin",
		// TODO stats
		health:               600,
		healthRegeneration:   2,
		mana:                 200,
		manaRegeneration:     3,
		armor:                defaultArmor,
		magicResistance:      defaultMagicResistance,
		criticalStrikeChance: defaultCriticalStrikeChance + 0.05,
		criticalStrikeFactor: defaultCriticalStrikeFactor + 0.5,
		cooldownReduction:    defaultCooldownReduction,
		damageThreatFactor:   defaultDamageThreatFactor,
		healingThreatFactor:  defaultHealingThreatFactor,
	}
	// Physical damage
	q = &ability{
		name:               "Q",
		targetType:         targetTypeEnemy,
		healthCost:         0,
		manaCost:           0,
		activationDuration: 0,
		cooldownDuration:   2 * second,
		disableTypes: []disableType{
			disableTypeStun,
		},
		perform: func(subject, object *unit) {
			newPhysicalDamage(subject, object, 140).perform(subject.game)
		},
	}
	// Physical damage / DoT / Increasing stacks
	w = &ability{
		name:               "W",
		targetType:         targetTypeEnemy,
		healthCost:         0,
		manaCost:           20,
		activationDuration: 0,
		cooldownDuration:   8 * second,
		disableTypes: []disableType{
			disableTypeStun,
		},
		perform: func(subject, object *unit) {
			newPhysicalDamage(subject, object, 80).perform(subject.game)
			object.AttachHandler(NewDoT(
				newPhysicalDamage(subject, object, 20),
				w,
				10*second,
			))
		},
	}
	// Increasing stacks / Decreasing armor and magic resistance
	e = &ability{
		name:               "E",
		targetType:         targetTypeNone,
		healthCost:         0,
		manaCost:           40,
		activationDuration: 0,
		cooldownDuration:   20 * second,
		disableTypes: []disableType{
			disableTypeStun,
		},
		perform: func(subject, object *unit) {
			subject.AttachHandler(NewModifier(
				subject,
				unitModification{
					armor:           -25,
					magicResistance: -25,
				},
				e.name,
				1,
				8*second,
			))
			for i := 0; i < 2; i++ {
				subject.AttachHandler(newAssassinStack(subject))
			}
		},
	}
	// Physical
	r = &ability{
		name:               "R",
		targetType:         targetTypeEnemy,
		healthCost:         0,
		manaCost:           120,
		activationDuration: 0,
		cooldownDuration:   60 * second,
		disableTypes: []disableType{
			disableTypeStun,
		},
		perform: func(subject, object *unit) {
			stack := statistic(0)
			subject.ForSubjectHandler(subject, func(ha Handler) {
				switch ha := ha.(type) {
				case *Modifier:
					if ha.name == assassinStack {
						stack += statistic(ha.nowStack)
					}
				}
			})
			newPhysicalDamage(subject, object, 400*stack*100).perform(subject.game)
			subject.ForSubjectHandler(subject, func(ha Handler) {
				switch ha := ha.(type) {
				case *Modifier:
					if ha.name == assassinStack {
						subject.DetachHandler(ha)
					}
				}
			})
		},
	}
	class.abilities = []*ability{q, w, e, r}
	return class
}
