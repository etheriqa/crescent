package main

type Cooldown struct {
	UnitObject
	ability        *Ability
	expirationTime InstanceTime

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
	h.writeOutputUnitCooldown()
	h.Object().AddEventHandler(h, EventGameTick)
}

// OnDetach does nothing
func (h *Cooldown) OnDetach() {
	h.Object().RemoveEventHandler(h, EventGameTick)
}

// HandleEvent handles the Event
func (h *Cooldown) HandleEvent(e Event) {
	switch e {
	case EventGameTick:
		if h.isActive() {
			return
		}
		h.writeOutputUnitCooldown()
		h.op.Handlers().Detach(h)
	}
}

// isActive returns true if the Cooldown is active
func (h *Cooldown) isActive() bool {
	return h.op.Clock().Before(h.expirationTime)
}

// writeOutputUnitCooldown writes a OutputUnitCooldown
func (h *Cooldown) writeOutputUnitCooldown() {
	// TODO write to only the object
	h.op.Writer().Write(OutputUnitCooldown{
		UnitID:         h.Object().ID(),
		AbilityName:    h.ability.Name,
		ExpirationTime: h.expirationTime,
		Active:         h.isActive(),
	})
}
