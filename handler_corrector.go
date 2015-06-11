package main

type Corrector struct {
	PartialHandler
	correction UnitCorrection
	name       string
	maxStack   Statistic
	stack      Statistic
}

// NewCorrector returns a Corrector handler
func NewCorrector(object *Unit, c UnitCorrection, name string, maxStack Statistic, duration GameDuration) *Corrector {
	return &Corrector{
		PartialHandler: MakePartialHandler(MakeObject(object), duration),
		correction:     c,
		name:           name,
		maxStack:       maxStack,
		stack:          Statistic(1),
	}
}

// Stack returns the number of stacks
func (m *Corrector) Stack() Statistic {
	return m.stack
}

// Armor returns amount of Armor correction
func (m *Corrector) Armor() Statistic {
	return m.stack * m.correction.Armor
}

// MagicResistance returns amount of Magic Resistance correction
func (m *Corrector) MagicResistance() Statistic {
	return m.stack * m.correction.MagicResistance
}

// CriticalStrikeChance returns amount of Critical Strike Chance correction
func (m *Corrector) CriticalStrikeChance() Statistic {
	return m.stack * m.correction.CriticalStrikeChance
}

// CriticalStrikeFactor returns amount of Critical Strike Factor correction
func (m *Corrector) CriticalStrikeFactor() Statistic {
	return m.stack * m.correction.CriticalStrikeFactor
}

// CooldownReduction returns amount of Cooldown Reduction correction
func (m *Corrector) CooldownReduction() Statistic {
	return m.stack * m.correction.CooldownReduction
}

// DamageThreatFactor returns amount of Damage Threat Factor correction
func (m *Corrector) DamageThreatFactor() Statistic {
	return m.stack * m.correction.DamageThreatFactor
}

// HealingThreatFactor returns amount of Healing Threat Factor correction
func (m *Corrector) HealingThreatFactor() Statistic {
	return m.stack * m.correction.HealingThreatFactor
}

// OnAttach updates the modificationStats of the unit
func (m *Corrector) OnAttach() {
	m.Object().AddEventHandler(m, EventDead)
	m.Object().AddEventHandler(m, EventGameTick)
	m.ForObjectHandler(m.Object(), func(ha Handler) {
		switch ha := ha.(type) {
		case *Corrector:
			if ha == m || ha.name != m.name {
				return
			}
			if ha.expirationTime > m.expirationTime {
				m.expirationTime = ha.expirationTime
			}
			m.stack += ha.stack
			if m.stack > m.maxStack {
				m.stack = m.maxStack
			}
			ha.Stop(ha)
		}
	})
	m.Object().ReloadCorrection()
}

// OnDetach updates the modificationStats of the unit
func (m *Corrector) OnDetach() {
	m.Object().RemoveEventHandler(m, EventDead)
	m.Object().RemoveEventHandler(m, EventGameTick)
	m.Object().ReloadCorrection()
}

// HandleEvent handles the event
func (m *Corrector) HandleEvent(e Event) {
	switch e {
	case EventDead:
		m.Stop(m)
	case EventGameTick:
		if m.IsExpired() {
			m.Up()
		}
	}
}

// Up ends the Corrector
func (m *Corrector) Up() {
	m.Stop(m)
	m.Publish(message{
	// TODO pack message
	})
}
