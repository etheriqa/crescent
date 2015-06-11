package main

import (
	"testing"

	"github.com/stretchr/testify/mock"
)

type MockedEventHandler struct {
	mock.Mock
}

func (m *MockedEventHandler) HandleEvent(e Event) {
	m.Called(e)
}

func TestEventDispatcher(t *testing.T) {
	dispatcher := NewEventDispatcher()

	a := new(MockedEventHandler)
	a.On("HandleEvent", EventGameTick).Return().Times(3)
	dispatcher.RemoveEventHandler(a, EventGameTick)
	dispatcher.TriggerEvent(EventGameTick)
	dispatcher.AddEventHandler(a, EventGameTick)
	dispatcher.AddEventHandler(a, EventGameTick)
	dispatcher.TriggerEvent(EventGameTick)
	dispatcher.TriggerEvent(EventGameTick)
	dispatcher.TriggerEvent(EventGameTick)
	dispatcher.RemoveEventHandler(a, EventGameTick)
	dispatcher.RemoveEventHandler(a, EventGameTick)
	dispatcher.TriggerEvent(EventGameTick)
	a.AssertExpectations(t)

	b := new(MockedEventHandler)
	b.On("HandleEvent", EventGameTick).Return().Times(3)
	b.On("HandleEvent", EventXoT).Return().Once()
	dispatcher.AddEventHandler(b, EventGameTick)
	dispatcher.TriggerEvent(EventGameTick)
	dispatcher.TriggerEvent(EventXoT)
	dispatcher.AddEventHandler(b, EventXoT)
	dispatcher.TriggerEvent(EventGameTick)
	dispatcher.TriggerEvent(EventXoT)
	dispatcher.RemoveEventHandler(b, EventXoT)
	dispatcher.TriggerEvent(EventGameTick)
	dispatcher.TriggerEvent(EventXoT)
	dispatcher.RemoveEventHandler(b, EventGameTick)
	b.AssertExpectations(t)

	c := new(MockedEventHandler)
	d := new(MockedEventHandler)
	c.On("HandleEvent", EventGameTick).Return().Once()
	d.On("HandleEvent", EventGameTick).Return().Once()
	dispatcher.AddEventHandler(c, EventGameTick)
	dispatcher.AddEventHandler(d, EventGameTick)
	dispatcher.TriggerEvent(EventGameTick)
	dispatcher.RemoveEventHandler(c, EventGameTick)
	dispatcher.RemoveEventHandler(d, EventGameTick)
	c.AssertExpectations(t)
	d.AssertExpectations(t)
}
