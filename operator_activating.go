package main

type activating struct {
	unit *unit
}

// onAttach confirms that the ability is available otherwise cancels activation of the ability
func (a *activating) onAttach() {
	a.unit.addEventHandler(a, eventDisable)
	a.unit.addEventHandler(a, eventGameTick)
	a.unit.addEventHandler(a, eventStats)
	if !a.satisfiesRequirements() {
		a.unit.detachOperator(a)
		return
	}
	if !a.isActivated() {
		return
	}
	a.performAbility()
}

// onDetach cleans up
func (a *activating) onDetach() {
	a.unit.removeEventHandler(a, eventDisable)
	a.unit.removeEventHandler(a, eventGameTick)
	a.unit.removeEventHandler(a, eventStats)
}

// handleEvent checks the ability has been activated or not and performs it
func (a *activating) handleEvent(e event) {
	switch e {
	case eventDisable:
	case eventGameTick:
	case eventStats:
	default:
		return
	}
	if !a.satisfiesRequirements() {
		a.unit.detachOperator(a)
		a.unit.publish(message{
			// todo pack message
			t: outInterrupt,
		})
		return
	}
	if !a.isActivated() {
		return
	}
	a.performAbility()
	a.unit.detachOperator(a)
}

// satisfiesRequirements returns true iff the ability satisfies requirements
func (a *activating) satisfiesRequirements() bool {
	// todo check cooldown time
	// todo check health and mana
	// todo check disable
	return true
}

// isActivated returns true iff the ability is activated
func (a *activating) isActivated() bool {
	// todo check charging/casting time
	return true
}

// performAbility performs the ability
func (a *activating) performAbility() {
	// todo perform the ability
}
