package main

type operator interface {
	onAttach()
	onDetach()
}

type partialOperator struct {
	unit           *unit
	performer      *unit
	expirationTime gameTime
}

// onAttach does nothing
func (p *partialOperator) onAttach() {}

// onDetach does nothing
func (p *partialOperator) onDetach() {}

// isExpired returns true iff it is expired
func (p *partialOperator) isExpired() bool {
	return p.expirationTime > p.unit.now()
}

// expire expires the operator iff it is expired
func (p *partialOperator) expire(o operator, m message) {
	if p.isExpired() {
		return
	}
	p.unit.detachOperator(o)
	p.unit.publish(m)
}
