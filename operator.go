package main

type operator interface {
	onAttach()
	onDetach()
}

type partialOperator struct {
	unit           *unit
	expirationTime gameTime
}

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
