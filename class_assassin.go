package main

const AssassinStackName = "Tenacity"

func AssassinStack(g Game, o Object) {
	c := UnitCorrection{
		CriticalStrikeChance: 0.05,
		CriticalStrikeFactor: 0.05,
	}
	g.Correction(o, c, 10, 10*Second, AssassinStackName)
}

func NewClassAssassin() *Class {
	var q, w, e, r Ability
	class := &Class{
		Name:                 "Assassin",
		Health:               880,
		HealthRegeneration:   21,
		Mana:                 280,
		ManaRegeneration:     15,
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
		Perform: func(g Game, s Subject, o *Unit) {
			for i := 0; i < 3; i++ {
				if o.IsDead() {
					return
				}
				_, _, crit, err := g.PhysicalDamage(s, o, 45).Perform()
				if err != nil {
					log.Fatal(err)
				}
				if crit {
					AssassinStack(g, MakeObject(s.Subject()))
				}
			}
		},
	}
	w = Ability{
		Name:               "Poison Dart",
		Description:        "Deals physcical damage / Grants a physical damage over time effect for 10 seconds to target",
		TargetType:         TargetTypeEnemy,
		HealthCost:         0,
		ManaCost:           15,
		ActivationDuration: 0,
		CooldownDuration:   8 * Second,
		DisableTypes: []DisableType{
			DisableTypeStun,
		},
		Perform: func(g Game, s Subject, o *Unit) {
			_, _, crit, err := g.PhysicalDamage(s, o, 80).Perform()
			if err != nil {
				log.Fatal(err)
			}
			if crit {
				AssassinStack(g, MakeObject(s.Subject()))
			}
			if o.IsDead() {
				return
			}
			g.DoT(g.PhysicalDamage(s, o, 10), 10*Second, w.Name)
		},
	}
	e = Ability{
		Name:               "Tenacity",
		Description:        "Gains three Tenacity effects / Loses armor and magic resistance for 8 seconds / Tenacity effect increases critical strike chance and critical strike damage",
		TargetType:         TargetTypeNone,
		HealthCost:         0,
		ManaCost:           32,
		ActivationDuration: 0,
		CooldownDuration:   20 * Second,
		DisableTypes: []DisableType{
			DisableTypeStun,
		},
		Perform: func(g Game, s Subject, o *Unit) {
			c := UnitCorrection{
				Armor:           -25,
				MagicResistance: -25,
			}
			g.Correction(s.Subject(), c, 1, 8*Second, e.Name)
			for i := 0; i < 3; i++ {
				AssassinStack(g, MakeObject(s.Subject()))
			}
		},
	}
	r = Ability{
		Name:               "Lethal Weapon",
		Description:        "Deals pure damage / Silences target for 1 second / Consumes all Tenacity effects",
		TargetType:         TargetTypeEnemy,
		HealthCost:         0,
		ManaCost:           58,
		ActivationDuration: 0,
		CooldownDuration:   60 * Second,
		DisableTypes: []DisableType{
			DisableTypeStun,
		},
		Perform: func(g Game, s Subject, o *Unit) {
			stack := Statistic(0)
			g.Effects().Each(func(h Effect) {
				switch h := h.(type) {
				case *Correction:
					if h.name == AssassinStackName {
						stack += h.Stack()
					}
				}
			})
			_, _, _, err := g.PureDamage(s, o, 400+stack*100).Perform()
			if err != nil {
				log.Fatal(err)
			}
			g.Effects().Each(func(h Effect) {
				switch h := h.(type) {
				case *Correction:
					if h.name == AssassinStackName {
						g.Effects().Detach(h)
					}
				}
			})
			if o.IsDead() {
				return
			}
			g.Disable(o, DisableTypeSilence, Second)
		},
	}
	return class
}
