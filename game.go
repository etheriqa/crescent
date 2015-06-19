package main

import (
	"errors"
)

type Game struct {
	clock    InstanceClock
	handlers HandlerContainer
	units    UnitContainer

	w InstanceOutputWriter
}

// NewGame returns a Game
func NewGame(clock InstanceClock, w InstanceOutputWriter) *Game {
	return &Game{
		clock:    clock,
		handlers: MakeHandlerSet(),
		units:    NewUnitMap(),

		w: w,
	}
}

// Clear clears the game state
func (g *Game) Clear() {
	g.handlers = MakeHandlerSet()
	g.units.Clear()
}

// SyncGame sends the game state
func (g *Game) SyncGame(w InstanceOutputWriter) {
	g.units.Each(func(u *Unit) {
		w.Write(OutputUnitJoin{
			UnitID:    u.ID(),
			UnitGroup: u.Group(),
			UnitName:  u.Name(),
			ClassName: u.ClassName(),
			Health:    u.Health(),
			HealthMax: u.HealthMax(),
			Mana:      u.Mana(),
			ManaMax:   u.ManaMax(),
		})
	})
}

// Join creates a Unit and adds it to the game
func (g *Game) Join(group UnitGroup, name UnitName, class *Class) (id UnitID, err error) {
	u, err := g.units.Join(group, name, class)
	if err != nil {
		return
	}
	id = u.ID()
	g.w.Write(OutputUnitJoin{
		UnitID:    u.ID(),
		UnitGroup: u.Group(),
		UnitName:  u.Name(),
		ClassName: u.ClassName(),
		Health:    u.Health(),
		HealthMax: u.HealthMax(),
		Mana:      u.Mana(),
		ManaMax:   u.ManaMax(),
	})
	return
}

// Leave removes the Unit
func (g *Game) Leave(id UnitID) (err error) {
	if err = g.units.Leave(id); err != nil {
		return
	}
	g.w.Write(OutputUnitLeave{
		UnitID: id,
	})
	return
}

// Ability activates the ability
func (g *Game) Ability(sid UnitID, oid *UnitID, abilityName string) error {
	s := g.units.Find(sid)
	if s == nil {
		return errors.New("Unknown subject UnitID")
	}
	var o *Unit
	if oid != nil {
		o = g.units.Find(*oid)
		if o == nil {
			return errors.New("Unknown object UnitID")
		}
	}
	a := s.Ability(abilityName)
	if a == nil {
		return errors.New("Unknown ability name")
	}
	g.Activating(s, o, a)
	return nil
}

// PerformGameTick performs the game tick routine
func (g *Game) PerformGameTick() {
	g.units.Each(func(u *Unit) {
		if u.IsDead() {
			return
		}
		u.TriggerEvent(EventGameTick)
	})
}

// PerformPeriodicalTick performs the periodical rick routine
func (g *Game) PerformPeriodicalTick() {
	g.units.Each(func(u *Unit) {
		if u.IsDead() {
			return
		}
		u.TriggerEvent(EventPeriodicalTick)
	})
}

// PerformRegenerationTick performs the regeneration tick routine
func (g *Game) PerformRegenerationTick() {
	g.units.Each(func(u *Unit) {
		if u.IsDead() {
			return
		}
		u.ModifyHealth(g.w, u.HealthRegeneration())
		u.ModifyMana(g.w, u.ManaRegeneration())
	})
}
