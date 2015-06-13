package main

type UnitContainer interface {
	Each(func(*Unit))
	EachFriend(*Unit, func(*Unit))
	EachEnemy(*Unit, func(*Unit))
}
