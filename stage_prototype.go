package main

type StagePrototype struct {
	prototype             *Unit
	prototypePhase        int
	prototypeCooldownTime InstanceTime

	iron            *Unit
	ironInitialized bool
	ironLeavingTime *InstanceTime

	silver            *Unit
	silverInitialized bool
	silverLeavingTime *InstanceTime
}

func NewStagePrototype() Stage {
	return &StagePrototype{}
}

func (s *StagePrototype) Initialize(g Game) error {
	id, err := g.Join(UnitGroupAI, "P-0", NewClassStagePrototype())
	if err != nil {
		return err
	}
	s.prototype = g.Units().Find(id)
	s.prototypeCooldownTime = g.Clock().Now()
	return nil
}

func (s *StagePrototype) OnTick(g Game) {
	s.syncMinion(g)
	s.actPrototype(g)
	s.actIron(g)
	s.actSilver(g)
}

func (s *StagePrototype) syncMinion(g Game) {
	g.Units().EachFriend(s.prototype, func(u *Unit) {
		switch u.ClassName() {
		case ClassName("Iron"):
			s.iron = u
		case ClassName("Silver"):
			s.silver = u
		default:

		}
	})

	if s.ironLeavingTime == nil && s.iron != nil && s.iron.IsDead() {
		s.ironLeavingTime = new(InstanceTime)
		*s.ironLeavingTime = g.Clock().Add(10 * Second)
	}

	if s.ironLeavingTime != nil && g.Clock().After(*s.ironLeavingTime) {
		g.Leave(s.iron.ID())
		s.iron = nil
		s.ironInitialized = false
		s.ironLeavingTime = nil
	}

	if s.silverLeavingTime == nil && s.silver != nil && s.silver.IsDead() {
		s.silverLeavingTime = new(InstanceTime)
		*s.silverLeavingTime = g.Clock().Add(10 * Second)
	}

	if s.silverLeavingTime != nil && g.Clock().After(*s.silverLeavingTime) {
		g.Leave(s.silver.ID())
		s.silver = nil
		s.silverInitialized = false
		s.silverLeavingTime = nil
	}
}

func (s *StagePrototype) actPrototype(g Game) {
	switch {
	case g.Clock().Before(s.prototypeCooldownTime):
		return
	case s.isActivating(g, s.prototype):
		return
	case s.prototype.Health() == s.prototype.HealthMax():
		return
	case s.prototype.Health() > s.prototype.HealthMax()*0.5:
		if s.prototypePhase != 1 {
			s.prototypePhase = 1
			g.Cooldown(s.prototype, s.prototype.Ability("Attack"))
			g.Cooldown(s.prototype, s.prototype.Ability("Falcon"))
			g.Cooldown(s.prototype, s.prototype.Ability("Shark"))
			g.Cooldown(s.prototype, s.prototype.Ability("Iron"))
		}
		s.actPrototypePhase1(g)
	default:
		if s.prototypePhase != 2 {
			s.prototypePhase = 2
			g.Cooldown(s.prototype, s.prototype.Ability("Attack"))
			g.Cooldown(s.prototype, s.prototype.Ability("Ray"))
			g.Cooldown(s.prototype, s.prototype.Ability("Bell"))
			g.Cooldown(s.prototype, s.prototype.Ability("Silver"))
		}
		s.actPrototypePhase2(g)
	}
}

func (s *StagePrototype) actPrototypePhase1(g Game) {
	o := s.maxThreatEnemy(g)
	if o == nil {
		return
	}
	if s.iron == nil {
		g.Activating(s.prototype, nil, s.prototype.Ability("Iron"))
	}
	g.Activating(s.prototype, o, s.prototype.Ability("Shark"))
	g.Activating(s.prototype, o, s.prototype.Ability("Falcon"))
	g.Activating(s.prototype, o, s.prototype.Ability("Attack"))
	s.prototypeCooldownTime = g.Clock().Add(3 * Second)
}

func (s *StagePrototype) actPrototypePhase2(g Game) {
	o := s.maxThreatEnemy(g)
	if o == nil {
		return
	}
	if s.silver == nil {
		g.Activating(s.prototype, nil, s.prototype.Ability("Silver"))
	}
	if s.iron == nil {
		g.Activating(s.prototype, nil, s.prototype.Ability("Iron"))
	}
	g.Activating(s.prototype, nil, s.prototype.Ability("Bell"))
	g.Activating(s.prototype, o, s.prototype.Ability("Ray"))
	g.Activating(s.prototype, o, s.prototype.Ability("Shark"))
	g.Activating(s.prototype, o, s.prototype.Ability("Falcon"))
	g.Activating(s.prototype, o, s.prototype.Ability("Attack"))
	s.prototypeCooldownTime = g.Clock().Add(2 * Second)
}

