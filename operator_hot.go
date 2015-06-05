package main

type hot struct {
	unit *unit
}

// onAttach removes duplicate HoTs and sends a message
func (h *hot) onAttach() {
	h.unit.addEventHandler(h, eventGameTick)
	h.unit.addEventHandler(h, eventStatsTick)
	// todo removes duplicate HoTs
	h.unit.publish(message{
		// todo pack message
		t: outHoTBegin,
	})
}

// onDetach removes the eventHandlers
func (h *hot) onDetach() {
	h.unit.removeEventHandler(h, eventGameTick)
	h.unit.removeEventHandler(h, eventStatsTick)
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
	h.unit.publish(message{
		// todo pack message
		t: outHoTEnd,
	})
}

// perform performs the HoT
func (h *hot) perform() {
	// todo perform the HoT
}
