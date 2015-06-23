package crescent

type Effect interface{}

type EffectWillAttach interface {
	EffectWillAttach(Game) error
}

type EffectDidAttach interface {
	EffectDidAttach(Game) error
}

type EffectWillDetach interface {
	EffectWillDetach(Game) error
}

type EffectDidDetach interface {
	EffectDidDetach(Game) error
}
