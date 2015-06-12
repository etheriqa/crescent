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

func newClassAssassin() *Class {
	var q, w, e, r *Ability
	Class := &Class{
		Name: "Assassin",
		// TODO stats
		Health:               600,
		HealthRegeneration:   2,
		Mana:                 200,
		ManaRegeneration:     3,
		Armor:                DefaultArmor,
		MagicResistance:      DefaultMagicResistance,
		CriticalStrikeChance: DefaultCriticalStrikeChance + 0.05,
		CriticalStrikeFactor: DefaultCriticalStrikeFactor + 0.5,
		CooldownReduction:    DefaultCooldownReduction,
		DamageThreatFactor:   DefaultDamageThreatFactor,
		HealingThreatFactor:  DefaultHealingThreatFactor,
		Abilities:            []*Ability{q, w, e, r},
	}
	// Physical damage
	q = &Ability{
		Name:               "Q",
		TargetType:         TargetTypeEnemy,
		HealthCost:         0,
		ManaCost:           0,
		ActivationDuration: 0,
		CooldownDuration:   2 * Second,
		DisableTypes: []DisableType{
			DisableTypeStun,
		},
		Perform: func(up UnitPair) {
			NewPhysicalDamage(up, 140).Perform()
		},
	}
	// Physical damage / DoT / Increasing stacks
	w = &Ability{
		Name:               "W",
		TargetType:         TargetTypeEnemy,
		HealthCost:         0,
		ManaCost:           20,
		ActivationDuration: 0,
		CooldownDuration:   8 * Second,
		DisableTypes: []DisableType{
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
	e = &Ability{
		Name:               "E",
		TargetType:         TargetTypeNone,
		HealthCost:         0,
		ManaCost:           40,
		ActivationDuration: 0,
		CooldownDuration:   20 * Second,
		DisableTypes: []DisableType{
			DisableTypeStun,
		},
		Perform: func(up UnitPair) {
			up.AttachHandler(NewCorrector(
				up.Subject(),
				UnitCorrection{
					Armor:           -25,
					MagicResistance: -25,
				},
				e.Name,
				1,
				8*Second,
			))
			for i := 0; i < 2; i++ {
				up.AttachHandler(newAssassinStack(up.Subject()))
			}
		},
	}
	// Physical
	r = &Ability{
		Name:               "R",
		TargetType:         TargetTypeEnemy,
		HealthCost:         0,
		ManaCost:           120,
		ActivationDuration: 0,
		CooldownDuration:   60 * Second,
		DisableTypes: []DisableType{
			DisableTypeStun,
		},
		Perform: func(up UnitPair) {
			stack := Statistic(0)
			up.ForSubjectHandler(func(ha Handler) {
				switch ha := ha.(type) {
				case *Corrector:
					if ha.name == assassinStack {
						stack += Statistic(ha.Stack())
					}
				}
			})
			NewPhysicalDamage(up, 400*stack*100).Perform()
			up.ForSubjectHandler(func(ha Handler) {
				switch ha := ha.(type) {
				case *Corrector:
					if ha.name == assassinStack {
						up.DetachHandler(ha)
					}
				}
			})
		},
	}
	return Class
}
