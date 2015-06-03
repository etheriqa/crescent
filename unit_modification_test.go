package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnitModificationAdd(t *testing.T) {
	assert := assert.New(t)
	um1 := &unitModification{
		armor:                10,
		magicResistance:      20,
		criticalStrikeChance: 30,
		criticalStrikeDamage: 40,
		cooldownReduction:    50,
		threatFactor:         60,
	}
	um2 := &unitModification{
		armor:                1000,
		magicResistance:      2000,
		criticalStrikeChance: 3000,
		criticalStrikeDamage: 4000,
		cooldownReduction:    5000,
		threatFactor:         6000,
	}
	um1.add(um2)
	assert.EqualValues(1010, um1.armor)
	assert.EqualValues(2020, um1.magicResistance)
	assert.EqualValues(3030, um1.criticalStrikeChance)
	assert.EqualValues(4040, um1.criticalStrikeDamage)
	assert.EqualValues(5050, um1.cooldownReduction)
	assert.EqualValues(6060, um1.threatFactor)
	assert.EqualValues(1000, um2.armor)
	assert.EqualValues(2000, um2.magicResistance)
	assert.EqualValues(3000, um2.criticalStrikeChance)
	assert.EqualValues(4000, um2.criticalStrikeDamage)
	assert.EqualValues(5000, um2.cooldownReduction)
	assert.EqualValues(6000, um2.threatFactor)
}
