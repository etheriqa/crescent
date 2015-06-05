package main

type dot struct {
	unit *unit
}

// onAttach removes duplicate DoTs and sends a message
func (d *dot) onAttach() {
	d.unit.addEventHandler(d, eventGameTick)
	d.unit.addEventHandler(d, eventStatsTick)
	// todo removes duplicate DoTs
	d.unit.publish(message{
		// todo pack message
		t: outDoTBegin,
	})
}

// onDetach removes the eventHandlers
func (d *dot) onDetach() {
	d.unit.removeEventHandler(d, eventGameTick)
	d.unit.removeEventHandler(d, eventStatsTick)
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
	d.unit.publish(message{
		// todo pack message
		t: outDoTEnd,
	})
}
