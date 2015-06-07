package main

type modifier struct {
	partialOperator
	um *unitModification
}

// onAttach updates the modificationStats of the unit
func (m *modifier) onAttach() {
	m.unit.addEventHandler(m, eventDead)
	m.unit.addEventHandler(m, eventGameTick)
	m.unit.updateModification()
	m.unit.publish(message{
		// TODO pack message
		t: outModifierBegin,
	})
}

// onDetach updates the modificationStats of the unit
func (m *modifier) onDetach() {
	m.unit.removeEventHandler(m, eventDead)
	m.unit.removeEventHandler(m, eventGameTick)
	m.unit.updateModification()
}

// handleEvent handles the event
func (m *modifier) handleEvent(e event) {
	switch e {
	case eventDead:
		m.terminate(m)
	case eventGameTick:
		m.expire(m, message{
			// TODO pack message
			t: outModifierEnd,
		})
	}
}
