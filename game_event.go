package main

type GameEventWriter interface {
	Write(interface{})
}
