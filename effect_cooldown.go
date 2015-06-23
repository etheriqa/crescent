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

// OnAttach removes other Cooldown effects
func (h *Cooldown) OnAttach() {
	h.op.Effects().BindObject(h).Each(func(o Effect) {
		switch o := o.(type) {
		case *Cooldown:
			if h == o {
				return
			}
			if h.ability != o.ability {
				return
			}
			h.op.Effects().Detach(o)
		}
	})

	h.writeOutputUnitCooldown()

	if !h.isActive() {
		h.op.Effects().Detach(h)
		return
	}

	h.Object().Register(h)
}

// OnDetach does nothing
func (h *Cooldown) OnDetach() {
	h.Object().Unregister(h)
}

// Handle handles the Event
func (h *Cooldown) Handle(p interface{}) {
	switch p.(type) {
	case *EventGameTick:
		if h.isActive() {
			return
		}
		h.writeOutputUnitCooldown()
		h.op.Effects().Detach(h)
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
