package main

import (
	"time"
)

type modifier struct {
	partialOperator
	unitModification
}

// newModifier initalizes a modifier
func newModifier(u *unit, duration time.Duration, m unitModification) *modifier {
	return &modifier{
		partialOperator: partialOperator{
			unit:           u,
			performer:      nil,
			expirationTime: u.now() + gameTime(duration/gameTick),
		},
		unitModification: m,
	}
}

// onAttach updates the modificationStats of the unit
func (m *modifier) onAttach() {
	// TODO stack
	// TODO remove duplicates
	m.addEventHandler(m, eventDead)
	m.addEventHandler(m, eventGameTick)
	m.updateModification()
	m.publish(message{
		// TODO pack message
		t: outModifierBegin,
	})
}

// onDetach updates the modificationStats of the unit
func (m *modifier) onDetach() {
	m.removeEventHandler(m, eventDead)
	m.removeEventHandler(m, eventGameTick)
	m.updateModification()
}

// handleEvent handles the event
func (m *modifier) handleEvent(e event) {
	switch e {
	case eventDead:
		m.detachOperator(m)
	case eventGameTick:
		m.expire(m, message{
			// TODO pack message
			t: outModifierEnd,
		})
	}
}
