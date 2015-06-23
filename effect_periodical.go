package main

import (
	"errors"
)

type Periodical struct {
	UnitPair
	name           string
	routine        func()
	expirationTime InstanceTime

	handler EventHandler
}

// NewPeriodical returns a Periodical
func NewPeriodical(g Game, s Subject, o Object, name string, routine func(), t InstanceTime) *Periodical {
	e := &Periodical{
		UnitPair:       MakePair(s, o),
		name:           name,
		routine:        routine,
		expirationTime: t,
		handler:        new(func(interface{})),
	}
	*e.handler = func(p interface{}) { e.handle(g, p) }
	return e
}

// EffectWillAttach removes duplicate Periodicals
func (e *Periodical) EffectWillAttach(g Game) error {
	ok := g.EffectQuery().BindSubject(e).BindObject(e).Every(func(f Effect) bool {
		switch f := f.(type) {
		case *Periodical:
			if e.name != f.name {
				return true
			}
			if e.expirationTime <= f.expirationTime {
				return false
			}
			g.DetachEffect(f)
		}
		return true
	})
	if !ok {
		return errors.New("Already attached")
	}
	return nil
}

// EffectDidAttach does nothing
func (e *Periodical) EffectDidAttach(g Game) error {
	e.writeOutputUnitAttach(g)
	e.Object().Register(e.handler)
	return nil
}

// EffectDidDetach does nothing
func (e *Periodical) EffectDidDetach(g Game) error {
	e.Object().Unregister(e.handler)
	return nil
}

// handle handles the Event
func (e *Periodical) handle(g Game, p interface{}) {
	switch p.(type) {
	case EventDead:
		g.DetachEffect(e)
	case EventGameTick:
		if g.Clock().Before(e.expirationTime) {
			return
		}
		e.writeOutputUnitDetach(g)
		g.DetachEffect(e)
	case EventPeriodicalTick:
		e.routine()
	}
}

// writeOutputUnitAttach writes a OutputUnitAttach
func (e *Periodical) writeOutputUnitAttach(g Game) {
	g.Writer().Write(OutputUnitAttach{
		UnitID:         e.Object().ID(),
		AttachmentName: e.name,
		ExpirationTime: e.expirationTime,
	})
}

// writeOutputUnitDetach writes a OutputUnitDetach
func (e *Periodical) writeOutputUnitDetach(g Game) {
	g.Writer().Write(OutputUnitDetach{
		UnitID:         e.Object().ID(),
		AttachmentName: e.name,
	})
}
