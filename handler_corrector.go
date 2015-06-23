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
	h.writeOutputUnitAttach()
	h.Object().Register(h)
	h.updateCorrection()
}

// OnDetach updates the UnitCorrection of the Object
func (h *Correction) OnDetach() {
	h.Object().Unregister(h)
	h.updateCorrection()
}

// Handle handles the Event
func (h *Correction) Handle(p interface{}) {
	switch p.(type) {
	case *EventGameTick:
		if h.op.Clock().Before(h.expirationTime) {
			return
		}
		h.writeOutputUnitDetach()
		h.op.Handlers().Detach(h)
	}
}

// updateCorrection updates the UnitCorrection of the Object
func (h *Correction) updateCorrection() {
	c := MakeUnitCorrection()
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

// writeOutputUnitAttach writes a OutputUnitAttach
func (h *Correction) writeOutputUnitAttach() {
	h.op.Writer().Write(OutputUnitAttach{
		UnitID:         h.Object().ID(),
		AttachmentName: h.name,
		Stack:          h.stack,
		ExpirationTime: h.expirationTime,
	})
}

// writeOutputUnitDetach writes a OutputUnitDetach
func (h *Correction) writeOutputUnitDetach() {
	h.op.Writer().Write(OutputUnitDetach{
		UnitID:         h.Object().ID(),
		AttachmentName: h.name,
	})
}
