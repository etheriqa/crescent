package main

type dot struct {
	u *unit
	o chan message
}

// isComplete returns false iff the DoT is effective
func (d *dot) isComplete(u *unit) bool {
	// todo implement
	return false
}

// onAttach removes duplicate DoTs and sends a message
func (d *dot) onAttach(u *unit) {
	// todo removes duplicate DoTs
	d.o <- message{
		// todo pack message
		t: outDoTBegin,
	}
}

// onTick performs DoT damage
func (d *dot) onTick(u *unit) {
	// todo implement
}

// onComplete sends a message
func (d *dot) onComplete(u *unit) {
	d.o <- message{
		// todo pack message
		t: outDoTEnd,
	}
}

// onDetach does nothing
func (d *dot) onDetach(u *unit) {}
