package main

type activating struct {
	partialOperator
}

// onAttach checks requirements
func (a *activating) onAttach() {
	a.unit.addEventHandler(a, eventDisable)
	a.unit.addEventHandler(a, eventGameTick)
	a.unit.addEventHandler(a, eventStats)
	if !a.satisfiesRequirements() {
		a.unit.detachOperator(a)
		return
	}
	if !a.isExpired() {
		return
	}
	a.performAbility()
}

// onDetach removes the eventHandlers
func (a *activating) onDetach() {
	a.unit.removeEventHandler(a, eventDisable)
	a.unit.removeEventHandler(a, eventGameTick)
	a.unit.removeEventHandler(a, eventStats)
}

// handleEvent handles the event
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
	if !a.isExpired() {
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

// performAbility performs the ability
func (a *activating) performAbility() {
	// todo perform the ability
}
