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
	expirationTime GameTime

	op Operator
}

// ArmorCorrection returns amount of armor correction
func (h *Correction) ArmorCorrection() Statistic {
	return h.correction.Armor
}

// MagicResistanceCorrection returns amount of magic resistance correction
func (h *Correction) MagicResistanceCorrection() Statistic {
	return h.correction.MagicResistance
}

// CriticalStrikeChanceCorrection returns amount of critical strike chance correction
func (h *Correction) CriticalStrikeChanceCorrection() Statistic {
	return h.correction.CriticalStrikeChance
}

// CriticalStrikeFactorCorrection returns amount of critical strike factor correction
func (h *Correction) CriticalStrikeFactorCorrection() Statistic {
	return h.correction.CriticalStrikeFactor
}

// CooldownReductionCorrection returns amount of cooldown reduction correction
func (h *Correction) CooldownReductionCorrection() Statistic {
	return h.correction.CooldownReduction
}

// DamageThreatFactorCorrection returns amount of damage threat factor correction
func (h *Correction) DamageThreatFactorCorrection() Statistic {
	return h.correction.DamageThreatFactor
}

// HealingThreatFactorCorrection returns amount of healing threat factor correction
func (h *Correction) HealingThreatFactorCorrection() Statistic {
	return h.correction.HealingThreatFactor
}

// Stack returns stack number
func (h *Correction) Stack() Statistic {
	return h.stack
}

// OnAttach merges Correction handlers and updates the UnitCorrection of the Object
func (h *Correction) OnAttach() {
	h.op.Handlers().BindObject(h).Each(func(o Handler) {
		switch o := o.(type) {
		case *Correction:
			if h == o || h.name != o.name {
				return
			}
			if h.expirationTime < o.expirationTime {
				h.expirationTime = o.expirationTime
				h.stack = o.stack
			} else {
				h.stack += o.stack
			}
			h.op.Handlers().Detach(o)
		}
	})

	if h.stack > h.stackLimit {
		h.stack = h.stackLimit
	}
	h.Object().AddEventHandler(h, EventGameTick)
	h.updateCorrection()
	h.op.Writer().Write(nil) // TODO
}

// OnDetach updates the UnitCorrection of the Object
func (h *Correction) OnDetach() {
	h.Object().RemoveEventHandler(h, EventGameTick)
	h.updateCorrection()
}

// HandleEvent handles the Event
func (h *Correction) HandleEvent(e Event) {
	switch e {
	case EventGameTick:
		if h.op.Clock().Before(h.expirationTime) {
			return
		}
		h.op.Handlers().Detach(h)
		h.op.Writer().Write(nil) // TODO
	}
}

// updateCorrection updates the UnitCorrection of the Object
func (h *Correction) updateCorrection() {
	c := UnitCorrection{}
	h.op.Handlers().BindObject(h).Each(func(o Handler) {
		switch o := o.(type) {
		case Corrector:
			c.Armor += o.ArmorCorrection()
			c.MagicResistance += o.MagicResistanceCorrection()
			c.CriticalStrikeChance += o.CriticalStrikeChanceCorrection()
			c.CriticalStrikeFactor += o.CriticalStrikeFactorCorrection()
			c.CooldownReduction += o.CooldownReductionCorrection()
			c.DamageThreatFactor += o.DamageThreatFactorCorrection()
			c.HealingThreatFactor += o.HealingThreatFactorCorrection()
		}
	})
	h.Object().UpdateCorrection(c)
}
