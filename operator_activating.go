package main

type activating struct {
	partialOperator
	ability ability
}

// onAttach checks requirements
func (a *activating) onAttach() {
	a.unit.addEventHandler(a, eventDead)
	a.unit.addEventHandler(a, eventDisable)
	a.unit.addEventHandler(a, eventGameTick)
	a.unit.addEventHandler(a, eventStats)
	a.checkRequirements()
}

// onDetach removes the eventHandlers
func (a *activating) onDetach() {
	a.unit.removeEventHandler(a, eventDead)
	a.unit.removeEventHandler(a, eventDisable)
	a.unit.removeEventHandler(a, eventGameTick)
	a.unit.removeEventHandler(a, eventStats)
}

// handleEvent handles the event
func (a *activating) handleEvent(e event) {
	switch e {
	case eventDead:
		a.terminate(a)
	case eventDisable:
		a.checkRequirements()
	case eventGameTick:
		a.perform()
	case eventStats:
		a.checkRequirements()
	}
}

func (a *activating) checkRequirements() bool {
	if a.ability.satisfiedRequirements(a.unit) {
		return true
	}
	a.terminate(a)
	a.unit.publish(message{
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
	a.ability.perform()
	a.terminate(a)
}
