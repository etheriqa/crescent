package main

type Game struct {
	clock    GameClock
	handlers HandlerContainer
	units    UnitContainer

	w InstanceOutputWriter
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
