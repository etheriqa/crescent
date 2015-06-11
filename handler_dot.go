package main

type dotType string

type DoT struct {
	partialHandler
	damage  *damage
	ability *ability
}

// NewDoT returns a DoT
func NewDoT(d *damage, a *ability, duration gameDuration) *DoT {
	return &DoT{
		partialHandler: partialHandler{
			unit:           d.receiver,
			performer:      d.performer,
			expirationTime: d.receiver.after(duration),
		},
		damage:  d,
		ability: a,
	}
}

// OnAttach removes duplicate DoTs
func (d *DoT) OnAttach() {
	d.AddEventHandler(d, EventDead)
	d.AddEventHandler(d, EventGameTick)
	d.AddEventHandler(d, EventXoT)
	for ha := range d.handlers {
		switch ha := ha.(type) {
		case *DoT:
			if ha == d || ha.performer != d.performer || ha.ability != d.ability {
				continue
			}
			if ha.expirationTime > d.expirationTime {
				d.detachHandler(d)
				return
			}
			d.detachHandler(ha)
		}
	}
	d.publish(message{
		// TODO pack message
		t: outDoTBegin,
	})
}

// OnDetach removes the EventHandlers
func (d *DoT) OnDetach() {
	d.RemoveEventHandler(d, EventDead)
	d.RemoveEventHandler(d, EventGameTick)
	d.RemoveEventHandler(d, EventXoT)
}

// HandleEvent handles the event
func (d *DoT) HandleEvent(e Event) {
	switch e {
	case EventDead:
		d.detachHandler(d)
	case EventGameTick:
		d.expire(d, message{
			// TODO pack message
			t: outDoTEnd,
		})
	case EventXoT:
		d.damage.perform(d.game)
	}
}
