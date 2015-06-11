package main

type Activating struct {
	*PartialHandler
	*ability
}

// NewActivating returns a Activating handler
func NewActivating(subject, object *unit, ability *ability) *Activating {
	return &Activating{
		PartialHandler: NewPartialHandler(subject, object, ability.activationDuration),
		ability:        ability,
	}
}

// OnAttach checks requirements
func (a *Activating) OnAttach() {
	a.Subject().AddEventHandler(a, EventDead)
	a.Subject().AddEventHandler(a, EventDisableInterrupt)
	a.Subject().AddEventHandler(a, EventGameTick)
	a.Subject().AddEventHandler(a, EventResourceDecreased)
	if a.Object() != nil {
		a.Object().AddEventHandler(a, EventDead)
	}
	ok := a.Container().EverySubjectHandler(a.Subject(), func(ha Handler) bool {
		switch ha.(type) {
		case *Activating:
			return false
		}
		return true
	})
	if !ok {
		a.Stop(a)
		return
	}
	if err := a.checkRequirements(a.Subject(), a.Object()); err != nil {
		a.Stop(a)
		return
	}
	if a.IsExpired() {
		a.perform()
		return
	}
	a.Publish(message{
	// TODO pack message
	})
}

// OnDetach removes the EventHandlers
func (a *Activating) OnDetach() {
	a.Subject().RemoveEventHandler(a, EventDead)
	a.Subject().RemoveEventHandler(a, EventDisableInterrupt)
	a.Subject().RemoveEventHandler(a, EventGameTick)
	a.Subject().RemoveEventHandler(a, EventResourceDecreased)
	if a.Object() != nil {
		a.Object().RemoveEventHandler(a, EventDead)
	}
}

// HandleEvent handles the event
func (a *Activating) HandleEvent(e Event) {
	switch e {
	case EventDead:
		a.Stop(a)
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
	if err := a.checkRequirements(a.Subject(), a.Object()); err != nil {
		a.Publish(message{
		// TODO pack message
		})
		a.Stop(a)
	}
	if !a.IsExpired() {
		return
	}
	a.Publish(message{
	// TODO pack message
	})
	// TODO consume health
	// TODO consume mana
	a.ability.perform(a.Subject(), a.Object())
	a.Stop(a)
}
