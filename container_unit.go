package main

type UnitContainer interface {
	Each(func(*Unit))
	EachFriend(*Unit, func(*Unit))
	EachEnemy(*Unit, func(*Unit))
}

type UnitMap map[UnitID]*Unit

// MakeUnitMap returns a UnitMap
func MakeUnitMap() UnitMap {
	return make(map[UnitID]*Unit)
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
