package main

type Threat struct {
	UnitPair
	threat Statistic

	op Operator
}

// OnAttach merges Threat effects
func (h *Threat) OnAttach() {
	h.op.Effects().BindSubject(h).BindObject(h).Each(func(o Effect) {
		switch o := o.(type) {
		case *Threat:
			if o == h {
				return
			}
			h.threat += o.threat
			h.op.Effects().Detach(o)
		}
	})
}

// OnDetach does nothing
func (h *Threat) OnDetach() {
}
