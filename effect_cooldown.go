package main

type Cooldown struct {
	UnitObject
	ability        *Ability
	expirationTime InstanceTime

	g Game
}

// Ability returns the Ability
func (h *Cooldown) Ability() *Ability {
	return h.ability
}

// EffectDidAttach removes other Cooldown effects
func (h *Cooldown) EffectDidAttach(g Game) error {
	h.g.EffectQuery().BindObject(h).Each(func(o Effect) {
		switch o := o.(type) {
		case *Cooldown:
			if h == o {
				return
			}
			if h.ability != o.ability {
				return
			}
			h.g.DetachEffect(o)
		}
	})

	h.writeOutputUnitCooldown()

	if !h.isActive() {
		h.g.DetachEffect(h)
		return nil
	}

	h.Object().Register(h)
	return nil
}

// EffectDidDetach does nothing
func (h *Cooldown) EffectDidDetach(g Game) error {
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
		h.g.DetachEffect(h)
	}
}

// isActive returns true if the Cooldown is active
func (h *Cooldown) isActive() bool {
	return h.g.Clock().Before(h.expirationTime)
}

// writeOutputUnitCooldown writes a OutputUnitCooldown
func (h *Cooldown) writeOutputUnitCooldown() {
	// TODO write to only the object
	h.g.Writer().Write(OutputUnitCooldown{
		UnitID:         h.Object().ID(),
		AbilityName:    h.ability.Name,
		ExpirationTime: h.expirationTime,
		Active:         h.isActive(),
	})
}
