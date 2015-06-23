package main

type EffectContainer interface {
	Attach(Effect)
	Detach(Effect)
	Bind(*Unit) EffectContainer
	BindSubject(Subject) EffectContainer
	BindObject(Object) EffectContainer
	Unbind() EffectContainer
	Each(func(Effect))
	Every(func(Effect) bool) bool
	Some(func(Effect) bool) bool
}

type EffectSet map[Effect]bool

type BoundEffectSet struct {
	effects EffectSet
	unit     *Unit
	subject  *Unit
	object   *Unit
}

// MakeEffectSet returns a EffectSet
func MakeEffectSet() EffectSet {
	return make(map[Effect]bool)
}

// Attach adds the Effect if not exists
func (hs EffectSet) Attach(h Effect) {
	if hs[h] {
		return
	}
	hs[h] = true
	h.OnAttach()
}

// Detach removes the Effect if exists
func (hs EffectSet) Detach(h Effect) {
	if !hs[h] {
		return
	}
	delete(hs, h)
	h.OnDetach()
}

// Bind binds the Unit
func (hs EffectSet) Bind(u *Unit) EffectContainer {
	return BoundEffectSet{
		effects: hs,
		unit:     u,
	}
}

// BindSubject binds the Subject
func (hs EffectSet) BindSubject(s Subject) EffectContainer {
	return BoundEffectSet{
		effects: hs,
		subject:  s.Subject(),
	}
}

// BindObject binds the Object
func (hs EffectSet) BindObject(o Object) EffectContainer {
	return BoundEffectSet{
		effects: hs,
		object:   o.Object(),
	}
}

// Unbind unbinds the Units
func (hs EffectSet) Unbind() EffectContainer {
	return hs
}

// Each calls the callback function with each the Effect
func (hs EffectSet) Each(callback func(Effect)) {
	for h := range hs {
		callback(h)
	}
}

// Every returns true if all of the callback result are true
func (hs EffectSet) Every(callback func(Effect) bool) bool {
	for h := range hs {
		if !callback(h) {
			return false
		}
	}
	return true
}

// Some returns true if any of the callback result are true
func (hs EffectSet) Some(callback func(Effect) bool) bool {
	for h := range hs {
		if callback(h) {
			return true
		}
	}
	return false
}

// Attach adds the Effect if not exists
func (bhs BoundEffectSet) Attach(h Effect) {
	bhs.effects.Attach(h)
}

// Detach removes the Effect if exists
func (bhs BoundEffectSet) Detach(h Effect) {
	bhs.effects.Detach(h)
}

// Bind binds the Unit
func (bhs BoundEffectSet) Bind(u *Unit) EffectContainer {
	return BoundEffectSet{
		effects: bhs.effects,
		unit:     u,
		subject:  bhs.subject,
		object:   bhs.object,
	}
}

// BindSubject binds the Subject
func (bhs BoundEffectSet) BindSubject(s Subject) EffectContainer {
	return BoundEffectSet{
		effects: bhs.effects,
		unit:     bhs.unit,
		subject:  s.Subject(),
		object:   bhs.object,
	}
}

// BindObject binds the Object
func (bhs BoundEffectSet) BindObject(o Object) EffectContainer {
	return BoundEffectSet{
		effects: bhs.effects,
		unit:     bhs.unit,
		subject:  bhs.subject,
		object:   o.Object(),
	}
}

// Unbind unbinds the Units
func (bhs BoundEffectSet) Unbind() EffectContainer {
	return bhs.effects
}

// Each calls the callback function with each the Effect
func (bhs BoundEffectSet) Each(callback func(Effect)) {
	bhs.effects.Each(func(h Effect) {
		var subject, object *Unit
		if _, ok := h.(Subject); ok {
			subject = h.(Subject).Subject()
		}
		if _, ok := h.(Object); ok {
			object = h.(Object).Object()
		}
		if bhs.unit != nil && bhs.unit != subject && bhs.unit != object {
			return
		}
		if bhs.subject != nil && bhs.subject != subject {
			return
		}
		if bhs.object != nil && bhs.object != object.Object() {
			return
		}
		callback(h)
	})
}

// Every returns true if all of the callback result are true
func (bhs BoundEffectSet) Every(callback func(Effect) bool) bool {
	ok := true
	bhs.Each(func(h Effect) {
		if !ok || !callback(h) {
			ok = false
		}
	})
	return ok
}

// Some returns true if any of the callback result are true
func (bhs BoundEffectSet) Some(callback func(Effect) bool) bool {
	ok := false
	bhs.Each(func(h Effect) {
		if ok || callback(h) {
			ok = true
		}
	})
	return ok
}
