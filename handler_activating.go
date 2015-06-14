package main

import (
	"errors"
)

type Activating struct {
	UnitSubject
	object         *Unit
	ability        *Ability
	expirationTime GameTime

	op Operator
}

// Ability returns the Ability
func (h *Activating) Ability() *Ability {
	return h.ability
}

// OnAttach checks requirements
func (h *Activating) OnAttach() {
	ok := h.op.Handlers().BindSubject(h).Every(func(o Handler) bool {
		switch o.(type) {
		case *Activating:
			if h != o {
				return false
			}
		}
		return true
	})
	if !ok {
		h.op.Handlers().Detach(h)
		return
	}

	if err := h.checkRequirements(); err != nil {
		log.Debug(err)
		h.op.Handlers().Detach(h)
		return
	}

	if h.ability.ActivationDuration == 0 {
		h.perform()
		return
	}

	h.Subject().AddEventHandler(h, EventGameTick)
	h.Subject().AddEventHandler(h, EventDead)
	h.Subject().AddEventHandler(h, EventDisabled)
	h.Subject().AddEventHandler(h, EventTakenDamage)
	if h.object != nil {
		h.object.AddEventHandler(h, EventDead)
	}
	h.op.Writer().Write(nil) // TODO
}

// OnDetach does nothing
func (h *Activating) OnDetach() {
	h.Subject().RemoveEventHandler(h, EventGameTick)
	h.Subject().RemoveEventHandler(h, EventDead)
	h.Subject().RemoveEventHandler(h, EventDisabled)
	h.Subject().RemoveEventHandler(h, EventTakenDamage)
	if h.object != nil {
		h.object.RemoveEventHandler(h, EventDead)
	}
}

// HandleEvent handles the Event
func (h *Activating) HandleEvent(e Event) {
	switch e {
	case EventGameTick:
		if h.op.Clock().Before(h.expirationTime) {
			return
		}
		h.perform()
	case EventDead:
		h.op.Handlers().Detach(h)
		if h.Subject().IsAlive() {
			h.op.Writer().Write(nil) // TODO
		}
	case EventDisabled:
		if err := h.checkDisable(); err != nil {
			h.op.Handlers().Detach(h)
			h.op.Writer().Write(nil) // TODO
		}
	case EventTakenDamage:
		if err := h.checkResource(); err != nil {
			h.op.Handlers().Detach(h)
			h.op.Writer().Write(nil) // TODO
		}
	}
}

// perform performs the Ability
func (h *Activating) perform() {
	if _, _, err := h.Subject().ModifyHealth(h.op.Writer(), -h.ability.HealthCost); err != nil {
		log.Fatal(err)
	}
	if _, _, err := h.Subject().ModifyMana(h.op.Writer(), -h.ability.ManaCost); err != nil {
		log.Fatal(err)
	}
	h.ability.Perform(h.op, h.Subject(), h.object)
	// TODO perform ability
	h.op.Handlers().Detach(h)
	h.op.Cooldown(h.Subject(), h.ability)
	h.op.Writer().Write(nil) // TODO
}

// checkRequirements checks all requirements
func (h *Activating) checkRequirements() error {
	if err := h.checkObject(); err != nil {
		return err
	}
	if err := h.checkCooldown(); err != nil {
		return err
	}
	if err := h.checkDisable(); err != nil {
		return err
	}
	if err := h.checkResource(); err != nil {
		return err
	}
	return nil
}

// checkObject checks the Object is valid
func (h *Activating) checkObject() error {
	switch h.ability.TargetType {
	case TargetTypeNone:
		if h.object != nil {
			return errors.New("The Object must be nil")
		}
	case TargetTypeFriend:
		if h.object == nil {
			return errors.New("The Object must be *Unit")
		}
		if h.object.Group() != h.Subject().Group() {
			return errors.New("The Object must be friend")
		}
		if h.object.IsDead() {
			return errors.New("The Object must be alive")
		}
	case TargetTypeEnemy:
		if h.object == nil {
			return errors.New("The Object must be *Unit")
		}
		if h.object.Group() == h.Subject().Group() {
			return errors.New("The Object must be enemy")
		}
		if h.object.IsDead() {
			return errors.New("The Object must be alive")
		}
	default:
		return errors.New("Unknown TargetType")
	}
	return nil
}

// checkCooldown checks the Subject does not have to wait the Cooldown
func (h *Activating) checkCooldown() error {
	ok := h.op.Handlers().BindObject(h.Subject()).Every(func(o Handler) bool {
		switch o := o.(type) {
		case *Cooldown:
			if h.ability == o.Ability() {
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
func (h *Activating) checkDisable() error {
	ok := h.op.Handlers().BindObject(h.Subject()).Every(func(o Handler) bool {
		switch o := o.(type) {
		case *Disable:
			for _, t := range h.ability.DisableTypes {
				if o.disableType == t {
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
func (h *Activating) checkResource() error {
	if h.Subject().Health() <= h.ability.HealthCost {
		return errors.New("The Subject does not have enough health")
	}
	if h.Subject().Mana() < h.ability.ManaCost {
		return errors.New("The Subject does not have enough mana")
	}
	return nil
}
