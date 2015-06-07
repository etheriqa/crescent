package main

type modifier struct {
	partialOperator
	um *unitModification
}

// onAttach updates the modificationStats of the unit
func (m *modifier) onAttach() {
	m.addEventHandler(m, eventDead)
	m.addEventHandler(m, eventGameTick)
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
