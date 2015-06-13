package main

type Game struct {
	clock    GameClock
	handlers HandlerContainer
	units    UnitContainer
	writer   GameEventWriter
}

// PerformRegeneration performs health regeneration and mana regeneration
func (g *Game) PerformRegeneration() {
	g.units.Each(func(u *Unit) {
		if u.IsDead() {
			return
		}
		u.ModifyHealth(g.writer, u.HealthRegeneration())
		u.ModifyMana(g.writer, u.ManaRegeneration())
	})
}
