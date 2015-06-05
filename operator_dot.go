package main

type dot struct {
	u *unit
	o chan message
}

// onAttach removes duplicate DoTs and sends a message
func (d *dot) onAttach(u *unit) {
	u.addEventHandler(d, eventGameTick)
	u.addEventHandler(d, eventStatsTick)
	// todo removes duplicate DoTs
	d.o <- message{
		// todo pack message
		t: outDoTBegin,
	}
}

// onDetach removes the eventHandlers
func (d *dot) onDetach(u *unit) {
	u.removeEventHandler(d, eventGameTick)
	u.removeEventHandler(d, eventStatsTick)
}

// handleEvent handles events
func (d *dot) handleEvent(e event) {
	switch e {
	case eventGameTick:
		d.expire()
	case eventStatsTick:
		d.perform()
	}
}

// expire expires the DoT iff it is expired
func (d *dot) expire() {
	// todo expire
}

// perform performs the DoT
func (d *dot) perform() {
	// todo perform the DoT
	d.o <- message{
		// todo pack message
		t: outDoTEnd,
	}
}
