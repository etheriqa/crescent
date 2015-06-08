package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDamageReductionFactor(t *testing.T) {
	assert := assert.New(t)
	assert.InDelta(1.000000, float64(damageReductionFactor(0)), 1e-6)
	assert.InDelta(0.909090, float64(damageReductionFactor(10)), 1e-6)
	assert.InDelta(0.833333, float64(damageReductionFactor(20)), 1e-6)
	assert.InDelta(0.800000, float64(damageReductionFactor(25)), 1e-6)
	assert.InDelta(0.769230, float64(damageReductionFactor(30)), 1e-6)
	assert.InDelta(0.714285, float64(damageReductionFactor(40)), 1e-6)
	assert.InDelta(0.666666, float64(damageReductionFactor(50)), 1e-6)
	assert.InDelta(0.500000, float64(damageReductionFactor(100)), 1e-6)
	assert.InDelta(0.333333, float64(damageReductionFactor(200)), 1e-6)
}
