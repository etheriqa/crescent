package crescent

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockedEventHandler struct{ mock.Mock }

type MockedEventDispatcher struct{ mock.Mock }

func (m *MockedEventHandler) Handle(p interface{}) {
	m.Called(p)
}

func (m *MockedEventDispatcher) Register(h EventHandler) {
	m.Called(h)
}

func (m *MockedEventDispatcher) Unregister(h EventHandler) {
	m.Called(h)
}

func (m *MockedEventDispatcher) Dispatch(p interface{}) {
	m.Called(p)
}

func TestEventHandlerBared(t *testing.T) {
	assert := assert.New(t)
	nCalled := 0
	payload := (interface{})(1000)
	f := func(p interface{}) {
		nCalled++
		assert.Equal(payload, p)
	}
	h := MakeEventHandler(f)
	h.Handle(payload)
	assert.Equal(1, nCalled)
}

func TestEventHandlerSet(t *testing.T) {
	hs := MakeEventDispatcher()
	h := new(MockedEventHandler)

	hs.Unregister(h)
	hs.Dispatch("foo")

	h.On("Handle", (interface{})("bar")).Return().Once()
	hs.Register(h)
	hs.Dispatch("bar")
	h.AssertExpectations(t)

	hs.Unregister(h)
	hs.Dispatch("baz")
}
