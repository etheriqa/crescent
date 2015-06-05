package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestThreatOnAttach(t *testing.T) {
	assert := assert.New(t)
	g := mockGame()
	u := newUnit(g)
	target1 := newUnit(g)
	target2 := newUnit(g)
	threat1 := &threat{
		unit:   u,
		target: target1,
		threat: 10,
	}
	threat2 := &threat{
		unit:   u,
		target: target1,
		threat: 10,
	}
	threat3 := &threat{
		unit:   u,
		target: target2,
		threat: 100,
	}
	u.attachOperator(&disable{
		partialOperator: partialOperator{
			unit: u,
		},
	})
	u.attachOperator(threat1)
	u.attachOperator(threat2)
	u.attachOperator(threat3)
	assert.False(u.operators[threat1])
	assert.True(u.operators[threat2])
	assert.True(u.operators[threat3])
	assert.EqualValues(20, threat2.threat)
	assert.EqualValues(100, threat3.threat)
}
