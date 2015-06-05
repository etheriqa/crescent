package main

type disable struct {
	unit *unit
}

// onAttach triggers eventDisable and sends a message
func (d *disable) onAttach() {
	d.unit.addEventHandler(d, eventGameTick)
	// todo remove duplicate disable
	d.unit.publish(message{
		// todo pack message
		t: outDisableBegin,
	})
}

// onDetach removes eventHandler
func (d *disable) onDetach() {
	d.unit.removeEventHandler(d, eventGameTick)
}

// handleEvent checks the disable has been expired or not
func (d *disable) handleEvent(e event) {
	switch e {
	case eventGameTick:
		d.expire()
	default:
		return
	}
}

// expire expires the disable iff it is expired
func (d *disable) expire() {
	// todo expire
}
