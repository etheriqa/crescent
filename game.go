package main

type Game interface {
	Clock() InstanceClock
	Effects() EffectContainer
	Units() UnitContainer

	Writer() InstanceOutputWriter

	Join(UnitGroup, UnitName, *Class) (UnitID, error)
	Leave(UnitID) error

	Activating(Subject, *Unit, *Ability)
	Cooldown(Object, *Ability)
	ResetCooldown(Object, *Ability)
	Correction(Object, UnitCorrection, Statistic, InstanceDuration, string)
	Disable(Object, DisableType, InstanceDuration)
	DamageThreat(Subject, Object, Statistic)
	HealingThreat(Subject, Object, Statistic)
	DoT(*Damage, InstanceDuration, string)
	HoT(*Healing, InstanceDuration, string)

	PhysicalDamage(Subject, Object, Statistic) *Damage
	MagicDamage(Subject, Object, Statistic) *Damage
	TrueDamage(Subject, Object, Statistic) *Damage
	PureDamage(Subject, Object, Statistic) *Damage

	Healing(Subject, Object, Statistic) *Healing
}

// Clock returns the InstanceClock
func (g *GameState) Clock() InstanceClock {
	return g.clock
}

// Effects returns the EffectContainer {
func (g *GameState) Effects() EffectContainer {
	return g.effects
}

// Units returns the UnitContainer {
func (g *GameState) Units() UnitContainer {
	return g.units
}

// Writer returns the InstanceOutputWriter
func (g *GameState) Writer() InstanceOutputWriter {
	return g.w
}

// Activating attaches a Activating Effect
func (g *GameState) Activating(s Subject, o *Unit, a *Ability) {
	g.effects.Attach(&Activating{
		UnitSubject:    MakeSubject(s),
		object:         o,
		ability:        a,
		expirationTime: g.clock.Add(a.ActivationDuration),

		g: g,
	})
}

// Cooldown attaches a Cooldown Effect
func (g *GameState) Cooldown(o Object, a *Ability) {
	g.effects.Attach(&Cooldown{
		UnitObject:     MakeObject(o),
		ability:        a,
		expirationTime: g.clock.Add(a.CooldownDuration),

		g: g,
	})
}

// ResetCooldown detaches Cooldown effects
func (g *GameState) ResetCooldown(o Object, a *Ability) {
	g.effects.Attach(&Cooldown{
		UnitObject:     MakeObject(o),
		ability:        a,
		expirationTime: g.clock.Now(),

		g: g,
	})
}

// Correction attaches a Correction Effect
func (g *GameState) Correction(o Object, c UnitCorrection, l Statistic, d InstanceDuration, name string) {
	g.effects.Attach(&Correction{
		UnitObject:     MakeObject(o),
		name:           name,
		correction:     c,
		stackLimit:     l,
		stack:          1,
		expirationTime: g.clock.Add(d),

		g: g,
	})
}

// Disable attaches a Disable Effect
func (g *GameState) Disable(o Object, t DisableType, d InstanceDuration) {
	g.effects.Attach(&Disable{
		UnitObject:     MakeObject(o),
		disableType:    t,
		expirationTime: g.clock.Add(d),

		g: g,
	})
}

// DamageThreat attaches a Threat Effect
func (g *GameState) DamageThreat(s Subject, o Object, d Statistic) {
	g.effects.Attach(&Threat{
		UnitPair: MakePair(s, o),
		threat:   d * s.Subject().DamageThreatFactor(),

		g: g,
	})
}

// HealingThreat attaches a Threat Effect
func (g *GameState) HealingThreat(s Subject, o Object, h Statistic) {
	g.effects.Attach(&Threat{
		UnitPair: MakePair(s, o),
		threat:   h * s.Subject().HealingThreatFactor(),

		g: g,
	})
}

// DoT attaches a Periodical Effect
func (g *GameState) DoT(damage *Damage, d InstanceDuration, name string) {
	g.effects.Attach(&Periodical{
		UnitPair:       MakePair(damage, damage),
		name:           name,
		routine:        func() { damage.Perform() },
		expirationTime: g.clock.Add(d),

		g: g,
	})
}

// HoT attaches a Periodical Effect
func (g *GameState) HoT(healing *Healing, d InstanceDuration, name string) {
	g.effects.Attach(&Periodical{
		UnitPair:       MakePair(healing, healing),
		name:           name,
		routine:        func() { healing.Perform() },
		expirationTime: g.clock.Add(d),

		g: g,
	})
}

// PhysicalDamage returns a Damage
func (g *GameState) PhysicalDamage(s Subject, o Object, d Statistic) *Damage {
	return &Damage{
		UnitPair:             MakePair(s, o),
		damage:               d * o.Object().PhysicalDamageReductionFactor(),
		criticalStrikeChance: s.Subject().CriticalStrikeChance(),
		criticalStrikeFactor: s.Subject().CriticalStrikeFactor(),

		g: g,
	}
}

// MagicDamage returns a Damage
func (g *GameState) MagicDamage(s Subject, o Object, d Statistic) *Damage {
	return &Damage{
		UnitPair:             MakePair(s, o),
		damage:               d * o.Object().MagicDamageReductionFactor(),
		criticalStrikeChance: s.Subject().CriticalStrikeChance(),
		criticalStrikeFactor: s.Subject().CriticalStrikeFactor(),

		g: g,
	}
}

// TrueDamage returns a Damage
func (g *GameState) TrueDamage(s Subject, o Object, d Statistic) *Damage {
	return &Damage{
		UnitPair:             MakePair(s, o),
		damage:               d,
		criticalStrikeChance: s.Subject().CriticalStrikeChance(),
		criticalStrikeFactor: s.Subject().CriticalStrikeFactor(),

		g: g,
	}
}

// PureDamage returns a Damage
func (g *GameState) PureDamage(s Subject, o Object, d Statistic) *Damage {
	return &Damage{
		UnitPair:             MakePair(s, o),
		damage:               d,
		criticalStrikeChance: 0,
		criticalStrikeFactor: 0,

		g: g,
	}
}

// Healing returns a Healing
func (g *GameState) Healing(s Subject, o Object, h Statistic) *Healing {
	return &Healing{
		UnitPair:             MakePair(s, o),
		healing:              h,
		criticalStrikeChance: s.Subject().CriticalStrikeChance(),
		criticalStrikeFactor: s.Subject().CriticalStrikeFactor(),

		g: g,
	}
}
