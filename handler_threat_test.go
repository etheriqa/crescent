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
	threat1 := &Threat{
		partialHandler: partialHandler{
			unit:      u,
			performer: performer1,
		},
		threat: 10,
	}
	threat2 := &Threat{
		partialHandler: partialHandler{
			unit:      u,
			performer: performer1,
		},
		threat: 10,
	}
	threat3 := &Threat{
		partialHandler: partialHandler{
			unit:      u,
			performer: performer2,
		},
		threat: 100,
	}
	u.attachHandler(&Disable{
		partialHandler: partialHandler{
			unit: u,
		},
	})
	u.attachHandler(threat1)
	u.attachHandler(threat2)
	u.attachHandler(threat3)
	assert.False(u.handlers[threat1])
	assert.True(u.handlers[threat2])
	assert.True(u.handlers[threat3])
	assert.EqualValues(20, threat2.threat)
	assert.EqualValues(100, threat3.threat)
}
