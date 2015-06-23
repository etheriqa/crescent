package main

type EventHandler interface {
	Handle(interface{})
}

type EventDispatcher interface {
	Register(EventHandler)
	Unregister(EventHandler)
	Dispatch(interface{})
}

type EventHandlerBared func(interface{})

type EventHandlerSet map[EventHandler]bool

// MakeEventHandler returns a EventHandlerBared
func MakeEventHandler(h func(interface{})) EventHandlerBared {
	return h
}

// Handle handles the payload
func (hb EventHandlerBared) Handle(p interface{}) {
	hb(p)
}

// MakeEventDispatcher returns a EventHandlerSet
func MakeEventDispatcher() EventHandlerSet {
	return make(map[EventHandler]bool)
}

// Register adds the EventHandler if not contains
func (hs EventHandlerSet) Register(h EventHandler) {
	hs[h] = true
}

// Unregister removes the EventHandler if contains
func (hs EventHandlerSet) Unregister(h EventHandler) {
	delete(hs, h)
}

// Dispatch calls Handle with the payload
func (hs EventHandlerSet) Dispatch(p interface{}) {
	for h := range hs {
		h.Handle(p)
	}
}
