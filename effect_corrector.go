package main

type Corrector interface {
	ArmorCorrection() Statistic
	MagicResistanceCorrection() Statistic
	CriticalStrikeChanceCorrection() Statistic
	CriticalStrikeFactorCorrection() Statistic
	CooldownReductionCorrection() Statistic
	DamageThreatFactorCorrection() Statistic
	HealingThreatFactorCorrection() Statistic
}

type Correction struct {
	UnitObject
	name           string
	correction     UnitCorrection
	stackLimit     Statistic
	stack          Statistic
	expirationTime InstanceTime

	handler EventHandler
}

// NewCorrection returns a Correction
func NewCorrection(g Game, o Object, c UnitCorrection, name string, l Statistic, t InstanceTime) *Correction {
	e := &Correction{
		UnitObject:     MakeObject(o),
		name:           name,
		correction:     c,
		stackLimit:     l,
		stack:          1,
		expirationTime: t,
		handler:        new(func(interface{})),
	}
	*e.handler = func(p interface{}) { e.handle(g, p) }
	return e
}

// ArmorCorrection returns amount of armor correction
func (e *Correction) ArmorCorrection() Statistic {
	return e.correction.Armor
}

// MagicResistanceCorrection returns amount of magic resistance correction
func (e *Correction) MagicResistanceCorrection() Statistic {
	return e.correction.MagicResistance
}

// CriticalStrikeChanceCorrection returns amount of critical strike chance correction
func (e *Correction) CriticalStrikeChanceCorrection() Statistic {
	return e.correction.CriticalStrikeChance
}

// CriticalStrikeFactorCorrection returns amount of critical strike factor correction
func (e *Correction) CriticalStrikeFactorCorrection() Statistic {
	return e.correction.CriticalStrikeFactor
}

// CooldownReductionCorrection returns amount of cooldown reduction correction
func (e *Correction) CooldownReductionCorrection() Statistic {
	return e.correction.CooldownReduction
}

// DamageThreatFactorCorrection returns amount of damage threat factor correction
func (e *Correction) DamageThreatFactorCorrection() Statistic {
	return e.correction.DamageThreatFactor
}

// HealingThreatFactorCorrection returns amount of healing threat factor correction
func (e *Correction) HealingThreatFactorCorrection() Statistic {
	return e.correction.HealingThreatFactor
}

// Stack returns stack number
func (e *Correction) Stack() Statistic {
	return e.stack
}

// EffectWillAttach merges Correction effects
func (e *Correction) EffectWillAttach(g Game) error {
	g.EffectQuery().BindObject(e).Each(func(f Effect) {
		switch f := f.(type) {
		case *Correction:
			if e.name != f.name {
				return
			}
			if e.expirationTime < f.expirationTime {
				e.expirationTime = f.expirationTime
				e.stack = f.stack
			} else {
				e.stack += f.stack
			}
			g.DetachEffect(f)
		}
	})

	if e.stack > e.stackLimit {
		e.stack = e.stackLimit
	}
	return nil
}

// EffectDidAttach updates the UnitCorrection of the object unit
func (e *Correction) EffectDidAttach(g Game) error {
	e.writeOutputUnitAttach(g)
	e.Object().Register(e.handler)
	e.updateCorrection(g)
	return nil
}

// EffectDidDetach updates the UnitCorrection of the object unit
func (e *Correction) EffectDidDetach(g Game) error {
	e.Object().Unregister(e.handler)
	e.updateCorrection(g)
	return nil
}

// handle handles the payload
func (e *Correction) handle(g Game, p interface{}) {
	switch p.(type) {
	case EventGameTick:
		if g.Clock().Before(e.expirationTime) {
			return
		}
		e.writeOutputUnitDetach(g)
		g.DetachEffect(e)
	}
}

// updateCorrection updates the UnitCorrection of the Object
func (e *Correction) updateCorrection(g Game) {
	c := MakeUnitCorrection()
	g.EffectQuery().BindObject(e).Each(func(f Effect) {
		switch f := f.(type) {
		case Corrector:
			c.Armor += f.ArmorCorrection()
			c.MagicResistance += f.MagicResistanceCorrection()
			c.CriticalStrikeChance += f.CriticalStrikeChanceCorrection()
			c.CriticalStrikeFactor += f.CriticalStrikeFactorCorrection()
			c.CooldownReduction += f.CooldownReductionCorrection()
			c.DamageThreatFactor += f.DamageThreatFactorCorrection()
			c.HealingThreatFactor += f.HealingThreatFactorCorrection()
		}
	})
	e.Object().UpdateCorrection(c)
}

// writeOutputUnitAttach writes a OutputUnitAttach
func (e *Correction) writeOutputUnitAttach(g Game) {
	g.Writer().Write(OutputUnitAttach{
		UnitID:         e.Object().ID(),
		AttachmentName: e.name,
		Stack:          e.stack,
		ExpirationTime: e.expirationTime,
	})
}

// writeOutputUnitDetach writes a OutputUnitDetach
func (e *Correction) writeOutputUnitDetach(g Game) {
	g.Writer().Write(OutputUnitDetach{
		UnitID:         e.Object().ID(),
		AttachmentName: e.name,
	})
}
