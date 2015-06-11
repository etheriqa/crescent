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
func newDamageThreat(performer, receiver *unit, d statistic) *threat {
	return newThreat(performer, receiver, d*performer.damageThreatFactor())
}

// newHealingThreat initializes a threat operator with healing
func newHealingThreat(performer, receiver *unit, h statistic) *threat {
	return newThreat(performer, receiver, h*performer.healingThreatFactor())
}

// onAttach merges threat operators they have same performer
func (t *threat) onAttach() {
	t.AddEventHandler(t, EventDead)
	for o := range t.operators {
		switch o := o.(type) {
		case *threat:
			if o == t || o.performer != t.performer {
				continue
			}
			t.threat += o.threat
			t.detachOperator(o)
		}
	}
}

// onDetach removes the EventHandler
func (t *threat) onDetach() {
	t.RemoveEventHandler(t, EventDead)
}

// HandleEvent handles the event
func (t *threat) HandleEvent(e Event) {
	switch e {
	case EventDead:
		t.detachOperator(t)
	}
}
