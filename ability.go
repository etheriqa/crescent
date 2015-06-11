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

type ability struct {
	name               string
	TargetType         TargetType
	healthCost         Statistic
	manaCost           Statistic
	activationDuration GameDuration
	cooldownDuration   GameDuration
	disableTypes       []DisableType
	Perform            func(up UnitPair)
}

// CheckRequirements checks the ability requirements are satisfied
func (a *ability) CheckRequirements(subject *Unit, object *Unit) error {
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
func (a *ability) CheckObject(subject, object *Unit) error {
	switch a.TargetType {
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
func (a *ability) CheckCooldown(subject *Unit) error {
	ok := subject.EverySubjectHandler(subject, func(ha Handler) bool {
		switch ha := ha.(type) {
		case *Cooldown:
			if ha.ability == a {
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
func (a *ability) CheckDisable(subject *Unit) error {
	ok := subject.EverySubjectHandler(subject, func(ha Handler) bool {
		switch ha := ha.(type) {
		case *Disable:
			for dt := range a.disableTypes {
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
func (a *ability) CheckResource(subject *Unit) error {
	if subject.health() < a.healthCost {
		return errors.New("The subject does not have enough health")
	}
	if subject.mana() < a.manaCost {
		return errors.New("The subject does not have enough mana")
	}
	return nil
}
