package crescent

type StageID uint64

type Stage interface {
	Initialize(Game) error
	OnTick(Game)
}
