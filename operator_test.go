package main

import (
	"github.com/stretchr/testify/mock"
)

type MockedOperator struct{ mock.Mock }

func (m *MockedOperator) Clock() InstanceClock {
	args := m.Called()
	return args.Get(0).(InstanceClock)
}

func (m *MockedOperator) Handlers() HandlerContainer {
	args := m.Called()
	return args.Get(0).(HandlerContainer)
}

func (m *MockedOperator) Units() UnitContainer {
	args := m.Called()
	return args.Get(0).(UnitContainer)
}

func (m *MockedOperator) Writer() InstanceOutputWriter {
	args := m.Called()
	return args.Get(0).(InstanceOutputWriter)
}

func (m *MockedOperator) Activating(s Subject, u *Unit, a *Ability) {
	m.Called(s, u, a)
}
func (m *MockedOperator) Cooldown(o Object, a *Ability) {
	m.Called(o, a)
}
func (m *MockedOperator) Correction(o Object, c UnitCorrection, l Statistic, d InstanceDuration, name string) {
	m.Called(o, c, l, d, name)
}
func (m *MockedOperator) Disable(o Object, t DisableType, d InstanceDuration) {
	m.Called(o, t, d)
}
func (m *MockedOperator) DamageThreat(s Subject, o Object, d Statistic) {
	m.Called(s, o, d)
}
func (m *MockedOperator) HealingThreat(s Subject, o Object, h Statistic) {
	m.Called(s, o, h)
}
func (m *MockedOperator) DoT(damage *Damage, d InstanceDuration, name string) {
	m.Called(damage, d, name)
}
func (m *MockedOperator) HoT(healing *Healing, d InstanceDuration, name string) {
	m.Called(healing, d, name)
}
func (m *MockedOperator) PhysicalDamage(s Subject, o Object, d Statistic) *Damage {
	args := m.Called(s, o, d)
	return args.Get(0).(*Damage)
}
func (m *MockedOperator) MagicDamage(s Subject, o Object, d Statistic) *Damage {
	args := m.Called(s, o, d)
	return args.Get(0).(*Damage)
}
func (m *MockedOperator) TrueDamage(s Subject, o Object, d Statistic) *Damage {
	args := m.Called(s, o, d)
	return args.Get(0).(*Damage)
}
func (m *MockedOperator) PureDamage(s Subject, o Object, d Statistic) *Damage {
	args := m.Called(s, o, d)
	return args.Get(0).(*Damage)
}
func (m *MockedOperator) Healing(s Subject, o Object, h Statistic) *Healing {
	args := m.Called(s, o, h)
	return args.Get(0).(*Healing)
}
