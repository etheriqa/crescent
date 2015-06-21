package main

type Operator interface {
	Clock() InstanceClock
	Handlers() HandlerContainer
	Units() UnitContainer

	Writer() InstanceOutputWriter

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
func (g *Game) Clock() InstanceClock {
	return g.clock
}

// Handlers returns the HandlerContainer {
func (g *Game) Handlers() HandlerContainer {
	return g.handlers
}

// Units returns the UnitContainer {
func (g *Game) Units() UnitContainer {
	return g.units
}

// Writer returns the InstanceOutputWriter
func (g *Game) Writer() InstanceOutputWriter {
	return g.w
}

// Activating attaches a Activating Handler
func (g *Game) Activating(s Subject, o *Unit, a *Ability) {
	g.handlers.Attach(&Activating{
		UnitSubject:    MakeSubject(s),
		object:         o,
		ability:        a,
		expirationTime: g.clock.Add(a.ActivationDuration),

		op: g,
	})
}

// Cooldown attaches a Cooldown Handler
func (g *Game) Cooldown(o Object, a *Ability) {
	g.handlers.Attach(&Cooldown{
		UnitObject:     MakeObject(o),
		ability:        a,
		expirationTime: g.clock.Add(a.CooldownDuration),

		op: g,
	})
}

// ResetCooldown detaches Cooldown handlers
func (g *Game) ResetCooldown(o Object, a *Ability) {
	g.handlers.Attach(&Cooldown{
		UnitObject:     MakeObject(o),
		ability:        a,
		expirationTime: g.clock.Now(),

		op: g,
	})
}

// Correction attaches a Correction Handler
func (g *Game) Correction(o Object, c UnitCorrection, l Statistic, d InstanceDuration, name string) {
	g.handlers.Attach(&Correction{
		UnitObject:     MakeObject(o),
		name:           name,
		correction:     c,
		stackLimit:     l,
		stack:          1,
		expirationTime: g.clock.Add(d),

		op: g,
	})
}

// Disable attaches a Disable Handler
func (g *Game) Disable(o Object, t DisableType, d InstanceDuration) {
	g.handlers.Attach(&Disable{
		UnitObject:     MakeObject(o),
		disableType:    t,
		expirationTime: g.clock.Add(d),

		op: g,
	})
}

// DamageThreat attaches a Threat Handler
func (g *Game) DamageThreat(s Subject, o Object, d Statistic) {
	g.handlers.Attach(&Threat{
		UnitPair: MakePair(s, o),
		threat:   d * s.Subject().DamageThreatFactor(),

		op: g,
	})
}

// HealingThreat attaches a Threat Handler
func (g *Game) HealingThreat(s Subject, o Object, h Statistic) {
	g.handlers.Attach(&Threat{
		UnitPair: MakePair(s, o),
		threat:   h * s.Subject().HealingThreatFactor(),

		op: g,
	})
}

// DoT attaches a Periodical Handler
func (g *Game) DoT(damage *Damage, d InstanceDuration, name string) {
	g.handlers.Attach(&Periodical{
		UnitPair:       MakePair(damage, damage),
		name:           name,
		routine:        func() { damage.Perform() },
		expirationTime: g.clock.Add(d),

		op: g,
	})
}

// HoT attaches a Periodical Handler
func (g *Game) HoT(healing *Healing, d InstanceDuration, name string) {
	g.handlers.Attach(&Periodical{
		UnitPair:       MakePair(healing, healing),
		name:           name,
		routine:        func() { healing.Perform() },
		expirationTime: g.clock.Add(d),

		op: g,
	})
}

// PhysicalDamage returns a Damage
func (g *Game) PhysicalDamage(s Subject, o Object, d Statistic) *Damage {
	return &Damage{
		UnitPair:             MakePair(s, o),
		damage:               d * o.Object().PhysicalDamageReductionFactor(),
		criticalStrikeChance: s.Subject().CriticalStrikeChance(),
		criticalStrikeFactor: s.Subject().CriticalStrikeFactor(),

		op: g,
	}
}

// MagicDamage returns a Damage
func (g *Game) MagicDamage(s Subject, o Object, d Statistic) *Damage {
	return &Damage{
		UnitPair:             MakePair(s, o),
		damage:               d * o.Object().MagicDamageReductionFactor(),
		criticalStrikeChance: s.Subject().CriticalStrikeChance(),
		criticalStrikeFactor: s.Subject().CriticalStrikeFactor(),

		op: g,
	}
}

// TrueDamage returns a Damage
func (g *Game) TrueDamage(s Subject, o Object, d Statistic) *Damage {
	return &Damage{
		UnitPair:             MakePair(s, o),
		damage:               d,
		criticalStrikeChance: s.Subject().CriticalStrikeChance(),
		criticalStrikeFactor: s.Subject().CriticalStrikeFactor(),

		op: g,
	}
}

// PureDamage returns a Damage
func (g *Game) PureDamage(s Subject, o Object, d Statistic) *Damage {
	return &Damage{
		UnitPair:             MakePair(s, o),
		damage:               d,
		criticalStrikeChance: 0,
		criticalStrikeFactor: 0,

		op: g,
	}
}

// Healing returns a Healing
func (g *Game) Healing(s Subject, o Object, h Statistic) *Healing {
	return &Healing{
		UnitPair:             MakePair(s, o),
		healing:              h,
		criticalStrikeChance: s.Subject().CriticalStrikeChance(),
		criticalStrikeFactor: s.Subject().CriticalStrikeFactor(),

		op: g,
	}
}
