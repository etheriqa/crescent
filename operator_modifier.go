package main

type modifier struct {
	um *unitModification
	o  chan message
}

// isComplete returns false iff the modifier is effective
func (m modifier) isComplete(u *unit) bool {
	// todo implement
	return false
}

// onAttach updates the modification of the unit
func (m modifier) onAttach(u *unit) {
	u.updateModification()
	m.o <- message{
		// todo pack message
		t: outModifierBegin,
	}
}

// onTick does nothing
func (m modifier) onTick(u *unit) {}

// onComplete sends a message
func (m modifier) onComplete(u *unit) {
	m.o <- message{
		// todo pack message
		t: outModifierEnd,
	}
}

// onDetach updates the modification of the unit
func (m modifier) onDetach(u *unit) {
	u.updateModification()
}
