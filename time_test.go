package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInstanceClock(t *testing.T) {
	assert := assert.New(t)
	clock := InstanceTime(1000)

	assert.Implements((*InstanceClock)(nil), clock)

	assert.Equal(clock, clock.Now())

	assert.Equal(InstanceTime(1100), clock.Add(100))
	assert.Equal(InstanceTime(1200), clock.Add(200))

	assert.True(clock.After(900))
	assert.False(clock.After(1000))
	assert.False(clock.After(1100))

	assert.False(clock.Before(900))
	assert.False(clock.Before(1000))
	assert.True(clock.Before(1100))

	assert.True(InstanceTime(0).IsPeriodicalTick())
	assert.False(InstanceTime(1).IsPeriodicalTick())
	assert.False(InstanceTime(2).IsPeriodicalTick())
	assert.True(InstanceTime(PeriodicalTick).IsPeriodicalTick())
	assert.False(InstanceTime(PeriodicalTick*100 - 1).IsPeriodicalTick())
	assert.True(InstanceTime(PeriodicalTick * 100).IsPeriodicalTick())
	assert.False(InstanceTime(PeriodicalTick*100 + 1).IsPeriodicalTick())

	assert.True(InstanceTime(0).IsRegenerationTick())
	assert.False(InstanceTime(1).IsRegenerationTick())
	assert.False(InstanceTime(2).IsRegenerationTick())
	assert.True(InstanceTime(RegenerationTick).IsRegenerationTick())
	assert.False(InstanceTime(RegenerationTick*100 - 1).IsRegenerationTick())
	assert.True(InstanceTime(RegenerationTick * 100).IsRegenerationTick())
	assert.False(InstanceTime(RegenerationTick*100 + 1).IsRegenerationTick())
}
