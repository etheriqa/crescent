package main

type activating struct {
	partialOperator
	*ability
	receiver *unit
}

// onAttach checks requirements
func (a *activating) onAttach() {
	a.addEventHandler(a, eventDead)
	a.addEventHandler(a, eventDisable)
	a.addEventHandler(a, eventGameTick)
	a.addEventHandler(a, eventStats)
	a.checkRequirements()
}

// onDetach removes the eventHandlers
func (a *activating) onDetach() {
	a.removeEventHandler(a, eventDead)
	a.removeEventHandler(a, eventDisable)
	a.removeEventHandler(a, eventGameTick)
	a.removeEventHandler(a, eventStats)
}

// handleEvent handles the event
func (a *activating) handleEvent(e event) {
	switch e {
	case eventDead:
		a.detachOperator(a)
	case eventDisable:
		a.checkRequirements()
	case eventGameTick:
		a.perform()
	case eventStats:
		a.checkRequirements()
	}
}

func (a *activating) checkRequirements() bool {
	if a.satisfiedRequirements(a.unit) {
		return true
	}
	a.detachOperator(a)
	a.publish(message{
		t: outInterrupt,
		// TODO pack message
	})
	return false
}

func (a *activating) perform() {
	if !a.checkRequirements() {
		return
	}
	if !a.isExpired() {
		return
	}
	a.ability.perform(a.unit, a.receiver)
	a.detachOperator(a)
}
