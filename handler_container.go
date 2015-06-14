package main

type HandlerContainer interface {
	Attach(Handler)
	Detach(Handler)
	Bind(*Unit) HandlerContainer
	BindSubject(Subject) HandlerContainer
	BindObject(Object) HandlerContainer
	Unbind() HandlerContainer
	Each(func(Handler))
	Every(func(Handler) bool) bool
	Some(func(Handler) bool) bool
}

type HandlerSet map[Handler]bool

type BoundHandlerSet struct {
	handlers HandlerSet
	unit     *Unit
	subject  *Unit
	object   *Unit
}

// MakeHandlerSet returns a HandlerSet
func MakeHandlerSet() HandlerSet {
	return make(map[Handler]bool)
}

// Attach adds the Handler if not exists
func (hs HandlerSet) Attach(h Handler) {
	if hs[h] {
		return
	}
	hs[h] = true
	h.OnAttach()
}

// Detach removes the Handler if exists
func (hs HandlerSet) Detach(h Handler) {
	if !hs[h] {
		return
	}
	delete(hs, h)
	h.OnDetach()
}

// Bind binds the Unit
func (hs HandlerSet) Bind(u *Unit) HandlerContainer {
	return BoundHandlerSet{
		handlers: hs,
		unit:     u,
	}
}

// BindSubject binds the Subject
func (hs HandlerSet) BindSubject(s Subject) HandlerContainer {
	return BoundHandlerSet{
		handlers: hs,
		subject:  s.Subject(),
	}
}

// BindObject binds the Object
func (hs HandlerSet) BindObject(o Object) HandlerContainer {
	return BoundHandlerSet{
		handlers: hs,
		object:   o.Object(),
	}
}

// Unbind unbinds the Units
func (hs HandlerSet) Unbind() HandlerContainer {
	return hs
}

// Each calls the callback function with each the Handler
func (hs HandlerSet) Each(callback func(Handler)) {
	for h := range hs {
		callback(h)
	}
}

// Every returns true if all of the callback result are true
func (hs HandlerSet) Every(callback func(Handler) bool) bool {
	for h := range hs {
		if !callback(h) {
			return false
		}
	}
	return true
}

// Some returns true if any of the callback result are true
func (hs HandlerSet) Some(callback func(Handler) bool) bool {
	for h := range hs {
		if callback(h) {
			return true
		}
	}
	return false
}

// Attach adds the Handler if not exists
func (bhs BoundHandlerSet) Attach(h Handler) {
	bhs.handlers.Attach(h)
}

// Detach removes the Handler if exists
func (bhs BoundHandlerSet) Detach(h Handler) {
	bhs.handlers.Detach(h)
}

// Bind binds the Unit
func (bhs BoundHandlerSet) Bind(u *Unit) HandlerContainer {
	return BoundHandlerSet{
		handlers: bhs.handlers,
		unit:     u,
		subject:  bhs.subject,
		object:   bhs.object,
	}
}

// BindSubject binds the Subject
func (bhs BoundHandlerSet) BindSubject(s Subject) HandlerContainer {
	return BoundHandlerSet{
		handlers: bhs.handlers,
		unit:     bhs.unit,
		subject:  s.Subject(),
		object:   bhs.object,
	}
}

// BindObject binds the Object
func (bhs BoundHandlerSet) BindObject(o Object) HandlerContainer {
	return BoundHandlerSet{
		handlers: bhs.handlers,
		unit:     bhs.unit,
		subject:  bhs.subject,
		object:   o.Object(),
	}
}

// Unbind unbinds the Units
func (bhs BoundHandlerSet) Unbind() HandlerContainer {
	return bhs.handlers
}

// Each calls the callback function with each the Handler
func (bhs BoundHandlerSet) Each(callback func(Handler)) {
	bhs.handlers.Each(func(h Handler) {
		subject, _ := h.(Subject)
		object, _ := h.(Object)
		if bhs.unit != nil && bhs.unit != subject && bhs.unit != object {
			return
		}
		if bhs.subject != nil && bhs.subject != subject {
			return
		}
		if bhs.object != nil && bhs.object != object {
			return
		}
		callback(h)
	})
}

// Every returns true if all of the callback result are true
func (bhs BoundHandlerSet) Every(callback func(Handler) bool) bool {
	return bhs.handlers.Every(func(h Handler) bool {
		subject, _ := h.(Subject)
		object, _ := h.(Object)
		if bhs.unit != nil && bhs.unit != subject && bhs.unit != object {
			return true
		}
		if bhs.subject != nil && bhs.subject != subject {
			return true
		}
		if bhs.object != nil && bhs.object != object {
			return true
		}
		return callback(h)
	})
}

// Some returns true if any of the callback result are true
func (bhs BoundHandlerSet) Some(callback func(Handler) bool) bool {
	return bhs.handlers.Every(func(h Handler) bool {
		subject, _ := h.(Subject)
		object, _ := h.(Object)
		if bhs.unit != nil && bhs.unit != subject && bhs.unit != object {
			return false
		}
		if bhs.subject != nil && bhs.subject != subject {
			return false
		}
		if bhs.object != nil && bhs.object != object {
			return false
		}
		return callback(h)
	})
}
