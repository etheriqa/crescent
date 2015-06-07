package main

type threat struct {
	partialOperator
	threat statistic
}

// newThreat initializes a threat operator
func newThreat(performer, receiver *unit, t statistic) *threat {
	return &threat{
		partialOperator: partialOperator{
			unit:      receiver,
			performer: performer,
		},
		threat: t,
	}
}

// newDamageThreat initializes a threat operator with damage
func newDamageThreat(performer, receiver *unit, damage statistic) *threat {
	return newThreat(performer, receiver, damage*performer.damageThreatFactor())
}

// newHealingThreat initializes a threat operator with healing
func newHealingThreat(performer, receiver *unit, healing statistic) *threat {
	return newThreat(performer, receiver, healing*performer.healingThreatFactor())
}

// onAttach merges threat operators they have same performer
func (t *threat) onAttach() {
	t.addEventHandler(t, eventDead)
	for o := range t.operators {
		if o == t {
			continue
		}
		if _, ok := o.(*threat); !ok {
			continue
		}
		if o.(*threat).performer != t.performer {
			continue
		}
		t.threat += o.(*threat).threat
		t.detachOperator(o)
	}
}

// onDetach removes the eventHandler
func (t *threat) onDetach() {
	t.removeEventHandler(t, eventDead)
}

// handleEvent handles the event
func (t *threat) handleEvent(e event) {
	switch e {
	case eventDead:
		t.detachOperator(t)
	}
}
