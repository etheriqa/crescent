package main

type disable struct {
	u *unit
	o chan message
}

// onAttach triggers eventDisable and sends a message
func (d *disable) onAttach(u *unit) {
	u.addEventHandler(d, eventGameTick)
	// todo remove duplicate disable
	d.o <- message{
		// todo pack message
		t: outDisableBegin,
	}
}

// onDetach removes eventHandler
func (d *disable) onDetach(u *unit) {
	u.removeEventHandler(d, eventGameTick)
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
