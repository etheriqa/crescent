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
func (a *activating) onTick(u *unit) {}

// onComplete performs the ability
func (a *activating) onComplete(u *unit) {
	// todo implement
}

// onDetach cleans up
func (a *activating) onDetach(u *unit) {
	a.u.detachStatsObserver(a)
	a.u.detachDisableObserver(a)
}

// update confirms that the ability is available otherwise interrupts activation of the ability
func (a *activating) update() {
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