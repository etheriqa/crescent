package main

type DisableType uint8

const (
	_ DisableType = iota

	DisableTypeSilence
	DisableTypeStun
)

type Disable struct {
	UnitObject
	disableType    DisableType
	expirationTime GameTime

	clock    GameClock
	handlers HandlerContainer
	writer   GameEventWriter
}

// OnAttach removes duplicate Disables
func (h *Disable) OnAttach() {
	ok := h.handlers.BindObject(h).Every(func(o Handler) bool {
		switch o := o.(type) {
		case *Disable:
			if h == o || h.disableType != o.disableType {
				return true
			}
			if h.expirationTime <= o.expirationTime {
				return false
			}
			h.handlers.Detach(o)
		}
		return true
	})
	if !ok {
		h.handlers.Detach(h)
		return
	}

	h.Object().AddEventHandler(h, EventGameTick)
	h.Object().AddEventHandler(h, EventDead)
	h.writer.Write(nil) // TODO
	h.Object().TriggerEvent(EventDisabled)
}

// OnDetach does nothing
func (h *Disable) OnDetach() {
	h.Object().RemoveEventHandler(h, EventGameTick)
	h.Object().RemoveEventHandler(h, EventDead)
}

// HandleEvent handles the Event
func (h *Disable) HandleEvent(e Event) {
	switch e {
	case EventGameTick:
		if h.clock.Before(h.expirationTime) {
			return
		}
		h.handlers.Detach(h)
		h.writer.Write(nil) // TODO
	case EventDead:
		h.handlers.Detach(h)
	}
}
