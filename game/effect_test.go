package game

import (
	"github.com/stretchr/testify/mock"
)

type MockedEffect struct{ mock.Mock }
type MockedEffectS struct{ MockedEffect }
type MockedEffectO struct{ MockedEffect }
type MockedEffectSO struct{ MockedEffect }
type MockedFullEffect struct{ MockedEffect }

func (m *MockedEffectS) Subject() *Unit {
	args := m.Called()
	return args.Get(0).(*Unit)
}

func (m *MockedEffectO) Object() *Unit {
	args := m.Called()
	return args.Get(0).(*Unit)
}

func (m *MockedEffectSO) Subject() *Unit {
	args := m.Called()
	return args.Get(0).(*Unit)
}

func (m *MockedEffectSO) Object() *Unit {
	args := m.Called()
	return args.Get(0).(*Unit)
}

func (m *MockedFullEffect) EffectWillAttach(g Game) error {
	args := m.Called(g)
	return args.Error(0)
}

func (m *MockedFullEffect) EffectDidAttach(g Game) error {
	args := m.Called(g)
	return args.Error(0)
}

func (m *MockedFullEffect) EffectWillDetach(g Game) error {
	args := m.Called(g)
	return args.Error(0)
}

func (m *MockedFullEffect) EffectDidDetach(g Game) error {
	args := m.Called(g)
	return args.Error(0)
}
