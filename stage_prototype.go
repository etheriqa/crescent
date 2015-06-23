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

func (s *StagePrototype) Initialize(op Operator) error {
	id, err := op.Join(UnitGroupAI, "P-0", NewClassStagePrototype())
	if err != nil {
		return err
	}
	s.prototype = op.Units().Find(id)
	s.prototypeCooldownTime = op.Clock().Now()
	return nil
}

func (s *StagePrototype) OnTick(op Operator) {
	s.syncMinion(op)
	s.actPrototype(op)
	s.actIron(op)
	s.actSilver(op)
}

func (s *StagePrototype) syncMinion(op Operator) {
	op.Units().EachFriend(s.prototype, func(u *Unit) {
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
		*s.ironLeavingTime = op.Clock().Add(10 * Second)
	}

	if s.ironLeavingTime != nil && op.Clock().After(*s.ironLeavingTime) {
		op.Leave(s.iron.ID())
		s.iron = nil
		s.ironInitialized = false
		s.ironLeavingTime = nil
	}

	if s.silverLeavingTime == nil && s.silver != nil && s.silver.IsDead() {
		s.silverLeavingTime = new(InstanceTime)
		*s.silverLeavingTime = op.Clock().Add(10 * Second)
	}

	if s.silverLeavingTime != nil && op.Clock().After(*s.silverLeavingTime) {
		op.Leave(s.silver.ID())
		s.silver = nil
		s.silverInitialized = false
		s.silverLeavingTime = nil
	}
}

func (s *StagePrototype) actPrototype(op Operator) {
	switch {
	case op.Clock().Before(s.prototypeCooldownTime):
		return
	case s.isActivating(op, s.prototype):
		return
	case s.prototype.Health() == s.prototype.HealthMax():
		return
	case s.prototype.Health() > s.prototype.HealthMax()*0.5:
		if s.prototypePhase != 1 {
			s.prototypePhase = 1
			op.Cooldown(s.prototype, s.prototype.Ability("Attack"))
			op.Cooldown(s.prototype, s.prototype.Ability("Falcon"))
			op.Cooldown(s.prototype, s.prototype.Ability("Shark"))
			op.Cooldown(s.prototype, s.prototype.Ability("Iron"))
		}
		s.actPrototypePhase1(op)
	default:
		if s.prototypePhase != 2 {
			s.prototypePhase = 2
			op.Cooldown(s.prototype, s.prototype.Ability("Attack"))
			op.Cooldown(s.prototype, s.prototype.Ability("Ray"))
			op.Cooldown(s.prototype, s.prototype.Ability("Bell"))
			op.Cooldown(s.prototype, s.prototype.Ability("Silver"))
		}
		s.actPrototypePhase2(op)
	}
}

func (s *StagePrototype) actPrototypePhase1(op Operator) {
	o := s.maxThreatEnemy(op)
	if o == nil {
		return
	}
	if s.iron == nil {
		op.Activating(s.prototype, nil, s.prototype.Ability("Iron"))
	}
	op.Activating(s.prototype, o, s.prototype.Ability("Shark"))
	op.Activating(s.prototype, o, s.prototype.Ability("Falcon"))
	op.Activating(s.prototype, o, s.prototype.Ability("Attack"))
	s.prototypeCooldownTime = op.Clock().Add(3 * Second)
}

func (s *StagePrototype) actPrototypePhase2(op Operator) {
	o := s.maxThreatEnemy(op)
	if o == nil {
		return
	}
	if s.silver == nil {
		op.Activating(s.prototype, nil, s.prototype.Ability("Silver"))
	}
	if s.iron == nil {
		op.Activating(s.prototype, nil, s.prototype.Ability("Iron"))
	}
	op.Activating(s.prototype, nil, s.prototype.Ability("Bell"))
	op.Activating(s.prototype, o, s.prototype.Ability("Ray"))
	op.Activating(s.prototype, o, s.prototype.Ability("Shark"))
	op.Activating(s.prototype, o, s.prototype.Ability("Falcon"))
	op.Activating(s.prototype, o, s.prototype.Ability("Attack"))
	s.prototypeCooldownTime = op.Clock().Add(2 * Second)
}

func (s *StagePrototype) actIron(op Operator) {
	if s.iron == nil {
		return
	}
	if !s.ironInitialized {
		op.Cooldown(s.iron, s.iron.Ability("Silence"))
		op.Cooldown(s.iron, s.iron.Ability("Iron"))
		s.ironInitialized = true
		return
	}
	var o *Unit
	op.Units().EachEnemy(s.iron, func(u *Unit) {
		if o == nil || u.ClassName() == "Healer" {
			o = u
		}
	})
	op.Activating(s.iron, o, s.iron.Ability("Silence"))
	op.Activating(s.iron, o, s.iron.Ability("Iron"))
}

func (s *StagePrototype) actSilver(op Operator) {
	if s.silver == nil {
		return
	}
	if !s.silverInitialized {
		op.Cooldown(s.silver, s.silver.Ability("Stun"))
		op.Cooldown(s.silver, s.silver.Ability("Silver"))
		s.silverInitialized = true
		return
	}
	var o *Unit
	op.Units().EachEnemy(s.silver, func(u *Unit) {
		if o == nil || u.ClassName() == "Tank" {
			o = u
		}
	})
	op.Activating(s.silver, o, s.silver.Ability("Stun"))
	op.Activating(s.silver, o, s.silver.Ability("Silver"))
}

func (s *StagePrototype) isActivating(op Operator, u *Unit) bool {
	return op.Handlers().BindSubject(u).Some(func(h Handler) bool {
		switch h.(type) {
		case *Activating:
			return true
		}
		return false
	})
}

func (s *StagePrototype) maxThreatEnemy(op Operator) *Unit {
	var u *Unit
	var threat Statistic
	op.Handlers().BindObject(s.prototype).Each(func(h Handler) {
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
		Perform: func(op Operator, s Subject, o *Unit) {
			baseDamage := 100 + 100*(1-s.Subject().Health()/s.Subject().HealthMax())
			op.PhysicalDamage(s, o, baseDamage).Perform()
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
		Perform: func(op Operator, s Subject, o *Unit) {
			baseDamage := Statistic(160)
			op.PhysicalDamage(s, o, baseDamage).Perform()
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
		Perform: func(op Operator, s Subject, o *Unit) {
			baseDamage := Statistic(310)
			op.PhysicalDamage(s, o, baseDamage).Perform()
		},
	}
	iron = Ability{
		Name:               "Iron",
		TargetType:         TargetTypeNone,
		ActivationDuration: 1 * Second,
		CooldownDuration:   30 * Second,
		DisableTypes:       []DisableType{},
		Perform: func(op Operator, s Subject, o *Unit) {
			op.Join(s.Subject().Group(), "Iron-PX1", NewClassStageIron())
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
		Perform: func(op Operator, s Subject, o *Unit) {
			op.MagicDamage(s, o, 100).Perform()
			if o.IsDead() {
				return
			}
			op.DoT(op.MagicDamage(s, o, 35), 10*Second, ray.Name)
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
		Perform: func(op Operator, s Subject, o *Unit) {
			op.Units().EachEnemy(s.Subject(), func(enemy *Unit) {
				if enemy.IsDead() {
					return
				}
				op.MagicDamage(s, enemy, 180).Perform()
			})
		},
	}
	silver = Ability{
		Name:               "Silver",
		TargetType:         TargetTypeNone,
		ActivationDuration: 1 * Second,
		CooldownDuration:   40 * Second,
		DisableTypes:       []DisableType{},
		Perform: func(op Operator, s Subject, o *Unit) {
			op.Join(s.Subject().Group(), "Silver-PX2", NewClassStageSilver())
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
		Perform: func(op Operator, s Subject, o *Unit) {
			op.Disable(o, DisableTypeStun, 3*Second)
		},
	}
	iron = Ability{
		Name:               "Iron",
		TargetType:         TargetTypeEnemy,
		ActivationDuration: 1 * Second,
		CooldownDuration:   5 * Second,
		DisableTypes:       []DisableType{},
		Perform: func(op Operator, s Subject, o *Unit) {
			op.MagicDamage(s, o, 120).Perform()
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
		Perform: func(op Operator, s Subject, o *Unit) {
			op.Disable(o, DisableTypeSilence, 3*Second)
		},
	}
	silver = Ability{
		Name:               "Silver",
		TargetType:         TargetTypeEnemy,
		ActivationDuration: 1 * Second,
		CooldownDuration:   5 * Second,
		DisableTypes:       []DisableType{},
		Perform: func(op Operator, s Subject, o *Unit) {
			op.MagicDamage(s, o, 270).Perform()
		},
	}
	return class
}