func (s *StagePrototype) actIron(g Game) {
	if s.iron == nil {
		return
	}
	if !s.ironInitialized {
		g.Cooldown(s.iron, s.iron.Ability("Silence"))
		g.Cooldown(s.iron, s.iron.Ability("Iron"))
		s.ironInitialized = true
		return
	}
	var o *Unit
	g.Units().EachEnemy(s.iron, func(u *Unit) {
		if o == nil || u.ClassName() == "Healer" {
			o = u
		}
	})
	g.Activating(s.iron, o, s.iron.Ability("Silence"))
	g.Activating(s.iron, o, s.iron.Ability("Iron"))
}

func (s *StagePrototype) actSilver(g Game) {
	if s.silver == nil {
		return
	}
	if !s.silverInitialized {
		g.Cooldown(s.silver, s.silver.Ability("Stun"))
		g.Cooldown(s.silver, s.silver.Ability("Silver"))
		s.silverInitialized = true
		return
	}
	var o *Unit
	g.Units().EachEnemy(s.silver, func(u *Unit) {
		if o == nil || u.ClassName() == "Tank" {
			o = u
		}
	})
	g.Activating(s.silver, o, s.silver.Ability("Stun"))
	g.Activating(s.silver, o, s.silver.Ability("Silver"))
}

func (s *StagePrototype) isActivating(g Game, u *Unit) bool {
	return g.Effects().BindSubject(u).Some(func(h Effect) bool {
		switch h.(type) {
		case *Activating:
			return true
		}
		return false
	})
}

func (s *StagePrototype) maxThreatEnemy(g Game) *Unit {
	var u *Unit
	var threat Statistic
	g.Effects().BindObject(s.prototype).Each(func(h Effect) {
		switch h := h.(type) {
		case *Threat:
			if h.Subject().IsDead() {
				return
			}
			if u == nil || h.threat > threat {
				u = h.Subject()
				threat = h.threat
			}
		}
	})
	return u
}

func NewClassStagePrototype() (class *Class) {
	var attack, falcon, shark, iron, ray, bell, silver, dandelion, tidalBore, diamond, pastorale, vines, waveCrest, morningLull, jadeite Ability
	class = &Class{
		Name:                 "Prototype",
		Health:               40000,
		HealthRegeneration:   0,
		Mana:                 1000,
		ManaRegeneration:     0,
		Armor:                DefaultArmor,
		MagicResistance:      DefaultMagicResistance,
		CriticalStrikeChance: DefaultCriticalStrikeChance,
		CriticalStrikeFactor: DefaultCriticalStrikeFactor,
		DamageThreatFactor:   DefaultDamageThreatFactor,
		HealingThreatFactor:  DefaultHealingThreatFactor,
		Abilities:            []*Ability{&attack, &falcon, &shark, &iron, &ray, &bell, &silver, &dandelion, &tidalBore, &diamond, &pastorale, &vines, &waveCrest, &morningLull, &jadeite},
	}
	attack = Ability{
		Name:               "Attack",
		TargetType:         TargetTypeEnemy,
		ActivationDuration: 0 * Second,
		CooldownDuration:   4 * Second,
		DisableTypes: []DisableType{
			DisableTypeStun,
		},
		Perform: func(g Game, s Subject, o *Unit) {
			baseDamage := 100 + 100*(1-s.Subject().Health()/s.Subject().HealthMax())
			g.PhysicalDamage(s, o, baseDamage).Perform()
		},
	}
	falcon = Ability{
		Name:               "Falcon",
		TargetType:         TargetTypeEnemy,
		ActivationDuration: 1 * Second,
		CooldownDuration:   8 * Second,
		DisableTypes: []DisableType{
			DisableTypeStun,
		},
		Perform: func(g Game, s Subject, o *Unit) {
			baseDamage := Statistic(160)
			g.PhysicalDamage(s, o, baseDamage).Perform()
		},
	}
	shark = Ability{
		Name:               "Shark",
		TargetType:         TargetTypeEnemy,
		ActivationDuration: 1 * Second,
		CooldownDuration:   12 * Second,
		DisableTypes: []DisableType{
			DisableTypeStun,
		},
		Perform: func(g Game, s Subject, o *Unit) {
			baseDamage := Statistic(310)
			g.PhysicalDamage(s, o, baseDamage).Perform()
		},
	}
	iron = Ability{
		Name:               "Iron",
		TargetType:         TargetTypeNone,
		ActivationDuration: 1 * Second,
		CooldownDuration:   30 * Second,
		DisableTypes:       []DisableType{},
		Perform: func(g Game, s Subject, o *Unit) {
			g.Join(s.Subject().Group(), "Iron-PX1", NewClassStageIron())
		},
	}
	ray = Ability{
		Name:               "Ray",
		TargetType:         TargetTypeEnemy,
		ActivationDuration: 1 * Second,
		CooldownDuration:   20 * Second,
		DisableTypes: []DisableType{
			DisableTypeSilence,
		},
		Perform: func(g Game, s Subject, o *Unit) {
			g.MagicDamage(s, o, 100).Perform()
			if o.IsDead() {
				return
			}
			g.DoT(g.MagicDamage(s, o, 35), 10*Second, ray.Name)
		},
	}
	bell = Ability{
		Name:               "Bell",
		TargetType:         TargetTypeNone,
		ActivationDuration: 1 * Second,
		CooldownDuration:   30 * Second,
		DisableTypes: []DisableType{
			DisableTypeSilence,
		},
		Perform: func(g Game, s Subject, o *Unit) {
			g.Units().EachEnemy(s.Subject(), func(enemy *Unit) {
				if enemy.IsDead() {
					return
				}
				g.MagicDamage(s, enemy, 180).Perform()
			})
		},
	}
	silver = Ability{
		Name:               "Silver",
		TargetType:         TargetTypeNone,
		ActivationDuration: 1 * Second,
		CooldownDuration:   40 * Second,
		DisableTypes:       []DisableType{},
		Perform: func(g Game, s Subject, o *Unit) {
			g.Join(s.Subject().Group(), "Silver-PX2", NewClassStageSilver())
		},
	}
	return
}

