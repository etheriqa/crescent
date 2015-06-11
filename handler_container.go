package main

type HandlerContainer struct {
	handlers map[Handler]bool
}

// NewHandlerContainer returns a HandlerContainer
func NewHandlerContainer() *HandlerContainer {
	return &HandlerContainer{
		handlers: make(map[Handler]bool),
	}
}

// AttachHandler adds the handler
func (hc *HandlerContainer) AttachHandler(ha Handler) {
	if hc.handlers[ha] {
		return
	}
	hc.handlers[ha] = true
	ha.OnAttach()
}

// DetachHandler removes the handler
func (hc *HandlerContainer) DetachHandler(ha Handler) {
	if !hc.handlers[ha] {
		return
	}
	delete(hc.handlers, ha)
	ha.OnDetach()
}

// ForSubjectHandler calls the callback with the handler has given the subject
func (hc *HandlerContainer) ForSubjectHandler(subject *Unit, callback func(Handler)) {
	for ha := range hc.handlers {
		if ha.Subject() == subject {
			callback(ha)
		}
	}
}

// ForObjectHandler calls the callback with the handler has given the object
func (hc *HandlerContainer) ForObjectHandler(object *Unit, callback func(Handler)) {
	for ha := range hc.handlers {
		if ha.Object() == object {
			callback(ha)
		}
	}
}

// EverySubjectHandler returns true if all of callback results are true
func (hc *HandlerContainer) EverySubjectHandler(subject *Unit, callback func(Handler) bool) bool {
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
func (hc *HandlerContainer) EveryObjectHandler(object *Unit, callback func(Handler) bool) bool {
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
func (hc *HandlerContainer) SomeSubjectHandler(subject *Unit, callback func(Handler) bool) bool {
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
func (hc *HandlerContainer) SomeObjectHandler(object *Unit, callback func(Handler) bool) bool {
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
