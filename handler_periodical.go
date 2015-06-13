package main

type Periodical struct {
	UnitPair
	name           string
	routine        func()
	expirationTime GameTime

	clock    GameClock
	handlers HandlerContainer
	writer   GameEventWriter
}

// OnAttach removes duplicate Periodicals
func (h *Periodical) OnAttach() {
	ok := h.handlers.BindSubject(h).BindObject(h).Every(func(o Handler) bool {
		switch o := o.(type) {
		case *Periodical:
			if h == o || h.name != o.name {
				return true
			}
			if h.expirationTime <= o.expirationTime {
				return false
			}
			h.handlers.Detach(o)
		}
		return true
	})
	if !ok {
		h.handlers.Detach(h)
		return
	}

	h.Object().AddEventHandler(h, EventGameTick)
	h.Object().AddEventHandler(h, EventPeriodicalTick)
	h.Object().AddEventHandler(h, EventDead)
	h.writer.Write(nil) // TODO
}

// OnDetach does nothing
func (h *Periodical) OnDetach() {
	h.Object().RemoveEventHandler(h, EventGameTick)
	h.Object().RemoveEventHandler(h, EventPeriodicalTick)
	h.Object().RemoveEventHandler(h, EventDead)
}

// HandleEvent handles the Event
func (h *Periodical) HandleEvent(e Event) {
	switch e {
	case EventDead:
		h.handlers.Detach(h)
	case EventGameTick:
		if h.clock.Before(h.expirationTime) {
			return
		}
		h.handlers.Detach(h)
		h.writer.Write(nil) // TODO
	case EventPeriodicalTick:
		h.routine()
	}
}
