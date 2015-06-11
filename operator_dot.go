package main

type dotType string

type dot struct {
	partialOperator
	damage  *damage
	ability *ability
}

// newDoT returns a DoT
func newDoT(d *damage, a *ability, duration gameDuration) *dot {
	return &dot{
		partialOperator: partialOperator{
			unit:           d.receiver,
			performer:      d.performer,
			expirationTime: d.receiver.after(duration),
		},
		damage:  d,
		ability: a,
	}
}

// onAttach removes duplicate DoTs
func (d *dot) onAttach() {
	d.AddEventHandler(d, EventDead)
	d.AddEventHandler(d, EventGameTick)
	d.AddEventHandler(d, EventXoT)
	for o := range d.operators {
		switch o := o.(type) {
		case *dot:
			if o == d || o.performer != d.performer || o.ability != d.ability {
				continue
			}
			if o.expirationTime > d.expirationTime {
				d.detachOperator(d)
				return
			}
			d.detachOperator(o)
		}
	}
	d.publish(message{
		// TODO pack message
		t: outDoTBegin,
	})
}

// onDetach removes the EventHandlers
func (d *dot) onDetach() {
	d.RemoveEventHandler(d, EventDead)
	d.RemoveEventHandler(d, EventGameTick)
	d.RemoveEventHandler(d, EventXoT)
}

// HandleEvent handles the event
func (d *dot) HandleEvent(e Event) {
	switch e {
	case EventDead:
		d.detachOperator(d)
	case EventGameTick:
		d.expire(d, message{
			// TODO pack message
			t: outDoTEnd,
		})
	case EventXoT:
		d.damage.perform(d.game)
	}
}
