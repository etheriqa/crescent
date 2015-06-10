package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnitUpdateModification(t *testing.T) {
	assert := assert.New(t)
	g := mockGame()
	c := &class{}
	u := newUnit(g, c)
	um1 := &modifier{
		partialOperator: partialOperator{
			unit: u,
		},
		unitModification: unitModification{
			armor:                10,
			magicResistance:      20,
			criticalStrikeChance: 30,
			criticalStrikeFactor: 40,
			cooldownReduction:    50,
			damageThreatFactor:   60,
			healingThreatFactor:  70,
		},
		ability: &ability{},
	}
	um2 := &modifier{
		partialOperator: partialOperator{
			unit: u,
		},
		unitModification: unitModification{
			armor:                1000,
			magicResistance:      2000,
			criticalStrikeChance: 3000,
			criticalStrikeFactor: 4000,
			cooldownReduction:    5000,
			damageThreatFactor:   6000,
			healingThreatFactor:  7000,
		},
		ability: &ability{},
	}
	assert.EqualValues(0, u.armor())
	assert.EqualValues(0, u.magicResistance())
	assert.EqualValues(0, u.criticalStrikeChance())
	assert.EqualValues(0, u.criticalStrikeFactor())
	assert.EqualValues(0, u.cooldownReduction())
	assert.EqualValues(0, u.damageThreatFactor())
	assert.EqualValues(0, u.healingThreatFactor())
	u.attachOperator(um1)
	assert.EqualValues(10, u.armor())
	assert.EqualValues(20, u.magicResistance())
	assert.EqualValues(30, u.criticalStrikeChance())
	assert.EqualValues(40, u.criticalStrikeFactor())
	assert.EqualValues(50, u.cooldownReduction())
	assert.EqualValues(60, u.damageThreatFactor())
	assert.EqualValues(70, u.healingThreatFactor())
	u.attachOperator(um2)
	assert.EqualValues(1010, u.armor())
	assert.EqualValues(2020, u.magicResistance())
	assert.EqualValues(3030, u.criticalStrikeChance())
	assert.EqualValues(4040, u.criticalStrikeFactor())
	assert.EqualValues(5050, u.cooldownReduction())
	assert.EqualValues(6060, u.damageThreatFactor())
	assert.EqualValues(7070, u.healingThreatFactor())
	u.detachOperator(um1)
	assert.EqualValues(1000, u.armor())
	assert.EqualValues(2000, u.magicResistance())
	assert.EqualValues(3000, u.criticalStrikeChance())
	assert.EqualValues(4000, u.criticalStrikeFactor())
	assert.EqualValues(5000, u.cooldownReduction())
	assert.EqualValues(6000, u.damageThreatFactor())
	assert.EqualValues(7000, u.healingThreatFactor())
}
