package main

type Handler interface {
	EventHandler
	Subject() *unit
	Object() *unit
	OnAttach()
	OnDetach()
}

type HandlerContainer interface {
	AttachHandler(Handler)
	DetachHandler(Handler)
	ForSubjectHandler(*unit, func(Handler))
	ForObjectHandler(*unit, func(Handler))
	EverySubjectHandler(*unit, func(Handler) bool) bool
	EveryObjectHandler(*unit, func(Handler) bool) bool
}

type implHandlerContainer struct {
	handlers map[Handler]bool
}

// NewHandlerContainer returns a implHandlerContainer
func NewHandlerContainer() *implHandlerContainer {
	return &implHandlerContainer{
		handlers: make(map[Handler]bool),
	}
}

// AttachHandler adds the handler
func (hc *implHandlerContainer) AttachHandler(ha Handler) {
	if hc.handlers[ha] {
		return
	}
	hc.handlers[ha] = true
	ha.OnAttach()
}

// DetachHandler removes the handler
func (hc *implHandlerContainer) DetachHandler(ha Handler) {
	if !hc.handlers[ha] {
		return
	}
	delete(hc.handlers, ha)
	ha.OnDetach()
}

// ForSubjectHandler calls the callback with the handler has given the subject
func (hc *implHandlerContainer) ForSubjectHandler(subject *unit, callback func(Handler)) {
	for ha := range hc.handlers {
		if ha.Subject() == subject {
			callback(ha)
		}
	}
}

// ForObjectHandler calls the callback with the handler has given the object
func (hc *implHandlerContainer) ForObjectHandler(object *unit, callback func(Handler)) {
	for ha := range hc.handlers {
		if ha.Object() == object {
			callback(ha)
		}
	}
}

// EverySubjectHandler returns true if all of callback results are true
func (hc *implHandlerContainer) EverySubjectHandler(subject *unit, callback func(Handler) bool) bool {
	for ha := range hc.handlers {
		if ha.Subject() != subject {
			continue
		}
		if !callback(ha) {
			return false
		}
	}
	return true
}

// EveryObjectHandler returns true if all of callback results are true
func (hc *implHandlerContainer) EveryObjectHandler(object *unit, callback func(Handler) bool) bool {
	for ha := range hc.handlers {
		if ha.Object() != object {
			continue
		}
		if !callback(ha) {
			return false
		}
	}
	return true
}

// SomeSubjectHandler returns true if any of callback results are true
func (hc *implHandlerContainer) SomeSubjectHandler(subject *unit, callback func(Handler) bool) bool {
	for ha := range hc.handlers {
		if ha.Subject() != subject {
			continue
		}
		if callback(ha) {
			return true
		}
	}
	return false
}

// SomeObjectHandler returns true if any of callback results are true
func (hc *implHandlerContainer) SomeObjectHandler(object *unit, callback func(Handler) bool) bool {
	for ha := range hc.handlers {
		if ha.Object() != object {
			continue
		}
		if callback(ha) {
			return true
		}
	}
	return false
}
