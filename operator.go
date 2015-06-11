package main

type Operator interface {
	Subject() *Unit
	Object() *Unit
	Perform() (before, after Statistic, crit bool, err error)
}
