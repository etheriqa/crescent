package game

import (
	"errors"
)

type EffectContainer interface {
	Attach(Game, Effect) error
	Detach(Game, Effect) error
}

type EffectQueryable interface {
	Bind(*Unit) EffectQueryable
	BindSubject(Subject) EffectQueryable
	BindObject(Object) EffectQueryable
	Unbind() EffectQueryable
	Each(func(Effect))
	Every(func(Effect) bool) bool
	Some(func(Effect) bool) bool
}

type EffectSet map[Effect]bool

type BoundEffectSet struct {
	effects EffectSet
	unit    *Unit
	subject *Unit
	object  *Unit
}

// MakeEffectSet returns a EffectSet
func MakeEffectSet() EffectSet {
	return make(map[Effect]bool)
}

// Attach adds the Effect if not contains
func (es EffectSet) Attach(g Game, e Effect) error {
	if es[e] {
		return errors.New("Already attached")
	}
	if e, ok := e.(EffectWillAttach); ok {
		if err := e.EffectWillAttach(g); err != nil {
			return err
		}
	}
	es[e] = true
	if e, ok := e.(EffectDidAttach); ok {
		if err := e.EffectDidAttach(g); err != nil {
			return err
		}
	}
	return nil
}

// Detach removes the Effect if contains
func (es EffectSet) Detach(g Game, e Effect) error {
	if !es[e] {
		return errors.New("Never attached")
	}
	if e, ok := e.(EffectWillDetach); ok {
		if err := e.EffectWillDetach(g); err != nil {
			return err
		}
	}
	delete(es, e)
	if e, ok := e.(EffectDidDetach); ok {
		if err := e.EffectDidDetach(g); err != nil {
			return err
		}
	}
	return nil
}

// Bind binds the Unit
func (es EffectSet) Bind(u *Unit) EffectQueryable {
	return BoundEffectSet{
		effects: es,
		unit:    u,
	}
}

// BindSubject binds the Subject
func (es EffectSet) BindSubject(s Subject) EffectQueryable {
	return BoundEffectSet{
		effects: es,
		subject: s.Subject(),
	}
}

// BindObject binds the Object
func (es EffectSet) BindObject(o Object) EffectQueryable {
	return BoundEffectSet{
		effects: es,
		object:  o.Object(),
	}
}

// Unbind unbinds the Units
func (es EffectSet) Unbind() EffectQueryable {
	return es
}

// Each calls the callback function with each the Effect
func (es EffectSet) Each(callback func(Effect)) {
	for e := range es {
		callback(e)
	}
}

// Every returns true if all of the callback result are true
func (es EffectSet) Every(callback func(Effect) bool) bool {
	for e := range es {
		if !callback(e) {
			return false
		}
	}
	return true
}

// Some returns true if any of the callback result are true
func (es EffectSet) Some(callback func(Effect) bool) bool {
	for e := range es {
		if callback(e) {
			return true
		}
	}
	return false
}

// Bind binds the Unit
func (bes BoundEffectSet) Bind(u *Unit) EffectQueryable {
	return BoundEffectSet{
		effects: bes.effects,
		unit:    u,
		subject: bes.subject,
		object:  bes.object,
	}
}

// BindSubject binds the Subject
func (bes BoundEffectSet) BindSubject(s Subject) EffectQueryable {
	return BoundEffectSet{
		effects: bes.effects,
		unit:    bes.unit,
		subject: s.Subject(),
		object:  bes.object,
	}
}

// BindObject binds the Object
func (bes BoundEffectSet) BindObject(o Object) EffectQueryable {
	return BoundEffectSet{
		effects: bes.effects,
		unit:    bes.unit,
		subject: bes.subject,
		object:  o.Object(),
	}
}

// Unbind unbinds the Units
func (bes BoundEffectSet) Unbind() EffectQueryable {
	return bes.effects
}

// Each calls the callback function with each the Effect
func (bes BoundEffectSet) Each(callback func(Effect)) {
	bes.effects.Each(func(e Effect) {
		var subject, object *Unit
		if _, ok := e.(Subject); ok {
			subject = e.(Subject).Subject()
		}
		if _, ok := e.(Object); ok {
			object = e.(Object).Object()
		}
		if bes.unit != nil && bes.unit != subject && bes.unit != object {
			return
		}
		if bes.subject != nil && bes.subject != subject {
			return
		}
		if bes.object != nil && bes.object != object.Object() {
			return
		}
		callback(e)
	})
}

// Every returns true if all of the callback result are true
func (bes BoundEffectSet) Every(callback func(Effect) bool) bool {
	ok := true
	bes.Each(func(e Effect) {
		if !ok || !callback(e) {
			ok = false
		}
	})
	return ok
}

// Some returns true if any of the callback result are true
func (bes BoundEffectSet) Some(callback func(Effect) bool) bool {
	ok := false
	bes.Each(func(e Effect) {
		if ok || callback(e) {
			ok = true
		}
	})
	return ok
}
