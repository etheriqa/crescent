package main

type hotType string

type hot struct {
	partialOperator
	healing *healing
}

// newHoT returns a HoT
func newHoT(h *healing, duration gameDuration) *hot {
	return &hot{
		partialOperator: partialOperator{
			unit:           h.receiver,
			performer:      h.performer,
			expirationTime: h.receiver.after(duration),
		},
		healing: h,
	}
}

// onAttach removes duplicate HoTs
func (h *hot) onAttach() {
	h.addEventHandler(h, eventDead)
	h.addEventHandler(h, eventGameTick)
	h.addEventHandler(h, eventXoT)
	for o := range h.operators {
		switch o := o.(type) {
		case *hot:
			if o == h || o.performer != h.performer || o.healing.name != h.healing.name {
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

// onDetach removes the eventHandlers
func (h *hot) onDetach() {
	h.removeEventHandler(h, eventDead)
	h.removeEventHandler(h, eventGameTick)
	h.removeEventHandler(h, eventXoT)
}

// handleEvent handles the event
func (h *hot) handleEvent(e event) {
	switch e {
	case eventDead:
		h.detachOperator(h)
	case eventGameTick:
		h.expire(h, message{
			// TODO pack message
			t: outHoTEnd,
		})
	case eventXoT:
		h.healing.perform(h.game)
	}
}
