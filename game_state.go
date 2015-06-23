package main

import (
	"errors"
)

type GameState struct {
	clock InstanceClock

	stage   Stage
	effects EffectSet
	units   *UnitMap

	w InstanceOutputWriter
}

// NewGameState returns a GameState
func NewGameState(clock InstanceClock, stage Stage, w InstanceOutputWriter) *GameState {
	g := &GameState{
		clock: clock,

		stage:   stage,
		effects: MakeEffectSet(),
		units:   NewUnitMap(),

		w: w,
	}
	stage.Initialize(g)
	return g
}

// SyncGameState sends the game state
func (g *GameState) SyncGameState(w InstanceOutputWriter) {
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

// SyncUnit sends the unit information
func (g *GameState) SyncUnit(w InstanceOutputWriter, id UnitID) {
	// TODO refactor
	u := g.units.Find(id)
	// TODO handle the error
	if u == nil {
		return
	}
	as := make([]OutputPlayerAbility, 4)
	for i := 0; i < 4; i++ {
		dts := make([]string, 0)
		for _, dt := range u.Abilities()[i].DisableTypes {
			switch dt {
			case DisableTypeSilence:
				dts = append(dts, "Silence")
			case DisableTypeStun:
				dts = append(dts, "Stun")
			}
		}
		as[i] = OutputPlayerAbility{
			Name:               u.Abilities()[i].Name,
			Description:        u.Abilities()[i].Description,
			TargetType:         u.Abilities()[i].TargetType,
			HealthCost:         u.Abilities()[i].HealthCost,
			ManaCost:           u.Abilities()[i].ManaCost,
			ActivationDuration: u.Abilities()[i].ActivationDuration,
			CooldownDuration:   u.Abilities()[i].CooldownDuration,
			DisableTypes:       dts,
		}
	}
	w.Write(OutputPlayer{
		UnitID: id,
		Q:      as[0],
		W:      as[1],
		E:      as[2],
		R:      as[3],
	})
}

// Ability activates the ability
func (g *GameState) Ability(sid UnitID, oid *UnitID, abilityName string) error {
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
	// TODO refactor
	switch abilityName {
	case "Q":
		g.Activating(s, o, s.Abilities()[0])
	case "W":
		g.Activating(s, o, s.Abilities()[1])
	case "E":
		g.Activating(s, o, s.Abilities()[2])
	case "R":
		g.Activating(s, o, s.Abilities()[3])
	default:
		return errors.New("Unknown ability name")
	}
	return nil
}

// Interrupt interrupts ability activation
func (g *GameState) Interrupt(id UnitID) error {
	u := g.units.Find(id)
	if u == nil {
		return errors.New("Unknown UnitID")
	}
	u.Dispatch(EventInterrupt{
		UnitID: id,
	})
	return nil
}

// PerformGameTick performs the game tick routine
func (g *GameState) PerformGameTick() {
	g.stage.OnTick(g)
	g.units.Each(func(u *Unit) {
		if u.IsDead() {
			return
		}
		u.Dispatch(EventGameTick{})
	})
}

// PerformPeriodicalTick performs the periodical rick routine
func (g *GameState) PerformPeriodicalTick() {
	g.units.Each(func(u *Unit) {
		if u.IsDead() {
			return
		}
		u.Dispatch(EventPeriodicalTick{})
	})
}

// PerformRegenerationTick performs the regeneration tick routine
func (g *GameState) PerformRegenerationTick() {
	g.units.Each(func(u *Unit) {
		if u.IsDead() {
			return
		}
		u.ModifyHealth(g.w, u.HealthRegeneration())
		u.ModifyMana(g.w, u.ManaRegeneration())
	})
}
