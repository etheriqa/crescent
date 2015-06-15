package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnit(t *testing.T) {
	assert := assert.New(t)

	class := &Class{
		Name:                 "Healer",
		Health:               1000,
		HealthRegeneration:   5,
		Mana:                 400,
		ManaRegeneration:     2,
		Armor:                DefaultArmor,
		MagicResistance:      DefaultMagicResistance,
		CriticalStrikeChance: DefaultCriticalStrikeChance,
		CriticalStrikeFactor: DefaultCriticalStrikeFactor,
		CooldownReduction:    DefaultCooldownReduction,
		DamageThreatFactor:   DefaultDamageThreatFactor,
		HealingThreatFactor:  DefaultHealingThreatFactor,
		Abilities:            []*Ability{},
	}
	u := NewUnit(100, 0, 1, "user", class)

	assert.Implements((*Subject)(nil), u)
	assert.Implements((*Object)(nil), u)

	assert.Equal(u, u.Subject())
	assert.Equal(u, u.Object())
	assert.Equal(UnitID(100), u.ID())
	assert.Equal(UnitGroup(0), u.Group())
	assert.Equal(UnitPosition(1), u.Position())
	assert.Equal(UnitName("user"), u.Name())
	assert.Equal(ClassName("Healer"), u.ClassName())
	assert.True(u.IsAlive())
	assert.False(u.IsDead())
	assert.Equal(Statistic(1000), u.Health())
	assert.Equal(Statistic(1000), u.HealthMax())
	assert.Equal(Statistic(5), u.HealthRegeneration())
	assert.Equal(Statistic(400), u.Mana())
	assert.Equal(Statistic(400), u.ManaMax())
	assert.Equal(Statistic(2), u.ManaRegeneration())
	assert.Equal(DefaultArmor, u.Armor())
	assert.Equal(DefaultMagicResistance, u.MagicResistance())
	assert.Equal(damageReductionFactor(DefaultArmor), u.PhysicalDamageReductionFactor())
	assert.Equal(damageReductionFactor(DefaultMagicResistance), u.MagicDamageReductionFactor())
	assert.Equal(DefaultCriticalStrikeChance, u.CriticalStrikeChance())
	assert.Equal(DefaultCriticalStrikeFactor, u.CriticalStrikeFactor())
	assert.Equal(DefaultCooldownReduction, u.CooldownReduction())
	assert.Equal(DefaultDamageThreatFactor, u.DamageThreatFactor())
	assert.Equal(DefaultHealingThreatFactor, u.HealingThreatFactor())
	assert.Nil(u.Ability("Q"))

	u.UpdateCorrection(UnitCorrection{
		Armor:                75,
		MagicResistance:      25,
		CriticalStrikeChance: 0.05,
		CriticalStrikeFactor: 0.5,
		CooldownReduction:    0.1,
		DamageThreatFactor:   2,
		HealingThreatFactor:  1,
	})
	assert.Equal(75+DefaultArmor, u.Armor())
	assert.Equal(25+DefaultMagicResistance, u.MagicResistance())
	assert.Equal(damageReductionFactor(75+DefaultArmor), u.PhysicalDamageReductionFactor())
	assert.Equal(damageReductionFactor(25+DefaultMagicResistance), u.MagicDamageReductionFactor())
	assert.Equal(0.05+DefaultCriticalStrikeChance, u.CriticalStrikeChance())
	assert.Equal(0.5+DefaultCriticalStrikeFactor, u.CriticalStrikeFactor())
	assert.Equal(0.1+DefaultCooldownReduction, u.CooldownReduction())
	assert.Equal(2+DefaultDamageThreatFactor, u.DamageThreatFactor())
	assert.Equal(1+DefaultHealingThreatFactor, u.HealingThreatFactor())

	w := new(MockedInstanceOutputWriter)

	{
		before, after, err := u.ModifyMana(w, 800)
		if assert.Nil(err) {
			assert.Equal(Statistic(400), before)
			assert.Equal(Statistic(400), after)
		}
		assert.Equal(Statistic(400), u.Mana())
		assert.Equal(Statistic(400), u.ManaMax())
	}
	{
		before, after, err := u.ModifyHealth(w, 2000)
		if assert.Nil(err) {
			assert.Equal(Statistic(1000), before)
			assert.Equal(Statistic(1000), after)
		}
		assert.Equal(Statistic(1000), u.Health())
		assert.Equal(Statistic(1000), u.HealthMax())
	}

	{
		w.On("Write", OutputUnitResource{
			UnitID: 100,
			Health: 1000,
			Mana:   200,
		}).Return()
		before, after, err := u.ModifyMana(w, -200)
		if assert.Nil(err) {
			assert.Equal(Statistic(400), before)
			assert.Equal(Statistic(200), after)
		}
		assert.Equal(Statistic(200), u.Mana())
		assert.Equal(Statistic(400), u.ManaMax())
		w.AssertExpectations(t)
	}

	{
		w.On("Write", OutputUnitResource{
			UnitID: 100,
			Health: 500,
			Mana:   200,
		}).Return()
		before, after, err := u.ModifyHealth(w, -500)
		if assert.Nil(err) {
			assert.Equal(Statistic(1000), before)
			assert.Equal(Statistic(500), after)
		}
		assert.Equal(Statistic(500), u.Health())
		assert.Equal(Statistic(1000), u.HealthMax())
		w.AssertExpectations(t)
	}

	{
		w.On("Write", OutputUnitResource{
			UnitID: 100,
			Health: 500,
			Mana:   0,
		}).Return()
		before, after, err := u.ModifyMana(w, -400)
		if assert.Nil(err) {
			assert.Equal(Statistic(200), before)
			assert.Equal(Statistic(0), after)
		}
		assert.Equal(Statistic(0), u.Mana())
		assert.Equal(Statistic(400), u.ManaMax())
		w.AssertExpectations(t)
	}

	{
		w.On("Write", OutputUnitResource{
			UnitID: 100,
			Health: 0,
			Mana:   0,
		}).Return()
		before, after, err := u.ModifyHealth(w, -2000)
		if assert.Nil(err) {
			assert.Equal(Statistic(500), before)
			assert.Equal(Statistic(0), after)
		}
		assert.Equal(Statistic(0), u.Health())
		assert.Equal(Statistic(1000), u.HealthMax())
		w.AssertExpectations(t)
	}

	assert.False(u.IsAlive())
	assert.True(u.IsDead())

	{
		_, _, err := u.ModifyMana(w, 400)
		assert.NotNil(err)
	}

	{
		_, _, err := u.ModifyHealth(w, 2000)
		assert.NotNil(err)
	}
}
