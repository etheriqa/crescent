package main

type hotType string

type HoT struct {
	partialHandler
	healing *healing
	ability *ability
}

// NewHoT returns a HoT
func NewHoT(h *healing, a *ability, duration gameDuration) *HoT {
	return &HoT{
		partialHandler: partialHandler{
			unit:           h.receiver,
			performer:      h.performer,
			expirationTime: h.receiver.after(duration),
		},
		healing: h,
		ability: a,
	}
}

// OnAttach removes duplicate HoTs
func (h *HoT) OnAttach() {
	h.AddEventHandler(h, EventDead)
	h.AddEventHandler(h, EventGameTick)
	h.AddEventHandler(h, EventXoT)
	for ha := range h.handlers {
		switch ha := ha.(type) {
		case *HoT:
			if ha == h || ha.performer != h.performer || ha.ability != h.ability {
				continue
			}
			if ha.expirationTime > h.expirationTime {
				h.detachHandler(h)
				return
			}
			h.detachHandler(ha)
		}
	}
	h.publish(message{
		// TODO pack message
		t: outHoTBegin,
	})
}

// OnDetach removes the EventHandlers
func (h *HoT) OnDetach() {
	h.RemoveEventHandler(h, EventDead)
	h.RemoveEventHandler(h, EventGameTick)
	h.RemoveEventHandler(h, EventXoT)
}

// HandleEvent handles the event
func (h *HoT) HandleEvent(e Event) {
	switch e {
	case EventDead:
		h.detachHandler(h)
	case EventGameTick:
		h.expire(h, message{
			// TODO pack message
			t: outHoTEnd,
		})
	case EventXoT:
		h.healing.perform(h.game)
	}
}
