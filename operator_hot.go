package main

type hot struct {
	u *unit
	o chan message
}

// isComplete returns false iff the HoT is effective
func (h hot) isComplete(u *unit) bool {
	// todo implement
	return false
}

// onAttach removes duplicate HoTs and sends a message
func (h hot) onAttach(u *unit) {
	// todo removes duplicate HoTs
	h.o <- message{
		// todo pack message
		t: outHoTBegin,
	}
}

// onTick performs HoT healing
func (h hot) onTick(u *unit) {
	// todo implement
}

// onComplete sends a message
func (h hot) onComplete(u *unit) {
	h.o <- message{
		// todo pack message
		t: outHoTEnd,
	}
}

// onDetach does nothing
func (h hot) onDetach(u *unit) {}
