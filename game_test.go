package main

import (
	"github.com/stretchr/testify/mock"
)

type MockedGame struct{ mock.Mock }

func (m *MockedGame) Clock() InstanceClock {
	args := m.Called()
	return args.Get(0).(InstanceClock)
}

func (m *MockedGame) Effects() EffectContainer {
	args := m.Called()
	return args.Get(0).(EffectContainer)
}

func (m *MockedGame) Units() UnitContainer {
	args := m.Called()
	return args.Get(0).(UnitContainer)
}

func (m *MockedGame) Writer() InstanceOutputWriter {
	args := m.Called()
	return args.Get(0).(InstanceOutputWriter)
}

func (m *MockedGame) Activating(s Subject, u *Unit, a *Ability) {
	m.Called(s, u, a)
}
func (m *MockedGame) Cooldown(o Object, a *Ability) {
	m.Called(o, a)
}
func (m *MockedGame) Correction(o Object, c UnitCorrection, l Statistic, d InstanceDuration, name string) {
	m.Called(o, c, l, d, name)
}
func (m *MockedGame) Disable(o Object, t DisableType, d InstanceDuration) {
	m.Called(o, t, d)
}
func (m *MockedGame) DamageThreat(s Subject, o Object, d Statistic) {
	m.Called(s, o, d)
}
func (m *MockedGame) HealingThreat(s Subject, o Object, h Statistic) {
	m.Called(s, o, h)
}
func (m *MockedGame) DoT(damage *Damage, d InstanceDuration, name string) {
	m.Called(damage, d, name)
}
func (m *MockedGame) HoT(healing *Healing, d InstanceDuration, name string) {
	m.Called(healing, d, name)
}
func (m *MockedGame) PhysicalDamage(s Subject, o Object, d Statistic) *Damage {
	args := m.Called(s, o, d)
	return args.Get(0).(*Damage)
}
func (m *MockedGame) MagicDamage(s Subject, o Object, d Statistic) *Damage {
	args := m.Called(s, o, d)
	return args.Get(0).(*Damage)
}
func (m *MockedGame) TrueDamage(s Subject, o Object, d Statistic) *Damage {
	args := m.Called(s, o, d)
	return args.Get(0).(*Damage)
}
func (m *MockedGame) PureDamage(s Subject, o Object, d Statistic) *Damage {
	args := m.Called(s, o, d)
	return args.Get(0).(*Damage)
}
func (m *MockedGame) Healing(s Subject, o Object, h Statistic) *Healing {
	args := m.Called(s, o, h)
	return args.Get(0).(*Healing)
}