func NewClassStageIron() (class *Class) {
	var silence, iron Ability
	class = &Class{
		Name:                 "Iron",
		Health:               500,
		HealthRegeneration:   50,
		Mana:                 100,
		ManaRegeneration:     0,
		Armor:                300,
		MagicResistance:      DefaultMagicResistance,
		CriticalStrikeChance: DefaultCriticalStrikeChance,
		CriticalStrikeFactor: DefaultCriticalStrikeFactor,
		DamageThreatFactor:   DefaultDamageThreatFactor,
		HealingThreatFactor:  DefaultHealingThreatFactor,
		Abilities:            []*Ability{&silence, &iron},
	}
	silence = Ability{
		Name:               "Stun",
		TargetType:         TargetTypeEnemy,
		ActivationDuration: 1 * Second,
		CooldownDuration:   8 * Second,
		DisableTypes:       []DisableType{},
		Perform: func(g Game, s Subject, o *Unit) {
			g.Disable(o, DisableTypeStun, 3*Second)
		},
	}
	iron = Ability{
		Name:               "Iron",
		TargetType:         TargetTypeEnemy,
		ActivationDuration: 1 * Second,
		CooldownDuration:   5 * Second,
		DisableTypes:       []DisableType{},
		Perform: func(g Game, s Subject, o *Unit) {
			g.MagicDamage(s, o, 120).Perform()
		},
	}
	return class
}

func NewClassStageSilver() (class *Class) {
	var stun, silver Ability
	class = &Class{
		Name:                 "Silver",
		Health:               1000,
		HealthRegeneration:   100,
		Mana:                 100,
		ManaRegeneration:     0,
		Armor:                DefaultArmor,
		MagicResistance:      300,
		CriticalStrikeChance: DefaultCriticalStrikeChance,
		CriticalStrikeFactor: DefaultCriticalStrikeFactor,
		DamageThreatFactor:   DefaultDamageThreatFactor,
		HealingThreatFactor:  DefaultHealingThreatFactor,
		Abilities:            []*Ability{&stun, &silver},
	}
	stun = Ability{
		Name:               "Silence",
		TargetType:         TargetTypeEnemy,
		ActivationDuration: 1 * Second,
		CooldownDuration:   8 * Second,
		DisableTypes:       []DisableType{},
		Perform: func(g Game, s Subject, o *Unit) {
			g.Disable(o, DisableTypeSilence, 3*Second)
		},
	}
	silver = Ability{
		Name:               "Silver",
		TargetType:         TargetTypeEnemy,
		ActivationDuration: 1 * Second,
		CooldownDuration:   5 * Second,
		DisableTypes:       []DisableType{},
		Perform: func(g Game, s Subject, o *Unit) {
			g.MagicDamage(s, o, 270).Perform()
		},
	}
	return class
}
