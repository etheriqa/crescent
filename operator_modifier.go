package main

type modifier struct {
	um *unitModification
	o  chan message
}

// onAttach updates the modification of the unit
func (m *modifier) onAttach(u *unit) {
	u.addEventHandler(m, eventGameTick)
	u.updateModification()
	m.o <- message{
		// todo pack message
		t: outModifierBegin,
	}
}

// onDetach updates the modification of the unit
func (m *modifier) onDetach(u *unit) {
	u.removeEventHandler(m, eventGameTick)
	u.updateModification()
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
	m.o <- message{
		// todo pack message
		t: outModifierEnd,
	}
}
