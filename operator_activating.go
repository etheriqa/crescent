package main

type activating struct {
	u *unit
	o chan message
}

// isComplete returns true iff ability is activated
func (a *activating) isComplete(u *unit) bool {
	// todo implement
	return true
}

// onAttach confirms that the ability is available otherwise cancels activation of the ability
func (a *activating) onAttach(u *unit) {
	u.addEventHandler(eventDisable, a)
	u.addEventHandler(eventStats, a)
	if !a.checkCondition() {
		u.detachOperator(a)
		return
	}
	a.o <- message{
		// todo pack message
		t: outActivate,
	}
}

// onTick does nothing
func (a *activating) onTick(u *unit) {}

// onComplete performs the ability
func (a *activating) onComplete(u *unit) {
	// todo implement
}

// onDetach cleans up
func (a *activating) onDetach(u *unit) {
	u.removeEventHandler(eventDisable, a)
	u.removeEventHandler(eventStats, a)
}

// handleEvent confirms that the ability is available otherwise interrupts activation of the ability
func (a *activating) handleEvent(e event) {
	switch e {
	case eventDisable:
	case eventStats:
	default:
		return
	}
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
func (a *activating) checkCondition() bool {
	// todo check cooldown time
	// todo check health and mana
	// todo check disabler
	return true
}
