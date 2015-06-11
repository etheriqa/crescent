package main

type UnitPair struct {
	*Game
	subject *Unit
	object  *Unit
}

// MakeUnitPair returns a UnitPair
func MakeUnitPair(subject, object *Unit) UnitPair {
	if subject == nil {
		log.Panic("The subject must be *Unit")
	}
	if object == nil {
		log.Panic("The object must be *Unit")
	}
	return UnitPair{
		Game:    subject.Game,
		subject: subject,
		object:  object,
	}
}

// MakeSubject returns a UnitPair
func MakeSubject(subject *Unit) UnitPair {
	if subject == nil {
		log.Panic("The subject must be *Unit")
	}
	return UnitPair{
		Game:    subject.Game,
		subject: subject,
		object:  nil,
	}
}

// MakeObject returns a UnitPair
func MakeObject(object *Unit) UnitPair {
	if object == nil {
		log.Panic("The object must be *Unit")
	}
	return UnitPair{
		Game:    object.Game,
		subject: nil,
		object:  object,
	}
}

// Subject returns the subject unit
func (p UnitPair) Subject() *Unit {
	return p.subject
}

// Object returns the object unit
func (p UnitPair) Object() *Unit {
	return p.object
}

// ForSubjectHandler calls the callback with the handler has the subject
func (p UnitPair) ForSubjectHandler(callback func(Handler)) {
	p.Subject().ForSubjectHandler(callback)
}

// ForObjectHandler calls the callback with the handler has the object
func (p UnitPair) ForObjectHandler(callback func(Handler)) {
	p.Object().ForObjectHandler(callback)
}

// EverySubjectHandler returns true if all of callback results are true
func (p UnitPair) EverySubjectHandler(callback func(Handler) bool) bool {
	return p.Subject().EverySubjectHandler(callback)
}

// EveryObjectHandler returns true if all of callback results are true
func (p UnitPair) EveryObjectHandler(callback func(Handler) bool) bool {
	return p.Object().EveryObjectHandler(callback)
}

// SomeSubjectHandler returns true if any of callback results are true
func (p UnitPair) SomeSubjectHandler(callback func(Handler) bool) bool {
	return p.Subject().SomeSubjectHandler(callback)
}

// SomeObjectHandler returns true if any of callback results are true
func (p UnitPair) SomeObjectHandler(callback func(Handler) bool) bool {
	return p.Object().SomeObjectHandler(callback)
}
