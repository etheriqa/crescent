package main

const AssassinStackName = "Assassin Stack"

func AssassinStack(op Operator, o Object) {
	c := UnitCorrection{
		CriticalStrikeChance: 0.05,
		CriticalStrikeFactor: 0.1,
	}
	op.Correction(o, c, 10, 10*Second, AssassinStackName)
}

func NewClassAssassin() *Class {
	var q, w, e, r *Ability
	class := &Class{
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
		Perform: func(op Operator, s Subject, o *Unit) {
			_, _, _, err := op.PhysicalDamage(s, o, 140).Perform()
			if err != nil {
				log.Fatal(err)
			}
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
		Perform: func(op Operator, s Subject, o *Unit) {
			_, _, _, err := op.PhysicalDamage(s, o, 80).Perform()
			if err != nil {
				log.Fatal(err)
			}
			op.DoT(op.PhysicalDamage(s, o, 10), 10*Second, w.Name)
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
		Perform: func(op Operator, s Subject, o *Unit) {
			c := UnitCorrection{
				Armor:           -25,
				MagicResistance: -25,
			}
			op.Correction(s.Subject(), c, 1, 8*Second, e.Name)
			for i := 0; i < 2; i++ {
				AssassinStack(op, MakeObject(s.Subject()))
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
		Perform: func(op Operator, s Subject, o *Unit) {
			stack := Statistic(0)
			op.Handlers().Each(func(h Handler) {
				switch h := h.(type) {
				case *Correction:
					if h.name == AssassinStackName {
						stack += h.Stack()
					}
				}
			})
			_, _, _, err := op.PhysicalDamage(s, o, 400+stack*100).Perform()
			if err != nil {
				log.Fatal(err)
			}
			op.Handlers().Each(func(h Handler) {
				switch h := h.(type) {
				case *Correction:
					if h.name == AssassinStackName {
						op.Handlers().Detach(h)
					}
				}
			})
		},
	}
	return class
}
