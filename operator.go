package main

type operator interface {
	isComplete(u *unit) bool
	onAttach(u *unit)
	onTick(u *unit)
	onComplete(u *unit)
	onDetach(u *unit)
}

type activating struct {
	u *unit
	o chan message
}

// isComplete returns true iff ability is activated
func (a activating) isComplete(u *unit) bool {
	// todo implement
	return true
}

// onAttach confirms that the ability is available otherwise cancels activation of the ability
func (a activating) onAttach(u *unit) {
	a.u.attachStatsObserver(a)
	a.u.attachDisableObserver(a)
	if !a.checkCondition() {
		a.u.detachOperator(a)
		return
	}
	a.o <- message{
		// todo pack message
		t: outActivate,
	}
}

// onTick does nothing
func (a activating) onTick(u *unit) {}

// onComplete performs the ability
func (a activating) onComplete(u *unit) {
	// todo implement
}

// onDetach cleans up
func (a activating) onDetach(u *unit) {
	a.u.detachStatsObserver(a)
	a.u.detachDisableObserver(a)
}

// update confirms that the ability is available otherwise interrupts activation of the ability
func (a activating) update() {
	if a.checkCondition() {
		return
	}
	a.u.detachOperator(a)
	a.o <- message{
		// todo pack message
		t: outInterrupt,
	}
}

// checkConditions returns true iff the ability satisfies prior condition
func (a activating) checkCondition() bool {
	// todo check cooldown time
	// todo check health and mana
	// todo check disabler
	return true
}

type modifier struct {
	um unitModification
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

type disable struct {
	u *unit
	o chan message
}

// isComplete returns false iff the disable is effective
func (d disable) isComplete(u *unit) bool {
	// todo implement
	return false
}

// onAttach triggers notifyDisable and sends a message
func (d disable) onAttach(u *unit) {
	d.u.notifyDisable()
	// todo remove duplicate disable
	d.o <- message{
		// todo pack message
		t: outDisableBegin,
	}
}

// onTick does nothing
func (d disable) onTick(u *unit) {}

// onComplete sends a message
func (d disable) onComplete(u *unit) {
	d.o <- message{
		// todo pack message
		t: outDisableEnd,
	}
}

// onDetach does nothing
func (d disable) onDetach(u *unit) {}

type dot struct {
	u *unit
	o chan message
}

// isComplete returns false iff the DoT is effective
func (d dot) isComplete(u *unit) bool {
	// todo implement
	return false
}

// onAttach removes duplicate DoTs and sends a message
func (d dot) onAttach(u *unit) {
	// todo removes duplicate DoTs
	d.o <- message{
		// todo pack message
		t: outDoTBegin,
	}
}

// onTick performs DoT damage
func (d dot) onTick(u *unit) {
	// todo implement
}

// onComplete sends a message
func (d dot) onComplete(u *unit) {
	d.o <- message{
		// todo pack message
		t: outDoTEnd,
	}
}

// onDetach does nothing
func (d dot) onDetach(u *unit) {}

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
