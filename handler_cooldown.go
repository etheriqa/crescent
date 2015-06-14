package main

type Cooldown struct {
	UnitObject
	ability        *Ability
	expirationTime GameTime

	op Operator
}

// Ability returns the Ability
func (h *Cooldown) Ability() *Ability {
	return h.ability
}

// OnAttach removes other Cooldown handlers
func (h *Cooldown) OnAttach() {
	h.op.Handlers().BindObject(h).Each(func(o Handler) {
		switch o := o.(type) {
		case *Cooldown:
			if h == o {
				return
			}
			h.op.Handlers().Detach(o)
		}
	})

	if h.ability.CooldownDuration == 0 {
		h.op.Handlers().Detach(h)
		return
	}

	h.expirationTime = h.op.Clock().Add(h.ability.CooldownDuration)
	h.Object().AddEventHandler(h, EventGameTick)
	h.op.Writer().Write(nil) // TODO
}

// OnDetach does nothing
func (h *Cooldown) OnDetach() {
	h.Object().RemoveEventHandler(h, EventGameTick)
}

// HandleEvent handles the Event
func (h *Cooldown) HandleEvent(e Event) {
	switch e {
	case EventGameTick:
		if h.op.Clock().Before(h.expirationTime) {
			return
		}
		h.op.Handlers().Detach(h)
		h.op.Writer().Write(nil) // TODO
	}
}
