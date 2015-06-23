package main

import (
	"errors"
)

type Game struct {
	clock InstanceClock

	stage    Stage
	handlers HandlerContainer
	units    UnitContainer

	w InstanceOutputWriter
}

// NewGame returns a Game
func NewGame(clock InstanceClock, stage Stage, w InstanceOutputWriter) *Game {
	g := &Game{
		clock: clock,

		stage:    stage,
		handlers: MakeHandlerSet(),
		units:    NewUnitMap(),

		w: w,
	}
	stage.Initialize(g)
	return g
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

// SyncUnit sends the unit information
func (g *Game) SyncUnit(w InstanceOutputWriter, id UnitID) {
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

// PerformGameTick performs the game tick routine
func (g *Game) PerformGameTick() {
	g.stage.OnTick(g)
	g.units.Each(func(u *Unit) {
		if u.IsDead() {
			return
		}
		u.Dispatch(EventGameTick{})
	})
}

// PerformPeriodicalTick performs the periodical rick routine
func (g *Game) PerformPeriodicalTick() {
	g.units.Each(func(u *Unit) {
		if u.IsDead() {
			return
		}
		u.Dispatch(EventPeriodicalTick{})
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
