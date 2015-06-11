package main

type Threat struct {
	PartialHandler
	threat Statistic
}

// NewThreat returns a Threat handler
func NewThreat(up UnitPair, t Statistic) *Threat {
	return &Threat{
		PartialHandler: MakePermanentPartialHandler(up),
		threat:         t,
	}
}

// NewDamageThreat initializes a threat handler with damage
func NewDamageThreat(up UnitPair, d Statistic) *Threat {
	return NewThreat(up, d*up.Object().damageThreatFactor())
}

// NewHealingThreat initializes a threat handler with healing
func NewHealingThreat(up UnitPair, h Statistic) *Threat {
	return NewThreat(up, h*up.Object().healingThreatFactor())
}

// OnAttach merges threat handlers they have same subject
func (t *Threat) OnAttach() {
	t.Subject().AddEventHandler(t, EventDead)
	t.Object().AddEventHandler(t, EventDead)
	t.ForSubjectHandler(func(ha Handler) {
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
