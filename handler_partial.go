package main

type PartialHandler struct {
	subject        *Unit
	object         *Unit
	expirationTime GameTime
}

// NewPartialHandler returns a PartialHandler
func NewPartialHandler(subject, object *Unit, duration GameDuration) *PartialHandler {
	ha := NewPermanentPartialHandler(subject, object)
	if subject != nil {
		ha.expirationTime = subject.After(duration)
	}
	if object != nil {
		ha.expirationTime = object.After(duration)
	}
	return ha
}

// NewPermanentPartialHandler returns a permanent PartialHandler
func NewPermanentPartialHandler(subject, object *Unit) *PartialHandler {
	return &PartialHandler{
		subject:        subject,
		object:         object,
		expirationTime: -1,
	}
}

// Container returns the HandlerContainer
func (p *PartialHandler) Container() HandlerContainer {
	if p.subject != nil {
		return p.subject
	}
	if p.object != nil {
		return p.object
	}
	// TODO return error
	log.Fatal("")
	return nil
}

// Subject returns the subject unit
func (p *PartialHandler) Subject() *Unit {
	return p.subject
}

// Object returns the object unit
func (p *PartialHandler) Object() *Unit {
	return p.object
}

// Now returns the current game time
func (p *PartialHandler) Now() GameTime {
	if p.subject != nil {
		return p.subject.Now()
	}
	if p.object != nil {
		return p.object.Now()
	}
	// TODO return error
	log.Fatal("")
	return GameTime(0)
}

// IsExpired returns whether the handler is expired or not
func (p *PartialHandler) IsExpired() bool {
	return p.expirationTime >= 0 && p.expirationTime <= p.Now()
}

// Stop detaches the handler from both the subject and the object
func (p *PartialHandler) Stop(ha Handler) {
	p.Container().DetachHandler(ha)
}

// Publish sends the message
func (p *PartialHandler) Publish(m message) {
	if p.subject != nil {
		p.subject.Publish(m)
		return
	}
	if p.object != nil {
		p.object.Publish(m)
		return
	}
}
