package main

type activating struct {
	partialOperator
	*ability
	receiver *unit
}

// onAttach checks requirements
func (a *activating) onAttach() {
	a.AddEventHandler(a, EventDead)
	a.AddEventHandler(a, EventDisableInterrupt)
	a.AddEventHandler(a, EventGameTick)
	a.AddEventHandler(a, EventResourceDecreased)
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

// onDetach removes the EventHandlers
func (a *activating) onDetach() {
	a.RemoveEventHandler(a, EventDead)
	a.RemoveEventHandler(a, EventDisableInterrupt)
	a.RemoveEventHandler(a, EventGameTick)
	a.RemoveEventHandler(a, EventResourceDecreased)
}

// HandleEvent handles the event
func (a *activating) HandleEvent(e Event) {
	switch e {
	case EventDead:
		a.detachOperator(a)
	case EventDisableInterrupt:
		a.perform()
	case EventGameTick:
		a.perform()
	case EventResourceDecreased:
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
