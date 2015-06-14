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
		clock: clock,

		w: w,
	}
}

// PerformRegeneration performs health regeneration and mana regeneration
func (g *Game) PerformRegeneration() {
	g.units.Each(func(u *Unit) {
		if u.IsDead() {
			return
		}
		u.ModifyHealth(g.w, u.HealthRegeneration())
		u.ModifyMana(g.w, u.ManaRegeneration())
	})
}
