package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDamage(t *testing.T) {
	assert := assert.New(t)
	s := NewUnit(0, 0, 0, "subject", &Class{
		Health: 1,
	})
	o := NewUnit(1, 1, 1, "object", &Class{
		Health: 1000,
	})
	op := new(MockedOperator)

	{
		damage := Damage{
			UnitPair:             MakePair(s, o),
			damage:               100,
			criticalStrikeChance: 0,
			criticalStrikeFactor: 0,
			op:                   op,
		}
		w := new(MockedInstanceOutputWriter)
		w.On("Write", OutputUnitResource{
			UnitID: 1,
			Health: 900,
			Mana:   0,
		}).Return().Once()
		w.On("Write", OutputDamage{
			SubjectUnitID: 0,
			ObjectUnitID:  1,
			Damage:        100,
			IsCritical:    false,
		}).Return().Once()
		op.On("Writer").Return(w).Twice()
		op.On("DamageThreat", &damage, &damage, Statistic(100)).Return().Once()
		before, after, crit, err := damage.Perform()
		assert.Equal(Statistic(1000), before)
		assert.Equal(Statistic(900), after)
		assert.False(crit)
		assert.Nil(err)
		w.AssertExpectations(t)
		op.AssertExpectations(t)
	}

	{
		damage := Damage{
			UnitPair:             MakePair(s, o),
			damage:               1000,
			criticalStrikeChance: 0,
			criticalStrikeFactor: 0,
			op:                   op,
		}
		w := new(MockedInstanceOutputWriter)
		w.On("Write", OutputUnitResource{
			UnitID: 1,
			Health: 0,
			Mana:   0,
		}).Return().Once()
		w.On("Write", OutputDamage{
			SubjectUnitID: 0,
			ObjectUnitID:  1,
			Damage:        1000,
			IsCritical:    false,
		}).Return().Once()
		op.On("Writer").Return(w).Twice()
		before, after, crit, err := damage.Perform()
		assert.Equal(Statistic(900), before)
		assert.Equal(Statistic(0), after)
		assert.False(crit)
		assert.Nil(err)
		w.AssertExpectations(t)
		op.AssertExpectations(t)
	}

	{
		damage := Damage{
			UnitPair:             MakePair(s, o),
			damage:               100,
			criticalStrikeChance: 0,
			criticalStrikeFactor: 0,
			op:                   op,
		}
		w := new(MockedInstanceOutputWriter)
		op.On("Writer").Return(w).Once()
		_, _, _, err := damage.Perform()
		assert.NotNil(err)
	}
}
