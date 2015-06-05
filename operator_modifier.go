package main

type modifier struct {
	partialOperator
	um *unitModification
}

// onAttach updates the modificationStats of the unit
func (m *modifier) onAttach() {
	m.unit.addEventHandler(m, eventGameTick)
	m.unit.updateModification()
	m.unit.publish(message{
		// todo pack message
		t: outModifierBegin,
	})
}

// onDetach updates the modificationStats of the unit
func (m *modifier) onDetach() {
	m.unit.removeEventHandler(m, eventGameTick)
	m.unit.updateModification()
}

// handleEvent handles the event
func (m *modifier) handleEvent(e event) {
	switch e {
	case eventGameTick:
		m.expire(m, message{
			// todo pack message
			t: outModifierEnd,
		})
	}
}
