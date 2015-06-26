package crescent

import (
	"math/rand"

	"github.com/stretchr/testify/mock"
)

type MockedGame struct{ mock.Mock }

func (m *MockedGame) Rand() *rand.Rand {
	args := m.Called()
	return args.Get(0).(*rand.Rand)
}

func (m *MockedGame) Clock() InstanceClock {
	args := m.Called()
	return args.Get(0).(InstanceClock)
}

func (m *MockedGame) Writer() InstanceOutputWriter {
	args := m.Called()
	return args.Get(0).(InstanceOutputWriter)
}

func (m *MockedGame) Join(g UnitGroup, name UnitName, c *Class) (UnitID, error) {
	args := m.Called(g, name, c)
	return args.Get(0).(UnitID), args.Error(1)
}

func (m *MockedGame) Leave(id UnitID) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockedGame) UnitQuery() UnitQueryable {
	args := m.Called()
	return args.Get(0).(UnitQueryable)
}

func (m *MockedGame) AttachEffect(e Effect) error {
	args := m.Called(e)
	return args.Error(0)
}

func (m *MockedGame) DetachEffect(e Effect) error {
	args := m.Called(e)
	return args.Error(0)
}

func (m *MockedGame) EffectQuery() EffectQueryable {
	args := m.Called()
	return args.Get(0).(EffectQueryable)
}

func (m *MockedGame) Activating(s Subject, u *Unit, a *Ability) {
	m.Called(s, u, a)
}

func (m *MockedGame) Cooldown(o Object, a *Ability) {
	m.Called(o, a)
}

func (m *MockedGame) ResetCooldown(o Object, a *Ability) {
	m.Called(o, a)
}

func (m *MockedGame) Correction(o Object, c UnitCorrection, name string, l Statistic, d InstanceDuration) {
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

func (m *MockedGame) DoT(damage Damage, name string, d InstanceDuration) {
	m.Called(damage, name, d)
}

func (m *MockedGame) HoT(healing Healing, name string, d InstanceDuration) {
	m.Called(healing, name, d)
}
