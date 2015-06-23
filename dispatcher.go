package main

type EventHandler interface {
	Handle(interface{})
}

type EventDispatcher interface {
	Register(EventHandler)
	Unregister(EventHandler)
	Dispatch(interface{})
}

type EventHandlerSet map[EventHandler]bool

// NewEventHandlerSet returns a EventHandlerSet
func MakeEventHandlerSet() EventHandlerSet {
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
