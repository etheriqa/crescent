package main

type UnitPair interface {
	HandlerContainer
	GameClock
	MessageWriter
	Subject() *Unit
	Object() *Unit
}

type implUnitPair struct {
	*Game
	subject *Unit
	object  *Unit
}

// MakeUnitPair returns a implUnitPair
func MakeUnitPair(subject, object *Unit) implUnitPair {
	if subject == nil {
		log.Panic("The subject must be *Unit")
	}
	if object == nil {
		log.Panic("The object must be *Unit")
	}
	return implUnitPair{
		Game:    subject.Game,
		subject: subject,
		object:  object,
	}
}

// MakeSubject returns a implUnitPair
func MakeSubject(subject *Unit) implUnitPair {
	if subject == nil {
		log.Panic("The subject must be *Unit")
	}
	return implUnitPair{
		Game:    subject.Game,
		subject: subject,
		object:  nil,
	}
}

// MakeObject returns a implUnitPair
func MakeObject(object *Unit) implUnitPair {
	if object == nil {
		log.Panic("The object must be *Unit")
	}
	return implUnitPair{
		Game:    object.Game,
		subject: nil,
		object:  object,
	}
}

// Subject returns the subject unit
func (p implUnitPair) Subject() *Unit {
	return p.subject
}

// Object returns the object unit
func (p implUnitPair) Object() *Unit {
	return p.object
}
