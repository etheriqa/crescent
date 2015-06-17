package main

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
		units:    MakeUnitMap(),

		w: w,
	}
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
