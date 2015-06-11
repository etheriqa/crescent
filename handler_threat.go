package main

type Threat struct {
	partialHandler
	threat statistic
}

// NewThreat initializes a threat handler
func NewThreat(performer, receiver *unit, t statistic) *Threat {
	return &Threat{
		partialHandler: partialHandler{
			unit:      receiver,
			performer: performer,
		},
		threat: t,
	}
}

// newDamageThreat initializes a threat handler with damage
func newDamageThreat(performer, receiver *unit, d statistic) *Threat {
	return NewThreat(performer, receiver, d*performer.damageThreatFactor())
}

// newHealingThreat initializes a threat handler with healing
func newHealingThreat(performer, receiver *unit, h statistic) *Threat {
	return NewThreat(performer, receiver, h*performer.healingThreatFactor())
}

// OnAttach merges threat handlers they have same performer
func (t *Threat) OnAttach() {
	t.AddEventHandler(t, EventDead)
	for ha := range t.handlers {
		switch ha := ha.(type) {
		case *Threat:
			if ha == t || ha.performer != t.performer {
				continue
			}
			t.threat += ha.threat
			t.detachHandler(ha)
		}
	}
}

// OnDetach removes the EventHandler
func (t *Threat) OnDetach() {
	t.RemoveEventHandler(t, EventDead)
}

// HandleEvent handles the event
func (t *Threat) HandleEvent(e Event) {
	switch e {
	case EventDead:
		t.detachHandler(t)
	}
}
