package crescent

type TargetType uint8

const (
	_ TargetType = iota
	TargetTypeNone
	TargetTypeFriend
	TargetTypeEnemy
)

type Ability struct {
	Name               string
	Description        string
	TargetType         TargetType
	HealthCost         Statistic
	ManaCost           Statistic
	ActivationDuration InstanceDuration
	CooldownDuration   InstanceDuration
	DisableTypes       []DisableType
	Perform            func(Game, Subject, *Unit)
}
