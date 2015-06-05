package main

type modifier struct {
	unit *unit
	um   *unitModification
}

// onAttach updates the modification of the unit
func (m *modifier) onAttach() {
	m.unit.addEventHandler(m, eventGameTick)
	m.unit.updateModification()
	m.unit.publish(message{
		// todo pack message
		t: outModifierBegin,
	})
}

// onDetach updates the modification of the unit
func (m *modifier) onDetach() {
	m.unit.removeEventHandler(m, eventGameTick)
	m.unit.updateModification()
}

// handleEvent han
func (m *modifier) handleEvent(e event) {
	switch e {
	case eventGameTick:
		m.expire()
	}
}

// expire expires the modifier iff it is expired
func (m *modifier) expire() {
	// todo expire
	m.unit.publish(message{
		// todo pack message
		t: outModifierEnd,
	})
}
