package main

type EventGameTick struct{}

type EventPeriodicalTick struct{}

type EventTakenDamage struct{}

type EventDisabled struct{}

type EventDead struct{}

type EventInterrupt struct {
	UnitID
}
