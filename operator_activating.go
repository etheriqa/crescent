package main

type activating struct {
	partialOperator
	*ability
	receiver *unit
}

// onAttach checks requirements
func (a *activating) onAttach() {
	a.addEventHandler(a, eventDead)
	a.addEventHandler(a, eventDisableInterrupt)
	a.addEventHandler(a, eventGameTick)
	a.addEventHandler(a, eventResourceDecreased)
	for o := range a.operators {
		switch o.(type) {
		case *activating:
			a.detachOperator(a)
			return
		}
	}
	if err := a.checkRequirements(a.unit, a.receiver); err != nil {
		a.detachOperator(a)
	}
	a.publish(message{
	// TODO pack message
	})
}

// onDetach removes the eventHandlers
func (a *activating) onDetach() {
	a.removeEventHandler(a, eventDead)
	a.removeEventHandler(a, eventDisableInterrupt)
	a.removeEventHandler(a, eventGameTick)
	a.removeEventHandler(a, eventResourceDecreased)
}

// handleEvent handles the event
func (a *activating) handleEvent(e event) {
	switch e {
	case eventDead:
		a.detachOperator(a)
	case eventDisableInterrupt:
		a.perform()
	case eventGameTick:
		a.perform()
	case eventResourceDecreased:
		a.perform()
	}
}

// perform performs the ability
func (a *activating) perform() {
	if err := a.checkRequirements(a.unit, a.receiver); err != nil {
		a.publish(message{
		// TODO pack message
		})
		a.detachOperator(a)
	}
	if !a.isExpired() {
		return
	}
	a.publish(message{
	// TODO pack message
	})
	a.ability.perform(a.unit, a.receiver)
	a.detachOperator(a)
}
