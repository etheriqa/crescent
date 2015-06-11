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
		abilities:            []*ability{q, w, e, r},
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
		Perform: func(up UnitPair) {
			NewPhysicalDamage(up, 140).Perform()
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
		Perform: func(up UnitPair) {
			NewPhysicalDamage(up, 80).Perform()
			up.AttachHandler(NewTicker(
				NewPhysicalDamage(up, 20),
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
		Perform: func(up UnitPair) {
			up.AttachHandler(NewCorrector(
				up.Subject(),
				UnitCorrection{
					Armor:           -25,
					MagicResistance: -25,
				},
				e.name,
				1,
				8*Second,
			))
			for i := 0; i < 2; i++ {
				up.AttachHandler(newAssassinStack(up.Subject()))
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
		Perform: func(up UnitPair) {
			stack := Statistic(0)
			up.ForSubjectHandler(up.Subject(), func(ha Handler) {
				switch ha := ha.(type) {
				case *Corrector:
					if ha.name == assassinStack {
						stack += Statistic(ha.Stack())
					}
				}
			})
			NewPhysicalDamage(up, 400*stack*100).Perform()
			up.ForSubjectHandler(up.Subject(), func(ha Handler) {
				switch ha := ha.(type) {
				case *Corrector:
					if ha.name == assassinStack {
						up.DetachHandler(ha)
					}
				}
			})
		},
	}
	return class
}
