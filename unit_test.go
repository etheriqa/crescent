package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnitUpdateModification(t *testing.T) {
	assert := assert.New(t)
	g := mockGame()
	u := newUnit(g)
	um1 := &modifier{
		partialOperator: partialOperator{
			unit: u,
		},
		um: &unitModification{
			armor:                10,
			magicResistance:      20,
			criticalStrikeChance: 30,
			criticalStrikeDamage: 40,
			cooldownReduction:    50,
			threatFactor:         60,
		},
	}
	um2 := &modifier{
		partialOperator: partialOperator{
			unit: u,
		},
		um: &unitModification{
			armor:                1000,
			magicResistance:      2000,
			criticalStrikeChance: 3000,
			criticalStrikeDamage: 4000,
			cooldownReduction:    5000,
			threatFactor:         6000,
		},
	}
	assert.EqualValues(0, u.armor())
	assert.EqualValues(0, u.magicResistance())
	assert.EqualValues(0, u.criticalStrikeChance())
	assert.EqualValues(0, u.criticalStrikeDamage())
	assert.EqualValues(0, u.cooldownReduction())
	assert.EqualValues(0, u.threatFactor())
	u.attachOperator(um1)
	assert.EqualValues(10, u.armor())
	assert.EqualValues(20, u.magicResistance())
	assert.EqualValues(30, u.criticalStrikeChance())
	assert.EqualValues(40, u.criticalStrikeDamage())
	assert.EqualValues(50, u.cooldownReduction())
	assert.EqualValues(60, u.threatFactor())
	u.attachOperator(um2)
	assert.EqualValues(1010, u.armor())
	assert.EqualValues(2020, u.magicResistance())
	assert.EqualValues(3030, u.criticalStrikeChance())
	assert.EqualValues(4040, u.criticalStrikeDamage())
	assert.EqualValues(5050, u.cooldownReduction())
	assert.EqualValues(6060, u.threatFactor())
	u.detachOperator(um1)
	assert.EqualValues(1000, u.armor())
	assert.EqualValues(2000, u.magicResistance())
	assert.EqualValues(3000, u.criticalStrikeChance())
	assert.EqualValues(4000, u.criticalStrikeDamage())
	assert.EqualValues(5000, u.cooldownReduction())
	assert.EqualValues(6000, u.threatFactor())
}
