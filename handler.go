package main

type Handler interface {
	OnAttach()
	OnDetach()
}
