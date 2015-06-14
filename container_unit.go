package main

import (
	"errors"
)

type UnitContainer interface {
	Join(UnitGroup, UnitName, *Class) (*Unit, error)
	Leave(UnitID) error
	Find(UnitID) *Unit
	Each(func(*Unit))
	EachFriend(*Unit, func(*Unit))
	EachEnemy(*Unit, func(*Unit))
}

type UnitMap map[UnitID]*Unit

// MakeUnitMap returns a UnitMap
func MakeUnitMap() UnitMap {
	return make(map[UnitID]*Unit)
}

// Join adds the Unit
func (um UnitMap) Join(group UnitGroup, name UnitName, class *Class) (*Unit, error) {
	// TODO WIP
	u := NewUnit(0, group, 0, name, class)
	// TODO ctor
	um[u.ID()] = u
	return u, nil
}

// Leave removes the Unit
func (um UnitMap) Leave(id UnitID) error {
	// TODO WIP
	if um[id] == nil {
		return errors.New("Unknown UnitID")
	}
	// TODO dtor
	delete(um, id)
	return nil
}

// Find finds the Unit with the UnitID
func (um UnitMap) Find(id UnitID) *Unit {
	return um[id]
}

// Each calls the callback function with each the Unit
func (um UnitMap) Each(callback func(*Unit)) {
	for _, u := range um {
		callback(u)
	}
}

// EachFriend calls the callback function with each the friend Unit
func (um UnitMap) EachFriend(u *Unit, callback func(*Unit)) {
	for _, v := range um {
		if u.Group() == v.Group() {
			callback(u)
		}
	}
}

// EachEnemy calls the callback function with each the enemy Unit
func (um UnitMap) EachEnemy(u *Unit, callback func(*Unit)) {
	for _, v := range um {
		if u.Group() != v.Group() {
			callback(u)
		}
	}
}
