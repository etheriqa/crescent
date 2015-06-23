package main

type Effect interface {
	OnAttach()
	OnDetach()
}
