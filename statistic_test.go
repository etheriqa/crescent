package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReductionFactor(t *testing.T) {
	assert := assert.New(t)
	assert.InDelta(1.000000, float64(reductionFactor(0)), 1e-6)
	assert.InDelta(0.909090, float64(reductionFactor(10)), 1e-6)
	assert.InDelta(0.833333, float64(reductionFactor(20)), 1e-6)
	assert.InDelta(0.800000, float64(reductionFactor(25)), 1e-6)
	assert.InDelta(0.769230, float64(reductionFactor(30)), 1e-6)
	assert.InDelta(0.714285, float64(reductionFactor(40)), 1e-6)
	assert.InDelta(0.666666, float64(reductionFactor(50)), 1e-6)
	assert.InDelta(0.500000, float64(reductionFactor(100)), 1e-6)
	assert.InDelta(0.333333, float64(reductionFactor(200)), 1e-6)
}
