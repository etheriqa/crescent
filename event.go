package main

type event uint8

const (
	eventDefault event = iota
	eventDisable
	eventStats
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

func (d *eventDispatcher) addEventHandler(e event, h eventHandler) {
	d.handlers[e][h] = true
}

func (d *eventDispatcher) removeEventHandler(e event, h eventHandler) {
	delete(d.handlers[e], h)
}

func (d *eventDispatcher) triggerEvent(e event) {
	for h := range d.handlers[e] {
		h.handleEvent(e)
	}
}
