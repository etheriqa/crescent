package main

import (
	"errors"
)

type Activating struct {
	UnitSubject
	object         *Unit
	ability        *Ability
	expirationTime InstanceTime

	op Operator
}

// Ability returns the Ability
func (h *Activating) Ability() *Ability {
	return h.ability
}

// OnAttach checks requirements
func (h *Activating) OnAttach() {
	ok := h.op.Effects().BindSubject(h).Every(func(o Effect) bool {
		switch o.(type) {
		case *Activating:
			if h != o {
				return false
			}
		}
		return true
	})
	if !ok {
		h.op.Effects().Detach(h)
		return
	}

	if err := h.checkRequirements(); err != nil {
		log.Debug(err)
		h.op.Effects().Detach(h)
		return
	}

	if h.ability.ActivationDuration == 0 {
		h.perform()
		return
	}

	h.Subject().Register(h)
	if h.object != nil {
		h.object.Register(h)
	}

	h.writeOutputUnitActivating()
}

// OnDetach does nothing
func (h *Activating) OnDetach() {
	h.Subject().Unregister(h)
	if h.object != nil {
		h.object.Unregister(h)
	}
}

// Handle handles the Event
func (h *Activating) Handle(p interface{}) {
	switch p.(type) {
	case *EventGameTick:
		if h.op.Clock().Before(h.expirationTime) {
			return
		}
		h.perform()
	case *EventDead:
		h.writeOutputUnitActivated(false)
		h.op.Effects().Detach(h)
		if h.Subject().IsDead() {
			return
		}
	case *EventDisabled:
		if err := h.checkDisable(); err == nil {
			return
		}
		h.writeOutputUnitActivated(false)
		h.op.Effects().Detach(h)
	case *EventTakenDamage:
		if err := h.checkResource(); err == nil {
			return
		}
		h.writeOutputUnitActivated(false)
		h.op.Effects().Detach(h)
	}
}

// perform performs the Ability
func (h *Activating) perform() {
	h.writeOutputUnitActivated(true)
	if _, _, err := h.Subject().ModifyHealth(h.op.Writer(), -h.ability.HealthCost); err != nil {
		log.Fatal(err)
	}
	if _, _, err := h.Subject().ModifyMana(h.op.Writer(), -h.ability.ManaCost); err != nil {
		log.Fatal(err)
	}
	h.ability.Perform(h.op, h.Subject(), h.object)
	h.op.Effects().Detach(h)
	h.op.Cooldown(h.Subject(), h.ability)
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
		// TODO WIP
		/*
			if h.object != nil {
				return errors.New("The Object must be nil")
			}
		*/
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
	ok := h.op.Effects().BindObject(h.Subject()).Every(func(o Effect) bool {
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
	ok := h.op.Effects().BindObject(h.Subject()).Every(func(o Effect) bool {
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

// writeOutputUnitActivating writes a OutputUnitActivating
func (h *Activating) writeOutputUnitActivating() {
	h.op.Writer().Write(OutputUnitActivating{
		UnitID:      h.Subject().ID(),
		AbilityName: h.ability.Name,
		StartTime:   h.op.Clock().Now(),
		EndTime:     h.expirationTime,
	})
}

// writeOutputUnitActivated writes a OutputUnitActivated
func (h *Activating) writeOutputUnitActivated(ok bool) {
	h.op.Writer().Write(OutputUnitActivated{
		UnitID:      h.Subject().ID(),
		AbilityName: h.ability.Name,
		OK:          ok,
	})
}
