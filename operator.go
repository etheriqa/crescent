package main

type operator interface {
	isComplete(u *unit) bool
	onAttach(u *unit)
	onTick(u *unit)
	onComplete(u *unit)
	onDetach(u *unit)
}
