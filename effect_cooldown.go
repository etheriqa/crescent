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

// EffectDidAttach removes other Cooldown effects
func (h *Cooldown) EffectDidAttach() error {
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
		return nil
	}

	h.Object().Register(h)
	return nil
}

// EffectDidDetach does nothing
func (h *Cooldown) EffectDidDetach() error {
	h.Object().Unregister(h)
	return nil
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
