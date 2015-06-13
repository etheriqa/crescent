package main

type Cooldown struct {
	UnitObject
	ability        *Ability
	expirationTime GameTime

	clock    GameClock
	handlers HandlerContainer
	writer   GameEventWriter
}

// Ability returns the Ability
func (h *Cooldown) Ability() *Ability {
	return h.ability
}

// OnAttach removes other Cooldown handlers
func (h *Cooldown) OnAttach() {
	h.handlers.BindObject(h).Each(func(o Handler) {
		switch o := o.(type) {
		case *Cooldown:
			if h == o {
				return
			}
			h.handlers.Detach(o)
		}
	})

	if h.ability.CooldownDuration == 0 {
		h.handlers.Detach(h)
		return
	}

	h.expirationTime = h.clock.Add(h.ability.CooldownDuration)
	h.Object().AddEventHandler(h, EventGameTick)
	h.writer.Write(nil) // TODO
}

// OnDetach does nothing
func (h *Cooldown) OnDetach() {
	h.Object().RemoveEventHandler(h, EventGameTick)
}

// HandleEvent handles the Event
func (h *Cooldown) HandleEvent(e Event) {
	switch e {
	case EventGameTick:
		if h.clock.Before(h.expirationTime) {
			return
		}
		h.handlers.Detach(h)
		h.writer.Write(nil) // TODO
	}
}
