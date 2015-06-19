package main

const AssassinStackName = "Tenecity"

func AssassinStack(op Operator, o Object) {
	c := UnitCorrection{
		CriticalStrikeChance: 0.05,
		CriticalStrikeFactor: 0.05,
	}
	op.Correction(o, c, 10, 10*Second, AssassinStackName)
}

func NewClassAssassin() *Class {
	var q, w, e, r Ability
	class := &Class{
		Name:                 "Assassin",
		Health:               600,
		HealthRegeneration:   11,
		Mana:                 200,
		ManaRegeneration:     14,
		Armor:                DefaultArmor,
		MagicResistance:      DefaultMagicResistance,
		CriticalStrikeChance: DefaultCriticalStrikeChance + 0.05,
		CriticalStrikeFactor: DefaultCriticalStrikeFactor + 0.5,
		CooldownReduction:    DefaultCooldownReduction,
		DamageThreatFactor:   DefaultDamageThreatFactor,
		HealingThreatFactor:  DefaultHealingThreatFactor,
		Abilities:            []*Ability{&q, &w, &e, &r},
	}
	q = Ability{
		Name:               "Triple Cleave",
		Description:        "Deals physical damage three times",
		TargetType:         TargetTypeEnemy,
		HealthCost:         0,
		ManaCost:           0,
		ActivationDuration: 0,
		CooldownDuration:   2 * Second,
		DisableTypes: []DisableType{
			DisableTypeStun,
		},
		Perform: func(op Operator, s Subject, o *Unit) {
			for i := 0; i < 3; i++ {
				_, _, crit, err := op.PhysicalDamage(s, o, 45).Perform()
				if err != nil {
					log.Fatal(err)
				}
				if crit {
					AssassinStack(op, MakeObject(s.Subject()))
				}
			}
		},
	}
	w = Ability{
		Name:               "Poison Dart",
		Description:        "Deals physcical damage / Grants a physical damage over time effect for 10 seconds to target",
		TargetType:         TargetTypeEnemy,
		HealthCost:         0,
		ManaCost:           20,
		ActivationDuration: 0,
		CooldownDuration:   8 * Second,
		DisableTypes: []DisableType{
			DisableTypeStun,
		},
		Perform: func(op Operator, s Subject, o *Unit) {
			_, _, crit, err := op.PhysicalDamage(s, o, 80).Perform()
			if err != nil {
				log.Fatal(err)
			}
			if crit {
				AssassinStack(op, MakeObject(s.Subject()))
			}
			op.DoT(op.PhysicalDamage(s, o, 10), 10*Second, w.Name)
		},
	}
	e = Ability{
		Name:               "Tenacity",
		Description:        "Gains three Tenacity effects / Loses armor and magic resistance for 8 seconds / Tenecity effect increases critical strike chance and critical strike damage",
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
			for i := 0; i < 3; i++ {
				AssassinStack(op, MakeObject(s.Subject()))
			}
		},
	}
	w = Ability{
		Name:               "Lethal Weapon",
		Description:        "Deals pure damage / Silences target for 1 second / Consumes all Tenecity effects",
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
			_, _, _, err := op.PureDamage(s, o, 400+stack*100).Perform()
			if err != nil {
				log.Fatal(err)
			}
			op.Disable(o, DisableTypeSilence, Second)
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
