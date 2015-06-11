package main

type hotType string

type HoT struct {
	*PartialHandler
	healing *healing
	ability *ability
}

// NewHoT returns a HoT handler
func NewHoT(h *healing, a *ability, duration GameDuration) *HoT {
	return &HoT{
		PartialHandler: NewPartialHandler(h.subject, h.object, duration),
		healing:        h,
		ability:        a,
	}
}

// OnAttach removes duplicate HoTs
func (h *HoT) OnAttach() {
	h.Object().AddEventHandler(h, EventDead)
	h.Object().AddEventHandler(h, EventGameTick)
	h.Object().AddEventHandler(h, EventXoT)
	ok := h.Container().EveryObjectHandler(h.Object(), func(ha Handler) bool {
		switch ha := ha.(type) {
		case *HoT:
			if ha == h || ha.Subject() != h.Subject() || ha.ability != h.ability {
				return true
			}
			if ha.expirationTime > h.expirationTime {
				return false
			}
			ha.Stop(ha)
		}
		return true
	})
	if !ok {
		h.Stop(h)
		return
	}
	h.Publish(message{
	// TODO pack message
	})
}

// OnDetach removes the EventHandlers
func (h *HoT) OnDetach() {
	h.Object().RemoveEventHandler(h, EventDead)
	h.Object().RemoveEventHandler(h, EventGameTick)
	h.Object().RemoveEventHandler(h, EventXoT)
}

// HandleEvent handles the event
func (h *HoT) HandleEvent(e Event) {
	switch e {
	case EventDead:
		h.Stop(h)
	case EventGameTick:
		if h.IsExpired() {
			h.Up()
		}
	case EventXoT:
		h.healing.Perform()
	}
}

// Up ends the HoT
func (h *HoT) Up() {
	h.Stop(h)
	h.Publish(message{})
}
