package main

type hotType string

type hot struct {
	partialOperator
	hotType hotType
	healing int32
	threat  int32
}

// onAttach removes duplicate HoTs
func (h *hot) onAttach() {
	h.unit.addEventHandler(h, eventDead)
	h.unit.addEventHandler(h, eventGameTick)
	h.unit.addEventHandler(h, eventStatsTick)
	for o := range h.unit.operators {
		if o == h {
			continue
		}
		if _, ok := o.(*hot); !ok {
			continue
		}
		if o.(*hot).hotType != h.hotType {
			continue
		}
		if o.(*disable).expirationTime >= h.expirationTime {
			h.unit.detachOperator(h)
			return
		}
		h.unit.detachOperator(o)
	}
	h.unit.publish(message{
		// TODO pack message
		t: outHoTBegin,
	})
}

// onDetach removes the eventHandlers
func (h *hot) onDetach() {
	h.unit.removeEventHandler(h, eventDead)
	h.unit.removeEventHandler(h, eventGameTick)
	h.unit.removeEventHandler(h, eventStatsTick)
}

// handleEvent handles the event
func (h *hot) handleEvent(e event) {
	switch e {
	case eventDead:
		h.terminate(h)
	case eventGameTick:
		h.expire(h, message{
			// TODO pack message
			t: outHoTEnd,
		})
	case eventStatsTick:
		h.perform()
	}
}

// perform performs the HoT
func (h *hot) perform() {
	h.unit.addHealth(h.healing)
	if h.performer.isDead() {
		return
	}
	for _, enemy := range h.unit.game.enemies(h.performer) {
		if enemy.isDead() {
			return
		}
		enemy.attachOperator(newThreat(enemy, h.performer, h.threat))
	}
}
