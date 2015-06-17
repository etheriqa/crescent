package main

import (
	"errors"
)

const GroupCapacity = 5

type UnitContainer interface {
	Join(UnitGroup, UnitName, *Class) (*Unit, error)
	Leave(UnitID) error
	Find(UnitID) *Unit
	Each(func(*Unit))
	EachFriend(*Unit, func(*Unit))
	EachEnemy(*Unit, func(*Unit))
}

type UnitMap struct {
	seq UnitID
	id  map[UnitID]*Unit
}

// MakeUnitMap returns a UnitMap
func MakeUnitMap() UnitMap {
	return UnitMap{
		seq: 0,
		id:  make(map[UnitID]*Unit),
	}
}

// Join creates a Unit and adds it
func (um UnitMap) Join(group UnitGroup, name UnitName, class *Class) (*Unit, error) {
	u := NewUnit(um.generateUnitID(), group, name, class)
	n := 0
	um.EachFriend(u, func(*Unit) { n++ })
	if n >= GroupCapacity {
		return nil, errors.New("There is no group capacity")
	}
	um.id[u.ID()] = u
	return u, nil
}

// Leave removes the Unit
func (um UnitMap) Leave(id UnitID) error {
	if um.id[id] == nil {
		return errors.New("Unknown UnitID")
	}
	delete(um.id, id)
	return nil
}

// Find finds the Unit with the UnitID
func (um UnitMap) Find(id UnitID) *Unit {
	return um.id[id]
}

// Each calls the callback function with each the Unit
func (um UnitMap) Each(callback func(*Unit)) {
	for _, u := range um.id {
		callback(u)
	}
}

// EachFriend calls the callback function with each the friend Unit
func (um UnitMap) EachFriend(u *Unit, callback func(*Unit)) {
	for _, v := range um.id {
		if u.Group() == v.Group() {
			callback(u)
		}
	}
}

// EachEnemy calls the callback function with each the enemy Unit
func (um UnitMap) EachEnemy(u *Unit, callback func(*Unit)) {
	for _, v := range um.id {
		if u.Group() != v.Group() {
			callback(u)
		}
	}
}

// generateUnitID generates a unique UnitID
func (um UnitMap) generateUnitID() UnitID {
	um.seq++
	return um.seq
}
