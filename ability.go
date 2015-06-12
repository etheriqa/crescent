package main

import (
	"errors"
)

type TargetType uint8

const (
	_ TargetType = iota
	TargetTypeNone
	TargetTypeFriend
	TargetTypeEnemy
)

type Ability struct {
	Name               string
	TargetType         TargetType
	HealthCost         Statistic
	ManaCost           Statistic
	ActivationDuration GameDuration
	CooldownDuration   GameDuration
	DisableTypes       []DisableType
	Perform            func(up UnitPair)
}

// CheckRequirements checks the ability requirements are satisfied
func (a *Ability) CheckRequirements(subject *Unit, object *Unit) error {
	if err := a.CheckObject(subject, object); err != nil {
		return err
	}
	if err := a.CheckCooldown(subject); err != nil {
		return err
	}
	if err := a.CheckDisable(subject); err != nil {
		return err
	}
	if err := a.CheckResource(subject); err != nil {
		return err
	}
	return nil
}

// CheckObject checks the object is valid
func (a *Ability) CheckObject(subject, object *Unit) error {
	switch a.TargetType {
	case TargetTypeNone:
		if object != nil {
			return errors.New("The object must be nil")
		}
	case TargetTypeFriend:
		if object == nil || subject.group != object.group {
			return errors.New("The object must be friend")
		}
	case TargetTypeEnemy:
		if object == nil || subject.group == object.group {
			return errors.New("The object must be enemy")
		}
	}
	return nil
}

// CheckCooldown checks the subject does not have to wait the cooldown time expiration
func (a *Ability) CheckCooldown(subject *Unit) error {
	ok := subject.EverySubjectHandler(func(ha Handler) bool {
		switch ha := ha.(type) {
		case *Cooldown:
			if ha.Ability() == a {
				return false
			}
		}
		return true
	})
	if ok {
		return nil
	}
	return errors.New("The subject has to wait the cooldown time expiration")
}

// CheckDisable checks the subject is not interrupted by the disables
func (a *Ability) CheckDisable(subject *Unit) error {
	ok := subject.EverySubjectHandler(func(ha Handler) bool {
		switch ha := ha.(type) {
		case *Disable:
			for dt := range a.DisableTypes {
				if DisableType(dt) == ha.DisableType {
					return false
				}
			}
		}
		return true
	})
	if ok {
		return nil
	}
	return errors.New("The subject is interrupted by the disable")
}

// CheckResource checks the subject satisfies the ability cost
func (a *Ability) CheckResource(subject *Unit) error {
	if subject.Health() < a.HealthCost {
		return errors.New("The subject does not have enough health")
	}
	if subject.Mana() < a.ManaCost {
		return errors.New("The subject does not have enough mana")
	}
	return nil
}
