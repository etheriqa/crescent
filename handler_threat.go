package main

type Threat struct {
	UnitPair
	threat Statistic

	op Operator
}

// OnAttach merges Threat handlers
func (h *Threat) OnAttach() {
	h.op.Handlers().BindSubject(h).BindObject(h).Each(func(o Handler) {
		switch o := o.(type) {
		case *Threat:
			if o == h {
				return
			}
			h.threat += o.threat
			h.op.Handlers().Detach(o)
		}
	})
}

// OnDetach does nothing
func (h *Threat) OnDetach() {
}
