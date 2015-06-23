package crescent

type Subject interface {
	Subject() *Unit
}

type Object interface {
	Object() *Unit
}

type UnitSubject struct {
	subject *Unit
}

type UnitObject struct {
	object *Unit
}

type UnitPair struct {
	subject *Unit
	object  *Unit
}

// MakeSubject returns a UnitSubject
func MakeSubject(s Subject) UnitSubject {
	return UnitSubject{
		subject: s.Subject(),
	}
}

// Subject returns the Subject
func (s UnitSubject) Subject() *Unit {
	return s.subject
}

// MakeObject returns a UnitObject
func MakeObject(o Object) UnitObject {
	return UnitObject{
		object: o.Object(),
	}
}

// Object returns the Object
func (o UnitObject) Object() *Unit {
	return o.object
}

// MakePair returns a UnitPair
func MakePair(s Subject, o Object) UnitPair {
	return UnitPair{
		subject: s.Subject(),
		object:  o.Object(),
	}
}

// Subject returns the Subject
func (p UnitPair) Subject() *Unit {
	return p.subject
}

// Object returns the Object
func (p UnitPair) Object() *Unit {
	return p.object
}
