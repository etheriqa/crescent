package main

type Event uint8

const (
	_ Event = iota

	EventGameTick
	EventPeriodicalTick

	EventDead
	EventDisabled
	EventTakenDamage
)

type EventHandler interface {
	HandleEvent(Event)
}

type EventDispatcher interface {
	AddEventHandler(EventHandler, Event)
	RemoveEventHandler(EventHandler, Event)
	TriggerEvent(Event)
}

type EventHandlerSet map[Event]map[EventHandler]bool

// NewEventDispatcher returns a EventDispatcher
func MakeEventHandlerSet() EventHandlerSet {
	return make(map[Event]map[EventHandler]bool)
}

// AddEventHandler adds the EventHandler if not exists
func (hs EventHandlerSet) AddEventHandler(h EventHandler, e Event) {
	if hs[e] == nil {
		hs[e] = make(map[EventHandler]bool)
	}
	hs[e][h] = true
}

// RemoveEventHandler removes the EventHandler if exists
func (hs EventHandlerSet) RemoveEventHandler(h EventHandler, e Event) {
	if hs[e] == nil {
		return
	}
	delete(hs[e], h)
}

// TriggerEvent triggers the Event
func (hs EventHandlerSet) TriggerEvent(e Event) {
	for h := range hs[e] {
		h.HandleEvent(e)
	}
}
