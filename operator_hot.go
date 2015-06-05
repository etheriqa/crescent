package main

type hot struct {
	u *unit
	o chan message
}

// onAttach removes duplicate HoTs and sends a message
func (h *hot) onAttach(u *unit) {
	u.addEventHandler(h, eventGameTick)
	u.addEventHandler(h, eventStatsTick)
	// todo removes duplicate HoTs
	h.o <- message{
		// todo pack message
		t: outHoTBegin,
	}
}

// onDetach removes the eventHandlers
func (h *hot) onDetach(u *unit) {
	u.removeEventHandler(h, eventGameTick)
	u.removeEventHandler(h, eventStatsTick)
}

// handleEvent handles events
func (h *hot) handleEvent(e event) {
	switch e {
	case eventGameTick:
		h.expire()
	case eventStatsTick:
		h.perform()
	}
}

// expire expires the HoT iff it is expired
func (h *hot) expire() {
	// todo expire
	h.o <- message{
		// todo pack message
		t: outHoTEnd,
	}
}

// perform performs the HoT
func (h *hot) perform() {
	// todo perform the HoT
}
