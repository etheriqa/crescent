package main

import (
	"errors"
)

type Activating struct {
	UnitSubject
	object         *Unit
	ability        *Ability
	expirationTime InstanceTime

	handler EventHandler
}

// NewActivating returns a Activating
func NewActivating(g Game, s Subject, o *Unit, a *Ability, t InstanceTime) *Activating {
	e := &Activating{
		UnitSubject:    MakeSubject(s),
		object:         o,
		ability:        a,
		expirationTime: t,
	}
	e.handler = MakeEventHandler(func(p interface{}) { e.handle(g, p) })
	return e
}

// MaybeObject returns the object unit or nil
func (e *Activating) MaybeObject() *Unit {
	return e.object
}

// Ability returns the Ability
func (e *Activating) Ability() *Ability {
	return e.ability
}

// EffectWillAttach checks requirements
func (e *Activating) EffectWillAttach(g Game) error {
	ok := g.EffectQuery().BindSubject(e).Every(func(f Effect) bool {
		switch f.(type) {
		case *Activating:
			return false
		}
		return true
	})
	if !ok {
		return errors.New("Already activating")
	}

	if err := e.checkRequirements(g); err != nil {
		return err
	}

	return nil
}

// EffectDidAttach performs ability if it has no activation duration
func (e *Activating) EffectDidAttach(g Game) error {
	if e.ability.ActivationDuration == 0 {
		e.perform(g)
		return nil
	}

	e.Subject().Register(e.handler)
	if e.object != nil {
		e.object.Register(e.handler)
	}

	e.writeOutputUnitActivating(g)

	return nil
}

// EffectDidDetach does nothing
func (e *Activating) EffectDidDetach(g Game) error {
	e.Subject().Unregister(e.handler)
	if e.object != nil {
		e.object.Unregister(e.handler)
	}
	return nil
}

// handle handles the payload
func (e *Activating) handle(g Game, p interface{}) {
	switch p.(type) {
	case EventGameTick:
		if g.Clock().Before(e.expirationTime) {
			return
		}
		e.perform(g)
	case EventDead:
		if e.Subject().IsAlive() {
			e.writeOutputUnitActivated(g, false)
		}
		g.DetachEffect(e)
	case EventDisabled:
		if err := e.checkDisable(g); err == nil {
			return
		}
		e.writeOutputUnitActivated(g, false)
		g.DetachEffect(e)
	case EventTakenDamage:
		if err := e.checkResource(); err == nil {
			return
		}
		e.writeOutputUnitActivated(g, false)
		g.DetachEffect(e)
	}
}

// perform performs the Ability
func (e *Activating) perform(g Game) {
	e.writeOutputUnitActivated(g, true)
	if _, _, err := e.Subject().ModifyHealth(g.Writer(), -e.ability.HealthCost); err != nil {
		log.Fatal(err)
	}
	if _, _, err := e.Subject().ModifyMana(g.Writer(), -e.ability.ManaCost); err != nil {
		log.Fatal(err)
	}
	e.ability.Perform(g, e.Subject(), e.object)
	g.DetachEffect(e)
	g.Cooldown(e.Subject(), e.ability)
}

// checkRequirements checks all requirements
func (e *Activating) checkRequirements(g Game) error {
	if err := e.checkObject(); err != nil {
		return err
	}
	if err := e.checkCooldown(g); err != nil {
		return err
	}
	if err := e.checkDisable(g); err != nil {
		return err
	}
	if err := e.checkResource(); err != nil {
		return err
	}
	return nil
}

// checkObject checks the Object is valid
func (e *Activating) checkObject() error {
	switch e.ability.TargetType {
	case TargetTypeNone:
		// TODO WIP
		/*
			if h.object != nil {
				return errors.New("The Object must be nil")
			}
		*/
	case TargetTypeFriend:
		if e.object == nil {
			return errors.New("The Object must be *Unit")
		}
		if e.object.Group() != e.Subject().Group() {
			return errors.New("The Object must be friend")
		}
		if e.object.IsDead() {
			return errors.New("The Object must be alive")
		}
	case TargetTypeEnemy:
		if e.object == nil {
			return errors.New("The Object must be *Unit")
		}
		if e.object.Group() == e.Subject().Group() {
			return errors.New("The Object must be enemy")
		}
		if e.object.IsDead() {
			return errors.New("The Object must be alive")
		}
	default:
		return errors.New("Unknown TargetType")
	}
	return nil
}

// checkCooldown checks the Subject does not have to wait the Cooldown
func (e *Activating) checkCooldown(g Game) error {
	ok := g.EffectQuery().BindObject(e.Subject()).Every(func(f Effect) bool {
		switch f := f.(type) {
		case *Cooldown:
			if e.ability == f.Ability() {
				return false
			}
		}
		return true
	})
	if ok {
		return nil
	}
	return errors.New("The Object has to wait the Cooldown")
}

// checkDisable checks the Subject has not been interrupted by Disables
func (e *Activating) checkDisable(g Game) error {
	ok := g.EffectQuery().BindObject(e.Subject()).Every(func(f Effect) bool {
		switch f := f.(type) {
		case *Disable:
			for _, t := range e.ability.DisableTypes {
				if f.disableType == t {
					return false
				}
			}
		}
		return true
	})
	if ok {
		return nil
	}
	return errors.New("The Object has been interrupted by Disables")
}

// checkResource checks the Subject has enough resource
func (e *Activating) checkResource() error {
	if e.Subject().Health() <= e.ability.HealthCost {
		return errors.New("The Subject does not have enough health")
	}
	if e.Subject().Mana() < e.ability.ManaCost {
		return errors.New("The Subject does not have enough mana")
	}
	return nil
}

// writeOutputUnitActivating writes a OutputUnitActivating
func (e *Activating) writeOutputUnitActivating(g Game) {
	g.Writer().Write(OutputUnitActivating{
		UnitID:      e.Subject().ID(),
		AbilityName: e.ability.Name,
		StartTime:   g.Clock().Now(),
		EndTime:     e.expirationTime,
	})
}

// writeOutputUnitActivated writes a OutputUnitActivated
func (e *Activating) writeOutputUnitActivated(g Game, ok bool) {
	g.Writer().Write(OutputUnitActivated{
		UnitID:      e.Subject().ID(),
		AbilityName: e.ability.Name,
		OK:          ok,
	})
}
