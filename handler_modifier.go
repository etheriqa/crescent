package main

type Modifier struct {
	partialHandler
	unitModification
	name     string
	maxStack int
	nowStack int
}

// NewModifier initalizes a modifier
func NewModifier(receiver *unit, m unitModification, name string, maxStack int, duration gameDuration) *Modifier {
	return &Modifier{
		partialHandler: partialHandler{
			unit:           receiver,
			expirationTime: receiver.after(duration),
		},
		unitModification: m,
		name:             name,
		maxStack:         maxStack,
		nowStack:         1,
	}
}

// OnAttach updates the modificationStats of the unit
func (m *Modifier) OnAttach() {
	m.AddEventHandler(m, EventDead)
	m.AddEventHandler(m, EventGameTick)
	for ha := range m.handlers {
		switch ha := ha.(type) {
		case *Modifier:
			if ha == m || ha.name != m.name {
				continue
			}
			if ha.expirationTime > m.expirationTime {
				m.expirationTime = ha.expirationTime
			}
			m.nowStack += ha.nowStack
			if m.nowStack > m.maxStack {
				m.nowStack = m.maxStack
			}
			m.detachHandler(ha)
		}
	}
	m.updateModification()
	m.publish(message{
		// TODO pack message
		t: outModifierBegin,
	})
}

// OnDetach updates the modificationStats of the unit
func (m *Modifier) OnDetach() {
	m.RemoveEventHandler(m, EventDead)
	m.RemoveEventHandler(m, EventGameTick)
	m.updateModification()
	m.publish(message{
		// TODO pack message
		t: outModifierEnd,
	})
}

// HandleEvent handles the event
func (m *Modifier) HandleEvent(e Event) {
	switch e {
	case EventDead:
		m.detachHandler(m)
	case EventGameTick:
		if m.isExpired() {
			m.detachHandler(m)
		}
	}
}
