package main

type event uint8

const (
	_ event = iota
	eventDead
	eventDisable
	eventGameTick
	eventResourceDecreased
	eventXoT
)

type eventDispatcher struct {
	handlers map[event]map[eventHandler]bool
}

type eventHandler interface {
	handleEvent(event)
}

func newEventDispatcher() *eventDispatcher {
	return &eventDispatcher{
		handlers: make(map[event]map[eventHandler]bool),
	}
}

func (d *eventDispatcher) addEventHandler(h eventHandler, e event) {
	if d.handlers[e] == nil {
		d.handlers[e] = make(map[eventHandler]bool)
	}
	d.handlers[e][h] = true
}

func (d *eventDispatcher) removeEventHandler(h eventHandler, e event) {
	if d.handlers[e] == nil {
		d.handlers[e] = make(map[eventHandler]bool)
	}
	delete(d.handlers[e], h)
}

func (d *eventDispatcher) triggerEvent(e event) {
	for h := range d.handlers[e] {
		h.handleEvent(e)
	}
}
