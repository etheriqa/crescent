package main

type operator interface {
	onAttach(u *unit)
	onDetach(u *unit)
}
