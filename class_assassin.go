package main

const assassinStack string = "Assassin Stack"

func newAssassinStack(subject *Unit) *Corrector {
	return NewCorrector(
		subject,
		UnitCorrection{
			CriticalStrikeChance: 0.05,
			CriticalStrikeFactor: 0.1,
		},
		assassinStack,
		10,
		10*Second,
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
		armor:                DefaultArmor,
		magicResistance:      DefaultMagicResistance,
		criticalStrikeChance: DefaultCriticalStrikeChance + 0.05,
		criticalStrikeFactor: DefaultCriticalStrikeFactor + 0.5,
		cooldownReduction:    DefaultCooldownReduction,
		damageThreatFactor:   DefaultDamageThreatFactor,
		healingThreatFactor:  DefaultHealingThreatFactor,
	}
	// Physical damage
	q = &ability{
		name:               "Q",
		TargetType:         TargetTypeEnemy,
		healthCost:         0,
		manaCost:           0,
		activationDuration: 0,
		cooldownDuration:   2 * Second,
		disableTypes: []DisableType{
			DisableTypeStun,
		},
		Perform: func(subject, object *Unit) {
			NewPhysicalDamage(subject, object, 140).Perform()
		},
	}
	// Physical damage / DoT / Increasing stacks
	w = &ability{
		name:               "W",
		TargetType:         TargetTypeEnemy,
		healthCost:         0,
		manaCost:           20,
		activationDuration: 0,
		cooldownDuration:   8 * Second,
		disableTypes: []DisableType{
			DisableTypeStun,
		},
		Perform: func(subject, object *Unit) {
			NewPhysicalDamage(subject, object, 80).Perform()
			object.AttachHandler(NewDoT(
				NewPhysicalDamage(subject, object, 20),
				w,
				10*Second,
			))
		},
	}
	// Increasing stacks / Decreasing armor and magic resistance
	e = &ability{
		name:               "E",
		TargetType:         TargetTypeNone,
		healthCost:         0,
		manaCost:           40,
		activationDuration: 0,
		cooldownDuration:   20 * Second,
		disableTypes: []DisableType{
			DisableTypeStun,
		},
		Perform: func(subject, object *Unit) {
			subject.AttachHandler(NewCorrector(
				subject,
				UnitCorrection{
					Armor:           -25,
					MagicResistance: -25,
				},
				e.name,
				1,
				8*Second,
			))
			for i := 0; i < 2; i++ {
				subject.AttachHandler(newAssassinStack(subject))
			}
		},
	}
	// Physical
	r = &ability{
		name:               "R",
		TargetType:         TargetTypeEnemy,
		healthCost:         0,
		manaCost:           120,
		activationDuration: 0,
		cooldownDuration:   60 * Second,
		disableTypes: []DisableType{
			DisableTypeStun,
		},
		Perform: func(subject, object *Unit) {
			stack := Statistic(0)
			subject.ForSubjectHandler(subject, func(ha Handler) {
				switch ha := ha.(type) {
				case *Corrector:
					if ha.name == assassinStack {
						stack += Statistic(ha.Stack())
					}
				}
			})
			NewPhysicalDamage(subject, object, 400*stack*100).Perform()
			subject.ForSubjectHandler(subject, func(ha Handler) {
				switch ha := ha.(type) {
				case *Corrector:
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
