package main

import (
	"github.com/stretchr/testify/mock"
)

type MockedHandler struct{ mock.Mock }
type MockedHandlerS struct{ MockedHandler }
type MockedHandlerO struct{ MockedHandler }
type MockedHandlerSO struct{ MockedHandler }

func (m *MockedHandler) OnAttach() {
	m.Called()
}

func (m *MockedHandler) OnDetach() {
	m.Called()
}

func (m *MockedHandlerS) Subject() *Unit {
	args := m.Called()
	return args.Get(0).(*Unit)
}

func (m *MockedHandlerO) Object() *Unit {
	args := m.Called()
	return args.Get(0).(*Unit)
}

func (m *MockedHandlerSO) Subject() *Unit {
	args := m.Called()
	return args.Get(0).(*Unit)
}

func (m *MockedHandlerSO) Object() *Unit {
	args := m.Called()
	return args.Get(0).(*Unit)
}
