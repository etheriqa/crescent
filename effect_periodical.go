package main

type Periodical struct {
	UnitPair
	name           string
	routine        func()
	expirationTime InstanceTime

	g Game
}

// EffectDidAttach removes duplicate Periodicals
func (h *Periodical) EffectDidAttach() error {
	ok := h.g.EffectQuery().BindSubject(h).BindObject(h).Every(func(o Effect) bool {
		switch o := o.(type) {
		case *Periodical:
			if h == o || h.name != o.name {
				return true
			}
			if h.expirationTime <= o.expirationTime {
				return false
			}
			h.g.DetachEffect(o)
		}
		return true
	})
	if !ok {
		h.g.DetachEffect(h)
		return nil
	}

	h.writeOutputUnitAttach()
	h.Object().Register(h)
	return nil
}

// EffectDidDetach does nothing
func (h *Periodical) EffectDidDetach() error {
	h.Object().Unregister(h)
	return nil
}

// Handle handles the Event
func (h *Periodical) Handle(p interface{}) {
	switch p.(type) {
	case *EventDead:
		h.g.DetachEffect(h)
	case *EventGameTick:
		if h.g.Clock().Before(h.expirationTime) {
			return
		}
		h.writeOutputUnitDetach()
		h.g.DetachEffect(h)
	case *EventPeriodicalTick:
		h.routine()
	}
}

// writeOutputUnitAttach writes a OutputUnitAttach
func (h *Periodical) writeOutputUnitAttach() {
	h.g.Writer().Write(OutputUnitAttach{
		UnitID:         h.Object().ID(),
		AttachmentName: h.name,
		ExpirationTime: h.expirationTime,
	})
}

// writeOutputUnitDetach writes a OutputUnitDetach
func (h *Periodical) writeOutputUnitDetach() {
	h.g.Writer().Write(OutputUnitDetach{
		UnitID:         h.Object().ID(),
		AttachmentName: h.name,
	})
}
