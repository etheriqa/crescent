package main

type Event uint8

const (
	_ Event = iota
	EventDead
	EventDisableInterrupt
	EventGameTick
	EventResourceDecreased
	EventTicker
)

type EventDispatcher struct {
	handlers map[Event]map[EventHandler]bool
}

type EventHandler interface {
	HandleEvent(Event)
}

// NewEventDispatcher returns a EventDispatcher
func NewEventDispatcher() *EventDispatcher {
	return &EventDispatcher{
		handlers: make(map[Event]map[EventHandler]bool),
	}
}

// AddEventHandler adds the EventHandler if not exists
func (d *EventDispatcher) AddEventHandler(h EventHandler, e Event) {
	if d.handlers[e] == nil {
		d.handlers[e] = make(map[EventHandler]bool)
	}
	d.handlers[e][h] = true
}

// RemoveEventHandler removes the EventHandler if exists
func (d *EventDispatcher) RemoveEventHandler(h EventHandler, e Event) {
	if d.handlers[e] == nil {
		return
	}
	delete(d.handlers[e], h)
}

// TriggerEvent triggers the Event
func (d *EventDispatcher) TriggerEvent(e Event) {
	for h := range d.handlers[e] {
		h.HandleEvent(e)
	}
}
