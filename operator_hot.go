package main

type hotType string

type hot struct {
	partialOperator
	healing *healing
	ability *ability
}

// newHoT returns a HoT
func newHoT(h *healing, a *ability, duration gameDuration) *hot {
	return &hot{
		partialOperator: partialOperator{
			unit:           h.receiver,
			performer:      h.performer,
			expirationTime: h.receiver.after(duration),
		},
		healing: h,
		ability: a,
	}
}

// onAttach removes duplicate HoTs
func (h *hot) onAttach() {
	h.AddEventHandler(h, EventDead)
	h.AddEventHandler(h, EventGameTick)
	h.AddEventHandler(h, EventXoT)
	for o := range h.operators {
		switch o := o.(type) {
		case *hot:
			if o == h || o.performer != h.performer || o.ability != h.ability {
				continue
			}
			if o.expirationTime > h.expirationTime {
				h.detachOperator(h)
				return
			}
			h.detachOperator(o)
		}
	}
	h.publish(message{
		// TODO pack message
		t: outHoTBegin,
	})
}

// onDetach removes the EventHandlers
func (h *hot) onDetach() {
	h.RemoveEventHandler(h, EventDead)
	h.RemoveEventHandler(h, EventGameTick)
	h.RemoveEventHandler(h, EventXoT)
}

// HandleEvent handles the event
func (h *hot) HandleEvent(e Event) {
	switch e {
	case EventDead:
		h.detachOperator(h)
	case EventGameTick:
		h.expire(h, message{
			// TODO pack message
			t: outHoTEnd,
		})
	case EventXoT:
		h.healing.perform(h.game)
	}
}
