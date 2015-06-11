package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockedHandler struct {
	mock.Mock
	id string
}

func (m *MockedHandler) HandleEvent(Event) {}

func (m *MockedHandler) Subject() *Unit {
	args := m.Called()
	return args.Get(0).(*Unit)
}

func (m *MockedHandler) Object() *Unit {
	args := m.Called()
	return args.Get(0).(*Unit)
}

func (m *MockedHandler) OnAttach() {
	m.Called()
}

func (m *MockedHandler) OnDetach() {
	m.Called()
}

func TestHandlerContainerCallback(t *testing.T) {
	container := NewHandlerContainer()
	ha1 := new(MockedHandler)
	ha2 := new(MockedHandler)
	ha3 := new(MockedHandler)

	ha1.On("OnAttach").Return().Once()
	ha1.On("OnDetach").Return().Once()
	ha2.On("OnAttach").Return().Once()

	container.DetachHandler(ha1)
	container.DetachHandler(ha1)
	container.DetachHandler(ha2)
	container.DetachHandler(ha2)
	container.DetachHandler(ha3)
	container.DetachHandler(ha3)
	container.AttachHandler(ha1)
	container.AttachHandler(ha1)
	container.AttachHandler(ha2)
	container.AttachHandler(ha2)
	container.DetachHandler(ha1)
	container.DetachHandler(ha1)

	ha1.AssertExpectations(t)
	ha2.AssertExpectations(t)
	ha2.AssertNotCalled(t, "OnDetach")
	ha3.AssertNotCalled(t, "OnAttach")
	ha3.AssertNotCalled(t, "OnDetach")
}

func TestHandlerContainerForSubjectHandlerAndForObjectHandler(t *testing.T) {
	assert := assert.New(t)
	container := NewHandlerContainer()
	ha1 := &MockedHandler{id: "ha1"}
	ha2 := &MockedHandler{id: "ha2"}
	ha3 := &MockedHandler{id: "ha3"}
	u1 := new(Unit)
	u2 := new(Unit)
	u3 := new(Unit)

	ha1.On("OnAttach").Return()
	ha1.On("Subject").Return(u1)
	ha1.On("Object").Return(u2)
	ha2.On("OnAttach").Return()
	ha2.On("Subject").Return(u1)
	ha2.On("Object").Return(u3)
	ha3.On("OnAttach").Return()
	ha3.On("Subject").Return(u2)
	ha3.On("Object").Return(u3)

	container.AttachHandler(ha1)
	container.AttachHandler(ha2)
	container.AttachHandler(ha3)

	{
		hs := make([]Handler, 0)
		container.ForSubjectHandler(u1, func(ha Handler) {
			hs = append(hs, ha)
		})
		assert.Contains(hs, ha1)
		assert.Contains(hs, ha2)
		assert.NotContains(hs, ha3)
	}

	{
		hs := make([]Handler, 0)
		container.ForSubjectHandler(u2, func(ha Handler) {
			hs = append(hs, ha)
		})
		assert.NotContains(hs, ha1)
		assert.NotContains(hs, ha2)
		assert.Contains(hs, ha3)
	}

	{
		hs := make([]Handler, 0)
		container.ForSubjectHandler(u3, func(ha Handler) {
			hs = append(hs, ha)
		})
		assert.NotContains(hs, ha1)
		assert.NotContains(hs, ha2)
		assert.NotContains(hs, ha3)
	}

	{
		hs := make([]Handler, 0)
		container.ForObjectHandler(u1, func(ha Handler) {
			hs = append(hs, ha)
		})
		assert.NotContains(hs, ha1)
		assert.NotContains(hs, ha2)
		assert.NotContains(hs, ha3)
	}

	{
		hs := make([]Handler, 0)
		container.ForObjectHandler(u2, func(ha Handler) {
			hs = append(hs, ha)
		})
		assert.Contains(hs, ha1)
		assert.NotContains(hs, ha2)
		assert.NotContains(hs, ha3)
	}

	{
		hs := make([]Handler, 0)
		container.ForObjectHandler(u3, func(ha Handler) {
			hs = append(hs, ha)
		})
		assert.NotContains(hs, ha1)
		assert.Contains(hs, ha2)
		assert.Contains(hs, ha3)
	}

	assert.True(container.EverySubjectHandler(u1, func(Handler) bool { return true }))
	assert.True(container.EverySubjectHandler(u2, func(Handler) bool { return true }))
	assert.True(container.EverySubjectHandler(u3, func(Handler) bool { return true }))
	assert.False(container.EverySubjectHandler(u1, func(Handler) bool { return false }))
	assert.False(container.EverySubjectHandler(u2, func(Handler) bool { return false }))
	assert.True(container.EverySubjectHandler(u3, func(Handler) bool { return false }))

	assert.True(container.EveryObjectHandler(u1, func(Handler) bool { return true }))
	assert.True(container.EveryObjectHandler(u2, func(Handler) bool { return true }))
	assert.True(container.EveryObjectHandler(u3, func(Handler) bool { return true }))
	assert.True(container.EveryObjectHandler(u1, func(Handler) bool { return false }))
	assert.False(container.EveryObjectHandler(u2, func(Handler) bool { return false }))
	assert.False(container.EveryObjectHandler(u3, func(Handler) bool { return false }))

	assert.True(container.SomeSubjectHandler(u1, func(Handler) bool { return true }))
	assert.True(container.SomeSubjectHandler(u2, func(Handler) bool { return true }))
	assert.False(container.SomeSubjectHandler(u3, func(Handler) bool { return true }))
	assert.False(container.SomeSubjectHandler(u1, func(Handler) bool { return false }))
	assert.False(container.SomeSubjectHandler(u2, func(Handler) bool { return false }))
	assert.False(container.SomeSubjectHandler(u3, func(Handler) bool { return false }))

	assert.False(container.SomeObjectHandler(u1, func(Handler) bool { return true }))
	assert.True(container.SomeObjectHandler(u2, func(Handler) bool { return true }))
	assert.True(container.SomeObjectHandler(u3, func(Handler) bool { return true }))
	assert.False(container.SomeObjectHandler(u1, func(Handler) bool { return false }))
	assert.False(container.SomeObjectHandler(u2, func(Handler) bool { return false }))
	assert.False(container.SomeObjectHandler(u3, func(Handler) bool { return false }))
}
