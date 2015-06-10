package main

type operator interface {
	onAttach()
	onDetach()
}

type partialOperator struct {
	*unit
	performer      *unit
	expirationTime gameTime
}

// onAttach does nothing
func (p *partialOperator) onAttach() {}

// onDetach does nothing
func (p *partialOperator) onDetach() {}

// isExpired returns true iff it is expired
func (p *partialOperator) isExpired() bool {
	return p.expirationTime != 0 && p.expirationTime > p.now()
}

// expire expires the operator iff it is expired
func (p *partialOperator) expire(o operator, m message) {
	if p.isExpired() {
		return
	}
	p.detachOperator(o)
	p.publish(m)
}
