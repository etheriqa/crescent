package main

type Threat struct {
	UnitPair
	threat Statistic

	handlers HandlerContainer
}

// OnAttach merges Threat handlers
func (h *Threat) OnAttach() {
	h.handlers.BindSubject(h).BindObject(h).Each(func(o Handler) {
		switch o := o.(type) {
		case *Threat:
			if o == h {
				return
			}
			h.threat += o.threat
			h.handlers.Detach(o)
		}
	})
}

// OnDetach does nothing
func (h *Threat) OnDetach() {
}
