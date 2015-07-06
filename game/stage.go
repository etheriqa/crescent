package game

type StageID uint64

type Stage interface {
	Initialize(Game) error
	OnTick(Game)
}

type StageFactory interface {
	New(StageID) Stage
}

type StageFactories map[StageID](func() Stage)

// New creates a Stage
func (sf StageFactories) New(id StageID) Stage {
	if f, ok := sf[id]; !ok {
		return nil
	} else {
		return f()
	}
}
