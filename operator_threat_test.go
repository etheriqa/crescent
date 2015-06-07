package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestThreatOnAttach(t *testing.T) {
	assert := assert.New(t)
	g := mockGame()
	c := &class{}
	u := newUnit(g, c)
	performer1 := newUnit(g, c)
	performer2 := newUnit(g, c)
	threat1 := &threat{
		partialOperator: partialOperator{
			unit:      u,
			performer: performer1,
		},
		threat: 10,
	}
	threat2 := &threat{
		partialOperator: partialOperator{
			unit:      u,
			performer: performer1,
		},
		threat: 10,
	}
	threat3 := &threat{
		partialOperator: partialOperator{
			unit:      u,
			performer: performer2,
		},
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
