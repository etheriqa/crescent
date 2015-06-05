package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestThreatOnAttach(t *testing.T) {
	assert := assert.New(t)
	u := newUnit()
	target1 := newUnit()
	target2 := newUnit()
	threat1 := &threat{
		target: target1,
		v:      10,
	}
	threat2 := &threat{
		target: target1,
		v:      10,
	}
	threat3 := &threat{
		target: target2,
		v:      100,
	}
	u.attachOperator(&disable{
		u: u,
		o: make(chan message, 100),
	})
	u.attachOperator(threat1)
	u.attachOperator(threat2)
	u.attachOperator(threat3)
	assert.False(u.operators[threat1])
	assert.True(u.operators[threat2])
	assert.True(u.operators[threat3])
	assert.EqualValues(20, threat2.v)
	assert.EqualValues(100, threat3.v)
}
