package main

type Threat struct {
	UnitPair
	threat Statistic

	handler EventHandler
}

// NewThreat returns a Threat
func NewThreat(g Game, s Subject, o Object, t Statistic) *Threat {
	e := &Threat{
		UnitPair: MakePair(s, o),
		threat:   t,
		handler:  new(func(interface{})),
	}
	*e.handler = func(p interface{}) { e.handle(g, p) }
	return e
}

// EffectWillAttach merges Threat effects
func (e *Threat) EffectWillAttach(g Game) error {
	g.EffectQuery().BindSubject(e).BindObject(e).Each(func(f Effect) {
		switch f := f.(type) {
		case *Threat:
			e.threat += f.threat
			g.DetachEffect(f)
		}
	})
	return nil
}

// EffectDidAttach does nothing
func (e *Threat) EffectDidAttach(g Game) error {
	e.Subject().Register(e.handler)
	e.Object().Register(e.handler)
	return nil
}

// EffectDidDetach does nothing
func (e *Threat) EffectDidDetach(g Game) error {
	e.Subject().Unregister(e.handler)
	e.Object().Unregister(e.handler)
	return nil
}

// handle handles the payload
func (e *Threat) handle(g Game, p interface{}) {
	switch p := p.(type) {
	case EventDead:
		g.DetachEffect(p)
	}
}
