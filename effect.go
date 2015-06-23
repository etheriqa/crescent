package main

type Effect interface{}

type EffectWillAttach interface {
	EffectWillAttach() error
}

type EffectDidAttach interface {
	EffectDidAttach() error
}

type EffectWillDetach interface {
	EffectWillDetach() error
}

type EffectDidDetach interface {
	EffectDidDetach() error
}
