package crescent

import (
	"math/rand"
)

type Game interface {
	Rand() *rand.Rand
	Clock() InstanceClock
	Writer() InstanceOutputWriter

	Join(UnitGroup, UnitName, *Class) (UnitID, error)
	Leave(UnitID) error
	UnitQuery() UnitQueryable

	AttachEffect(Effect) error
	DetachEffect(Effect) error
	EffectQuery() EffectQueryable

	Activating(Subject, *Unit, *Ability)
	Cooldown(Object, *Ability)
	ResetCooldown(Object, *Ability)
	Correction(Object, UnitCorrection, string, Statistic, InstanceDuration)
	Disable(Object, DisableType, InstanceDuration)
	DamageThreat(Subject, Object, Statistic)
	HealingThreat(Subject, Object, Statistic)
	DoT(Damage, string, InstanceDuration)
	HoT(Healing, string, InstanceDuration)
}
