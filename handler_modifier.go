package main

type Modifier struct {
	*PartialHandler
	unitModification
	name     string
	maxStack Statistic
	stack    Statistic
}

// NewModifier returns a Modifier handler
func NewModifier(object *Unit, m unitModification, name string, maxStack Statistic, duration GameDuration) *Modifier {
	return &Modifier{
		PartialHandler:   NewPartialHandler(nil, object, duration),
		unitModification: m,
		name:             name,
		maxStack:         maxStack,
		stack:            Statistic(1),
	}
}

// Stack returns the number of stacks
func (m *Modifier) Stack() Statistic {
	return m.stack
}

// OnAttach updates the modificationStats of the unit
func (m *Modifier) OnAttach() {
	m.Object().AddEventHandler(m, EventDead)
	m.Object().AddEventHandler(m, EventGameTick)
	m.Container().ForObjectHandler(m.Object(), func(ha Handler) {
		switch ha := ha.(type) {
		case *Modifier:
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
	m.Object().updateModification()
}

// OnDetach updates the modificationStats of the unit
func (m *Modifier) OnDetach() {
	m.Object().RemoveEventHandler(m, EventDead)
	m.Object().RemoveEventHandler(m, EventGameTick)
	m.Object().updateModification()
}

// HandleEvent handles the event
func (m *Modifier) HandleEvent(e Event) {
	switch e {
	case EventDead:
		m.Stop(m)
	case EventGameTick:
		if m.IsExpired() {
			m.Up()
		}
	}
}

// Up ends the Modifier
func (m *Modifier) Up() {
	m.Stop(m)
	m.Publish(message{
	// TODO pack message
	})
}
