package main

type disable struct {
	u *unit
	o chan message
}

// isComplete returns false iff the disable is effective
func (d *disable) isComplete(u *unit) bool {
	// todo implement
	return false
}

// onAttach triggers notifyDisable and sends a message
func (d *disable) onAttach(u *unit) {
	d.u.notifyDisable()
	// todo remove duplicate disable
	d.o <- message{
		// todo pack message
		t: outDisableBegin,
	}
}

// onTick does nothing
func (d *disable) onTick(u *unit) {}

// onComplete sends a message
func (d *disable) onComplete(u *unit) {
	d.o <- message{
		// todo pack message
		t: outDisableEnd,
	}
}

// onDetach does nothing
func (d *disable) onDetach(u *unit) {}
