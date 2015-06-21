package main

type StageID uint64

type Stage interface {
	Initialize(Operator) error
	OnTick(Operator)
}
