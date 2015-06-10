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
	perform            func(performer, receiver *unit)
}

// checkRequirements checks the ability requirements are satisfied
func (a *ability) checkRequirements(performer *unit, receiver *unit) error {
	if err := a.checkReceiver(performer, receiver); err != nil {
		return err
	}
	if err := a.checkCooldown(performer); err != nil {
		return err
	}
	if err := a.checkDisable(performer); err != nil {
		return err
	}
	if err := a.checkCost(performer); err != nil {
		return err
	}
	return nil
}

// checkReceiver checks the receiver is valid
func (a *ability) checkReceiver(performer, receiver *unit) error {
	switch a.targetType {
	case targetTypeFriend:
		if receiver == nil || performer.group != receiver.group {
			return errors.New("The receiver must be friend")
		}
	case targetTypeEnemy:
		if receiver == nil || performer.group == receiver.group {
			return errors.New("The receiver must be enemy")
		}
	}
	return nil
}

// checkCooldown checks the performer does not have to wait the cooldown time expiration
func (a *ability) checkCooldown(performer *unit) error {
	for o := range performer.operators {
		switch o := o.(type) {
		case *cooldown:
			if o.ability != a {
				return errors.New("The performer has to wait the cooldown time expiration")
			}
		}
	}
	return nil
}

// checkDisable checks the performer is not interrupted by the disables
func (a *ability) checkDisable(performer *unit) error {
	for o := range performer.operators {
		switch o := o.(type) {
		case *disable:
			for d := range a.disableTypes {
				if disableType(d) == o.disableType {
					return errors.New("The performer is interrupted by the disable")
				}
			}
		}
	}
	return nil
}

// checkCost checks the performer satisfies the ability cost
func (a *ability) checkCost(performer *unit) error {
	if performer.health() < a.healthCost {
		return errors.New("The performer does not have enough health")
	}
	if performer.mana() < a.manaCost {
		return errors.New("The performer does not have enough mana")
	}
	return nil
}
