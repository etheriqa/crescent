package game

type LevelID uint64

type Level interface {
	Initialize(Game) error
	OnTick(Game)
}

type LevelFactory interface {
	New(LevelID) Level
}

type LevelFactories map[LevelID](func() Level)

// New creates a Level
func (sf LevelFactories) New(id LevelID) Level {
	if f, ok := sf[id]; !ok {
		return nil
	} else {
		return f()
	}
}
