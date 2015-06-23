package crescent

import (
	"time"
)

type (
	InstanceTime     int64
	InstanceDuration int64
)

const (
	RealGameTick         = time.Second / 20
	RealPeriodicalTick   = time.Second / 2
	RealRegenerationTick = time.Second * 5

	GameTick         = InstanceDuration(RealGameTick / RealGameTick)
	PeriodicalTick   = InstanceDuration(RealPeriodicalTick / RealGameTick)
	RegenerationTick = InstanceDuration(RealRegenerationTick / RealGameTick)

	Second = InstanceDuration(time.Second / RealGameTick)
)

type InstanceClock interface {
	Now() InstanceTime
	Add(InstanceDuration) InstanceTime
	After(InstanceTime) bool
	Before(InstanceTime) bool
}

// Now returns the InstanceTime
func (t InstanceTime) Now() InstanceTime {
	return t
}

// Add returns the InstanceTime t+d
func (t InstanceTime) Add(d InstanceDuration) InstanceTime {
	return t + InstanceTime(d)
}

// Before returns true if t is after u
func (t InstanceTime) After(u InstanceTime) bool {
	return t > u
}

// After returns true if t is before u
func (t InstanceTime) Before(u InstanceTime) bool {
	return t < u
}

// IsPeriodicalTick returns true if the time t is periodical tick
func (t InstanceTime) IsPeriodicalTick() bool {
	return InstanceDuration(t)%PeriodicalTick == 0
}

// IsRegenerationTick returns true if the time t is regeneration tick
func (t InstanceTime) IsRegenerationTick() bool {
	return InstanceDuration(t)%RegenerationTick == 0
}
