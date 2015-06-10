package main

type modifier struct {
	partialOperator
	unitModification
	*ability
	maxStack int
	nowStack int
}

// newModifier initalizes a modifier
func newModifier(receiver *unit, duration gameDuration, m unitModification, a *ability, maxStack int) *modifier {
	return &modifier{
		partialOperator: partialOperator{
			unit:           receiver,
			expirationTime: receiver.after(duration),
		},
		unitModification: m,
		ability:          a,
		maxStack:         maxStack,
		nowStack:         1,
	}
}

// onAttach updates the modificationStats of the unit
func (m *modifier) onAttach() {
	m.addEventHandler(m, eventDead)
	m.addEventHandler(m, eventGameTick)
	for o := range m.operators {
		switch o := o.(type) {
		case *modifier:
			if o == m {
				continue
			}
			if o.ability != m.ability {
				continue
			}
			// TODO check the expiration time
			m.detachOperator(o)
			if m.nowStack < m.maxStack {
				m.nowStack++
			}
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
	m.removeEventHandler(m, eventDead)
	m.removeEventHandler(m, eventGameTick)
	m.updateModification()
}

// handleEvent handles the event
func (m *modifier) handleEvent(e event) {
	switch e {
	case eventDead:
		m.detachOperator(m)
	case eventGameTick:
		m.expire(m, message{
			// TODO pack message
			t: outModifierEnd,
		})
	}
}
