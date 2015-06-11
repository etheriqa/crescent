package main

type PartialHandler struct {
	UnitPair
	expirationTime GameTime
}

// MakePartialHandler returns a PartialHandler
func MakePartialHandler(up UnitPair, duration GameDuration) PartialHandler {
	return PartialHandler{
		UnitPair:       up,
		expirationTime: up.After(duration),
	}
}

// MakePermanentPartialHandler returns a permanent PartialHandler
func MakePermanentPartialHandler(up UnitPair) PartialHandler {
	return PartialHandler{
		UnitPair:       up,
		expirationTime: -1,
	}
}

// IsExpired returns whether the handler is expired or not
func (p *PartialHandler) IsExpired() bool {
	return p.expirationTime >= 0 && p.expirationTime <= p.Now()
}

// Stop detaches the handler from both the subject and the object
func (p *PartialHandler) Stop(ha Handler) {
	p.DetachHandler(ha)
}
