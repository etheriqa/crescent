package main

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
	Perform            func(Operator, Subject, *Unit)
}
