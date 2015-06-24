package crescent

import (
	"math/rand"
)

type Game interface {
	Rand() *rand.Rand
	Clock() InstanceClock
	Writer() InstanceOutputWriter

	Join(UnitGroup, UnitName, *Class) (UnitID, error)
	Leave(UnitID) error
	UnitQuery() UnitQueryable

	AttachEffect(Effect) error
	DetachEffect(Effect) error
	EffectQuery() EffectQueryable

	Activating(Subject, *Unit, *Ability)
	Cooldown(Object, *Ability)
	ResetCooldown(Object, *Ability)
	Correction(Object, UnitCorrection, string, Statistic, InstanceDuration)
	Disable(Object, DisableType, InstanceDuration)
	DamageThreat(Subject, Object, Statistic)
	HealingThreat(Subject, Object, Statistic)
	DoT(*Damage, string, InstanceDuration)
	HoT(*Healing, string, InstanceDuration)

	PhysicalDamage(Subject, Object, Statistic) *Damage
	MagicDamage(Subject, Object, Statistic) *Damage
	TrueDamage(Subject, Object, Statistic) *Damage
	PureDamage(Subject, Object, Statistic) *Damage

	Healing(Subject, Object, Statistic) *Healing
}

// Rand returns the *rand.Rand
func (g *GameState) Rand() *rand.Rand {
	return g.rand
}

// Clock returns the InstanceClock
func (g *GameState) Clock() InstanceClock {
	return g.clock
}

// Writer returns the InstanceOutputWriter
func (g *GameState) Writer() InstanceOutputWriter {
	return g.w
}

// Join creates a Unit and adds it to the game
func (g *GameState) Join(group UnitGroup, name UnitName, class *Class) (id UnitID, err error) {
	u, err := g.units.Join(group, name, class)
	if err != nil {
		return
	}
	id = u.ID()
	g.w.Write(OutputUnitJoin{
		UnitID:    u.ID(),
		UnitGroup: u.Group(),
		UnitName:  u.Name(),
		ClassName: u.ClassName(),
		Health:    u.Health(),
		HealthMax: u.HealthMax(),
		Mana:      u.Mana(),
		ManaMax:   u.ManaMax(),
	})
	return
}

// Leave removes the Unit
func (g *GameState) Leave(id UnitID) (err error) {
	if err = g.units.Leave(id); err != nil {
		return
	}
	g.w.Write(OutputUnitLeave{
		UnitID: id,
	})
	return
}

// UnitQuery returns a UnitQueryable
func (g *GameState) UnitQuery() UnitQueryable {
	return g.units
}

// AddEffect adds the effect
func (g *GameState) AttachEffect(e Effect) error {
	return g.effects.Attach(g, e)
}

// RemoveEffect removes the effect
func (g *GameState) DetachEffect(e Effect) error {
	return g.effects.Detach(g, e)
}

// EffectQuery returns a EffectQueryable
func (g *GameState) EffectQuery() EffectQueryable {
	return g.effects
}

// Activating attaches a Activating Effect
func (g *GameState) Activating(s Subject, o *Unit, a *Ability) {
	g.AttachEffect(NewActivating(g, s, o, a, g.clock.Add(a.ActivationDuration)))
}

// Cooldown attaches a Cooldown Effect
func (g *GameState) Cooldown(o Object, a *Ability) {
	g.AttachEffect(NewCooldown(g, o, a, g.clock.Add(a.CooldownDuration)))
}

// ResetCooldown detaches Cooldown effects
func (g *GameState) ResetCooldown(o Object, a *Ability) {
	g.AttachEffect(NewCooldown(g, o, a, g.clock.Now()))
}

// Correction attaches a Correction Effect
func (g *GameState) Correction(o Object, c UnitCorrection, name string, l Statistic, d InstanceDuration) {
	g.AttachEffect(NewCorrection(g, o, c, name, l, g.clock.Add(d)))
}

// Disable attaches a Disable Effect
func (g *GameState) Disable(o Object, dt DisableType, d InstanceDuration) {
	g.AttachEffect(NewDisable(g, o, dt, g.clock.Add(d)))
}

// DamageThreat attaches a Threat Effect
func (g *GameState) DamageThreat(s Subject, o Object, d Statistic) {
	g.AttachEffect(NewThreat(g, s, o, d*s.Subject().DamageThreatFactor()))
}

// HealingThreat attaches a Threat Effect
func (g *GameState) HealingThreat(s Subject, o Object, h Statistic) {
	g.AttachEffect(NewThreat(g, s, o, h*s.Subject().HealingThreatFactor()))
}

// DoT attaches a Periodical Effect
func (g *GameState) DoT(damage *Damage, name string, d InstanceDuration) {
	g.AttachEffect(NewPeriodical(g, damage, damage, name, func() { damage.Perform() }, g.clock.Add(d)))
}

// HoT attaches a Periodical Effect
func (g *GameState) HoT(healing *Healing, name string, d InstanceDuration) {
	g.AttachEffect(NewPeriodical(g, healing, healing, name, func() { healing.Perform() }, g.clock.Add(d)))
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
