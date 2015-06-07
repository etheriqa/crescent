package main

type threat struct {
	partialOperator
	threat int32
}

// newThreat initializes a threat operator
func newThreat(unit, performer *unit, t int32) *threat {
	return &threat{
		partialOperator: partialOperator{
			unit:      unit,
			performer: performer,
		},
		threat: t,
	}
}

// newDamageThreat initializes a threat operator with damage
func newDamageThreat(unit, performer *unit, damage int32) *threat {
	return newThreat(unit, performer, damage*performer.damageThreatFactor())
}

// newHealingThreat initializes a threat operator with healing
func newHealingThreat(unit, performer *unit, healing int32) *threat {
	return newThreat(unit, performer, healing*performer.healingThreatFactor())
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
