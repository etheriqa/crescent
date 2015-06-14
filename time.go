package main

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
	Before(InstanceTime) bool
	After(InstanceTime) bool
}
