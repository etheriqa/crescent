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

// PerformGameTick performs the game tick routine
func (g *Game) PerformGameTick() {
	// TODO
	log.WithField("time", g.clock.Now()).Debug("GameTick")
}

// PerformPeriodicalTick performs the periodical rick routine
func (g *Game) PerformPeriodicalTick() {
	// TODO
	log.WithField("time", g.clock.Now()).Debug("PeriodicalTick")
}

// PerformRegenerationTick performs the regeneration tick routine
func (g *Game) PerformRegenerationTick() {
	// TODO
	log.WithField("time", g.clock.Now()).Debug("RegenerationTick")
}
