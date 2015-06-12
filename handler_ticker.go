package main

type Ticker struct {
	PartialHandler
	operator Operator
	ability  *Ability
}

// NewTicker
func NewTicker(op Operator, a *Ability, duration GameDuration) *Ticker {
	return &Ticker{
		PartialHandler: MakePartialHandler(MakeUnitPair(op.Subject(), op.Object()), duration), // TODO refactor
		operator:       op,
		ability:        a,
	}
}

// OnAttach adds the EventHandlers and removes duplicate Tickers
func (t *Ticker) OnAttach() {
	t.Object().AddEventHandler(t, EventDead)
	t.Object().AddEventHandler(t, EventGameTick)
	t.Object().AddEventHandler(t, EventTicker)
	ok := t.EveryObjectHandler(func(ha Handler) bool {
		switch ha := ha.(type) {
		case *Ticker:
			if ha == t || ha.Subject() != t.Subject() || ha.ability != t.ability {
				return true
			}
			if ha.expirationTime > t.expirationTime {
				return false
			}
			ha.Stop(ha)
		}
		return true
	})
	if !ok {
		t.Stop(t)
		return
	}
	t.Publish(message{
	// TODO pack message
	})
}

// OnDetach removes the EventHandlers
func (t *Ticker) OnDetach() {
	t.Object().RemoveEventHandler(t, EventDead)
	t.Object().RemoveEventHandler(t, EventGameTick)
	t.Object().RemoveEventHandler(t, EventTicker)
}

// HandleEvent handles the event
func (t *Ticker) HandleEvent(e Event) {
	switch e {
	case EventDead:
		t.Stop(t)
	case EventGameTick:
		if t.IsExpired() {
			t.Up()
		}
	case EventTicker:
		t.operator.Perform()
	}
}

// Up ends the Ticker
func (t *Ticker) Up() {
	t.Stop(t)
	t.Publish(message{
	// TODO pack message
	})
}
