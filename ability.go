package main

import (
	"errors"
)

type targetType uint8

const (
	_ targetType = iota
	targetTypeNone
	targetTypeFriend
	targetTypeEnemy
)

type ability struct {
	name               string
	targetType         targetType
	healthCost         statistic
	manaCost           statistic
	activationDuration gameDuration
	cooldownDuration   gameDuration
	disableTypes       []disableType
	perform            func(subject, object *unit)
}

// checkRequirements checks the ability requirements are satisfied
func (a *ability) checkRequirements(subject *unit, object *unit) error {
	if err := a.checkobject(subject, object); err != nil {
		return err
	}
	if err := a.checkCooldown(subject); err != nil {
		return err
	}
	if err := a.checkDisable(subject); err != nil {
		return err
	}
	if err := a.checkCost(subject); err != nil {
		return err
	}
	return nil
}

// checkobject checks the object is valid
func (a *ability) checkobject(subject, object *unit) error {
	switch a.targetType {
	case targetTypeFriend:
		if object == nil || subject.group != object.group {
			return errors.New("The object must be friend")
		}
	case targetTypeEnemy:
		if object == nil || subject.group == object.group {
			return errors.New("The object must be enemy")
		}
	}
	return nil
}

// checkCooldown checks the subject does not have to wait the cooldown time expiration
func (a *ability) checkCooldown(subject *unit) error {
	ok := subject.AllSubjectHandler(subject, func(ha Handler) bool {
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

// checkDisable checks the subject is not interrupted by the disables
func (a *ability) checkDisable(subject *unit) error {
	ok := subject.AllSubjectHandler(subject, func(ha Handler) bool {
		switch ha := ha.(type) {
		case *Disable:
			for dt := range a.disableTypes {
				if disableType(dt) == ha.disableType {
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

// checkCost checks the subject satisfies the ability cost
func (a *ability) checkCost(subject *unit) error {
	if subject.health() < a.healthCost {
		return errors.New("The subject does not have enough health")
	}
	if subject.mana() < a.manaCost {
		return errors.New("The subject does not have enough mana")
	}
	return nil
}
