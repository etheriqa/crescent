package main

type Activating struct {
	partialHandler
	*ability
	receiver *unit
}

// NewActivating returns a activating handler
func NewActivating(performer, receiver *unit, ability *ability) *Activating {
	return &Activating{
		partialHandler: partialHandler{
			unit:           performer,
			expirationTime: performer.after(ability.activationDuration),
		},
		ability:  ability,
		receiver: receiver,
	}
}

// OnAttach checks requirements
func (a *Activating) OnAttach() {
	a.AddEventHandler(a, EventDead)
	a.AddEventHandler(a, EventDisableInterrupt)
	a.AddEventHandler(a, EventGameTick)
	a.AddEventHandler(a, EventResourceDecreased)
	for ha := range a.handlers {
		switch ha.(type) {
		case *Activating:
			a.detachHandler(a)
			return
		}
	}
	if err := a.checkRequirements(a.unit, a.receiver); err != nil {
		a.detachHandler(a)
		return
	}
	if a.isExpired() {
		a.perform()
		return
	}
	a.publish(message{
	// TODO pack message
	})
}

// OnDetach removes the EventHandlers
func (a *Activating) OnDetach() {
	a.RemoveEventHandler(a, EventDead)
	a.RemoveEventHandler(a, EventDisableInterrupt)
	a.RemoveEventHandler(a, EventGameTick)
	a.RemoveEventHandler(a, EventResourceDecreased)
}

// HandleEvent handles the event
func (a *Activating) HandleEvent(e Event) {
	switch e {
	case EventDead:
		a.detachHandler(a)
	case EventDisableInterrupt:
		a.perform()
	case EventGameTick:
		a.perform()
	case EventResourceDecreased:
		a.perform()
	}
}

// perform performs the ability
func (a *Activating) perform() {
	if err := a.checkRequirements(a.unit, a.receiver); err != nil {
		a.publish(message{
		// TODO pack message
		})
		a.detachHandler(a)
	}
	if !a.isExpired() {
		return
	}
	a.publish(message{
	// TODO pack message
	})
	a.ability.perform(a.unit, a.receiver)
	a.detachHandler(a)
}
