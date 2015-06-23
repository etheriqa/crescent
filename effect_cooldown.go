package main

type Cooldown struct {
	UnitObject
	ability        *Ability
	expirationTime InstanceTime

	handler EventHandler
}

// NewCooldown returns a Cooldown
func NewCooldown(g Game, o Object, a *Ability, t InstanceTime) *Cooldown {
	e := &Cooldown{
		UnitObject:     MakeObject(o),
		ability:        a,
		expirationTime: t,
		handler:        new(func(interface{})),
	}
	*e.handler = func(p interface{}) { e.handle(g, p) }
	return e
}

// Ability returns the Ability
func (e *Cooldown) Ability() *Ability {
	return e.ability
}

// EffectWillAttach removes other Cooldown effects
func (e *Cooldown) EffectWillAttach(g Game) error {
	g.EffectQuery().BindObject(e).Each(func(f Effect) {
		switch f := f.(type) {
		case *Cooldown:
			if e.ability != f.ability {
				return
			}
			g.DetachEffect(f)
		}
	})
	return nil
}

// EffectDidAttach detaches itself if it is not active
func (e *Cooldown) EffectDidAttach(g Game) error {
	e.writeOutputUnitCooldown(g)

	if !e.isActive(g) {
		g.DetachEffect(e)
		return nil
	}

	e.Object().Register(e.handler)
	return nil
}

// EffectDidDetach does nothing
func (e *Cooldown) EffectDidDetach(g Game) error {
	e.Object().Unregister(e.handler)
	return nil
}

// handle handles the payload
func (e *Cooldown) handle(g Game, p interface{}) {
	switch p.(type) {
	case EventGameTick:
		if e.isActive(g) {
			return
		}
		e.writeOutputUnitCooldown(g)
		g.DetachEffect(e)
	}
}

// isActive returns true if the Cooldown is active
func (e *Cooldown) isActive(g Game) bool {
	return g.Clock().Before(e.expirationTime)
}

// writeOutputUnitCooldown writes a OutputUnitCooldown
func (e *Cooldown) writeOutputUnitCooldown(g Game) {
	// TODO write to only the object
	g.Writer().Write(OutputUnitCooldown{
		UnitID:         e.Object().ID(),
		AbilityName:    e.ability.Name,
		ExpirationTime: e.expirationTime,
		Active:         e.isActive(g),
	})
}
