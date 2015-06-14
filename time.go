package main

import (
	"time"
)

type (
	GameTime     int64
	GameDuration int64
)

const (
	RealGameTick         = time.Second / 20
	RealPeriodicalTick   = time.Second / 2
	RealRegenerationTick = time.Second * 5

	GameTick         = GameDuration(RealGameTick / RealGameTick)
	PeriodicalTick   = GameDuration(RealPeriodicalTick / RealGameTick)
	RegenerationTick = GameDuration(RealRegenerationTick / RealGameTick)

	Second = GameDuration(time.Second / RealGameTick)
)

type GameClock interface {
	Now() GameTime
	Add(GameDuration) GameTime
	Before(GameTime) bool
	After(GameTime) bool
}
