package main

import (
	"errors"
)

type DisableType uint8

const (
	_ DisableType = iota

	DisableTypeSilence
	DisableTypeStun
)

func (d DisableType) String() string {
	switch d {
	case DisableTypeSilence:
		return "Silence"
	case DisableTypeStun:
		return "Stun"
	default:
		return ""
	}
}

type Disable struct {
	UnitObject
	disableType    DisableType
	expirationTime InstanceTime

	handler EventHandler
}

// NewDisable returns a Disable
func NewDisable(g Game, o Object, d DisableType, t InstanceTime) *Disable {
	e := &Disable{
		UnitObject:     MakeObject(o),
		disableType:    d,
		expirationTime: t,
		handler:        new(func(interface{})),
	}
	*e.handler = func(p interface{}) { e.handle(g, p) }
	return e
}

// EffectWillAttach removes duplicate Disables
func (e *Disable) EffectWillAttach(g Game) error {
	ok := g.EffectQuery().BindObject(e).Every(func(f Effect) bool {
		switch f := f.(type) {
		case *Disable:
			if e.disableType != f.disableType {
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
		return errors.New("Already disabled")
	}
	return nil
}

// EffectDidAttach dispatches a EventDisabled
func (e *Disable) EffectDidAttach(g Game) error {
	e.writeOutputUnitAttach(g)
	e.Object().Register(e.handler)
	e.Object().Dispatch(EventDisabled{})
	return nil
}

// EffectDidDetach does nothing
func (e *Disable) EffectDidDetach(g Game) error {
	e.Object().Unregister(e.handler)
	return nil
}

// handle handles the payload
func (e *Disable) handle(g Game, p interface{}) {
	switch p.(type) {
	case EventGameTick:
		if g.Clock().Before(e.expirationTime) {
			return
		}
		e.writeOutputUnitDetach(g)
		g.DetachEffect(e)
	case EventDead:
		g.DetachEffect(e)
	}
}

// writeOutputUnitAttach writes a OutputUnitAttach
func (e *Disable) writeOutputUnitAttach(g Game) {
	g.Writer().Write(OutputUnitAttach{
		UnitID:         e.Object().ID(),
		AttachmentName: e.disableType.String(),
		ExpirationTime: e.expirationTime,
	})
}

// writeOutputUnitDetach writes a OutputUnitDetach
func (e *Disable) writeOutputUnitDetach(g Game) {
	g.Writer().Write(OutputUnitDetach{
		UnitID:         e.Object().ID(),
		AttachmentName: e.disableType.String(),
	})
}
