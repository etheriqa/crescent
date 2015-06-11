package main

type modifier struct {
	partialOperator
	unitModification
	name     string
	maxStack int
	nowStack int
}

// newModifier initalizes a modifier
func newModifier(receiver *unit, m unitModification, name string, maxStack int, duration gameDuration) *modifier {
	return &modifier{
		partialOperator: partialOperator{
			unit:           receiver,
			expirationTime: receiver.after(duration),
		},
		unitModification: m,
		name:             name,
		maxStack:         maxStack,
		nowStack:         1,
	}
}

// onAttach updates the modificationStats of the unit
func (m *modifier) onAttach() {
	m.AddEventHandler(m, EventDead)
	m.AddEventHandler(m, EventGameTick)
	for o := range m.operators {
		switch o := o.(type) {
		case *modifier:
			if o == m || o.name != m.name {
				continue
			}
			if o.expirationTime > m.expirationTime {
				m.expirationTime = o.expirationTime
			}
			m.nowStack += o.nowStack
			if m.nowStack > m.maxStack {
				m.nowStack = m.maxStack
			}
			m.detachOperator(o)
		}
	}
	m.updateModification()
	m.publish(message{
		// TODO pack message
		t: outModifierBegin,
	})
}

// onDetach updates the modificationStats of the unit
func (m *modifier) onDetach() {
	m.RemoveEventHandler(m, EventDead)
	m.RemoveEventHandler(m, EventGameTick)
	m.updateModification()
	m.publish(message{
		// TODO pack message
		t: outModifierEnd,
	})
}

// HandleEvent handles the event
func (m *modifier) HandleEvent(e Event) {
	switch e {
	case EventDead:
		m.detachOperator(m)
	case EventGameTick:
		if m.isExpired() {
			m.detachOperator(m)
		}
	}
}
