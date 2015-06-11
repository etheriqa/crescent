package main

type Threat struct {
	*PartialHandler
	threat statistic
}

// NewThreat returns a Threat handler
func NewThreat(subject, object *unit, t statistic) *Threat {
	return &Threat{
		PartialHandler: NewPermanentPartialHandler(subject, object),
		threat:         t,
	}
}

// newDamageThreat initializes a threat handler with damage
func newDamageThreat(subject, object *unit, d statistic) *Threat {
	return NewThreat(subject, object, d*object.damageThreatFactor())
}

// newHealingThreat initializes a threat handler with healing
func newHealingThreat(subject, object *unit, h statistic) *Threat {
	return NewThreat(subject, object, h*object.healingThreatFactor())
}

// OnAttach merges threat handlers they have same subject
func (t *Threat) OnAttach() {
	t.Subject().AddEventHandler(t, EventDead)
	t.Object().AddEventHandler(t, EventDead)
	t.Container().ForSubjectHandler(t.Subject(), func(ha Handler) {
		switch ha := ha.(type) {
		case *Threat:
			if ha == t || ha.Object() != t.Object() {
				return
			}
			t.threat += ha.threat
			ha.Stop(ha)
		}
	})
}

// OnDetach removes the EventHandler
func (t *Threat) OnDetach() {
	t.Subject().RemoveEventHandler(t, EventDead)
	t.Object().RemoveEventHandler(t, EventDead)
}

// HandleEvent handles the event
func (t *Threat) HandleEvent(e Event) {
	switch e {
	case EventDead:
		t.Stop(t)
	}
}
