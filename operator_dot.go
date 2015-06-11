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
	d.addEventHandler(d, eventDead)
	d.addEventHandler(d, eventGameTick)
	d.addEventHandler(d, eventXoT)
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

// onDetach removes the eventHandlers
func (d *dot) onDetach() {
	d.removeEventHandler(d, eventDead)
	d.removeEventHandler(d, eventGameTick)
	d.removeEventHandler(d, eventXoT)
}

// handleEvent handles the event
func (d *dot) handleEvent(e event) {
	switch e {
	case eventDead:
		d.detachOperator(d)
	case eventGameTick:
		d.expire(d, message{
			// TODO pack message
			t: outDoTEnd,
		})
	case eventXoT:
		d.damage.perform(d.game)
	}
}
