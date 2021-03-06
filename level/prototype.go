package level

import (
	. "github.com/etheriqa/crescent/game"
)

type LevelPrototype struct {
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

func NewLevelPrototype() Level {
	return &LevelPrototype{}
}

func (s *LevelPrototype) Initialize(g Game) error {
	id, err := g.Join(UnitGroupAI, "P-0", NewClassLevelPrototype())
	if err != nil {
		return err
	}
	s.prototype = g.UnitQuery().Find(id)
	s.prototypeCooldownTime = g.Clock().Now()
	return nil
}

func (s *LevelPrototype) OnTick(g Game) {
	s.syncMinion(g)
	s.actPrototype(g)
	s.actIron(g)
	s.actSilver(g)
}

func (s *LevelPrototype) syncMinion(g Game) {
	g.UnitQuery().EachFriend(s.prototype, func(u *Unit) {
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

func (s *LevelPrototype) actPrototype(g Game) {
	switch {
	case g.Clock().Before(s.prototypeCooldownTime):
		return
	case s.isActivating(g, s.prototype):
		return
	case s.prototype.Health() == s.prototype.HealthMax():
		return
	case s.prototype.Health() > s.prototype.HealthMax()*0.6:
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

func (s *LevelPrototype) actPrototypePhase1(g Game) {
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

func (s *LevelPrototype) actPrototypePhase2(g Game) {
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

func (s *LevelPrototype) actIron(g Game) {
	if s.iron == nil {
		return
	}
	if !s.ironInitialized {
		g.Cooldown(s.iron, s.iron.Ability("Stun"))
		g.Cooldown(s.iron, s.iron.Ability("Iron"))
		s.ironInitialized = true
		return
	}
	var o *Unit
	g.UnitQuery().EachEnemy(s.iron, func(u *Unit) {
		if o == nil || u.ClassName() == "Healer" {
			o = u
		}
	})
	g.Activating(s.iron, o, s.iron.Ability("Stun"))
	g.Activating(s.iron, o, s.iron.Ability("Iron"))
}

func (s *LevelPrototype) actSilver(g Game) {
	if s.silver == nil {
		return
	}
	if !s.silverInitialized {
		g.Cooldown(s.silver, s.silver.Ability("Silence"))
		g.Cooldown(s.silver, s.silver.Ability("Silver"))
		s.silverInitialized = true
		return
	}
	var o *Unit
	g.UnitQuery().EachEnemy(s.silver, func(u *Unit) {
		if o == nil || u.ClassName() == "Tank" {
			o = u
		}
	})
	g.Activating(s.silver, o, s.silver.Ability("Silence"))
	g.Activating(s.silver, o, s.silver.Ability("Silver"))
}

func (s *LevelPrototype) isActivating(g Game, u *Unit) bool {
	return g.EffectQuery().BindSubject(u).Some(func(h Effect) bool {
		switch h.(type) {
		case *Activating:
			return true
		}
		return false
	})
}

func (s *LevelPrototype) maxThreatEnemy(g Game) *Unit {
	var u *Unit
	var threat Statistic
	g.EffectQuery().BindObject(s.prototype).Each(func(h Effect) {
		switch h := h.(type) {
		case *Threat:
			if h.Subject().IsDead() {
				return
			}
			if u == nil || h.Threat() > threat {
				u = h.Subject()
				threat = h.Threat()
			}
		}
	})
	return u
}

func NewClassLevelPrototype() (class *Class) {
	var attack, falcon, shark, iron, ray, bell, silver, dandelion, tidalBore, diamond, pastorale, vines, waveCrest, morningLull, jadeite Ability
	class = &Class{
		Name:                 "Prototype",
		Health:               50000,
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
			MakePhysicalDamage(s, o, baseDamage).Perform(g)
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
			MakePhysicalDamage(s, o, baseDamage).Perform(g)
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
			MakePhysicalDamage(s, o, baseDamage).Perform(g)
		},
	}
	iron = Ability{
		Name:               "Iron",
		TargetType:         TargetTypeNone,
		ActivationDuration: 1 * Second,
		CooldownDuration:   30 * Second,
		DisableTypes:       []DisableType{},
		Perform: func(g Game, s Subject, o *Unit) {
			g.Join(s.Subject().Group(), "Iron-PX1", NewClassLevelIron())
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
			MakeMagicDamage(s, o, 100).Perform(g)
			if o.IsDead() {
				return
			}
			g.DoT(MakeMagicDamage(s, o, 45), ray.Name, 10*Second)
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
			g.UnitQuery().EachEnemy(s.Subject(), func(enemy *Unit) {
				if enemy.IsDead() {
					return
				}
				MakeMagicDamage(s, enemy, 280).Perform(g)
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
			g.Join(s.Subject().Group(), "Silver-PX2", NewClassLevelSilver())
		},
	}
	return
}

func NewClassLevelIron() (class *Class) {
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
			MakeMagicDamage(s, o, 120).Perform(g)
		},
	}
	return class
}

func NewClassLevelSilver() (class *Class) {
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
			MakeMagicDamage(s, o, 270).Perform(g)
		},
	}
	return class
}
