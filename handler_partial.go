package main

type PartialHandler struct {
	subject        *unit
	object         *unit
	expirationTime gameTime
}

// NewPartialHandler returns a PartialHandler
func NewPartialHandler(subject, object *unit, duration gameDuration) *PartialHandler {
	ha := NewPermanentPartialHandler(subject, object)
	if subject != nil {
		ha.expirationTime = subject.after(duration)
	}
	if object != nil {
		ha.expirationTime = object.after(duration)
	}
	return ha
}

// NewPermanentPartialHandler returns a permanent PartialHandler
func NewPermanentPartialHandler(subject, object *unit) *PartialHandler {
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
func (p *PartialHandler) Subject() *unit {
	return p.subject
}

// Object returns the object unit
func (p *PartialHandler) Object() *unit {
	return p.object
}

// Now returns the current game time
func (p *PartialHandler) Now() gameTime {
	if p.subject != nil {
		return p.subject.now()
	}
	if p.object != nil {
		return p.object.now()
	}
	// TODO return error
	log.Fatal("")
	return gameTime(0)
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
		p.subject.publish(m)
		return
	}
	if p.object != nil {
		p.object.publish(m)
		return
	}
}
