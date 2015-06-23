package main

type Threat struct {
	UnitPair
	threat Statistic

	g Game
}

// EffectDidAttach merges Threat effects
func (h *Threat) EffectDidAttach() error {
	h.g.Effects().BindSubject(h).BindObject(h).Each(func(o Effect) {
		switch o := o.(type) {
		case *Threat:
			if o == h {
				return
			}
			h.threat += o.threat
			h.g.Effects().Detach(o)
		}
	})
	return nil
}
