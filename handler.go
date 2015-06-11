package main

type Handler interface {
	EventHandler
	OnAttach()
	OnDetach()
}

type partialHandler struct {
	*unit
	performer      *unit
	expirationTime gameTime
}

// OnAttach does nothing
func (p *partialHandler) OnAttach() {}

// OnDetach does nothing
func (p *partialHandler) OnDetach() {}

// isExpired returns true iff it is expired
func (p *partialHandler) isExpired() bool {
	return p.expirationTime != 0 && p.expirationTime > p.now()
}

// expire expires the handler iff it is expired
func (p *partialHandler) expire(ha Handler, m message) {
	if p.isExpired() {
		return
	}
	p.detachHandler(ha)
	p.publish(m)
}
