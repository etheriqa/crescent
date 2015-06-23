package main

import (
	"github.com/stretchr/testify/mock"
)

type MockedEffect struct{ mock.Mock }
type MockedEffectS struct{ MockedEffect }
type MockedEffectO struct{ MockedEffect }
type MockedEffectSO struct{ MockedEffect }

func (m *MockedEffect) EffectDidAttach() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockedEffect) EffectDidDetach() error {
	args := m.Called()
	return args.Error(0)
}

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
