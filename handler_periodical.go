package main

type Periodical struct {
	UnitPair
	name           string
	routine        func()
	expirationTime InstanceTime

	op Operator
}

// OnAttach removes duplicate Periodicals
func (h *Periodical) OnAttach() {
	ok := h.op.Handlers().BindSubject(h).BindObject(h).Every(func(o Handler) bool {
		switch o := o.(type) {
		case *Periodical:
			if h == o || h.name != o.name {
				return true
			}
			if h.expirationTime <= o.expirationTime {
				return false
			}
			h.op.Handlers().Detach(o)
		}
		return true
	})
	if !ok {
		h.op.Handlers().Detach(h)
		return
	}

	h.writeOutputUnitAttach()
	h.Object().Register(h)
}

// OnDetach does nothing
func (h *Periodical) OnDetach() {
	h.Object().Unregister(h)
}

// Handle handles the Event
func (h *Periodical) Handle(p interface{}) {
	switch p.(type) {
	case *EventDead:
		h.op.Handlers().Detach(h)
	case *EventGameTick:
		if h.op.Clock().Before(h.expirationTime) {
			return
		}
		h.writeOutputUnitDetach()
		h.op.Handlers().Detach(h)
	case *EventPeriodicalTick:
		h.routine()
	}
}

// writeOutputUnitAttach writes a OutputUnitAttach
func (h *Periodical) writeOutputUnitAttach() {
	h.op.Writer().Write(OutputUnitAttach{
		UnitID:         h.Object().ID(),
		AttachmentName: h.name,
		ExpirationTime: h.expirationTime,
	})
}

// writeOutputUnitDetach writes a OutputUnitDetach
func (h *Periodical) writeOutputUnitDetach() {
	h.op.Writer().Write(OutputUnitDetach{
		UnitID:         h.Object().ID(),
		AttachmentName: h.name,
	})
}
