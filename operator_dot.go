package main

type dotType string

type dot struct {
	partialOperator
	dotType dotType
	damage  int32
	threat  int32
}

// onAttach removes duplicate DoTs
func (d *dot) onAttach() {
	d.addEventHandler(d, eventDead)
	d.addEventHandler(d, eventGameTick)
	d.addEventHandler(d, eventXoT)
	for o := range d.operators {
		if o == d {
			continue
		}
		if _, ok := o.(*dot); !ok {
			continue
		}
		if o.(*dot).dotType != d.dotType {
			continue
		}
		if o.(*disable).expirationTime >= d.expirationTime {
			d.detachOperator(d)
			return
		}
		d.detachOperator(o)
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
		d.perform()
	}
}

// perform performs the DoT
func (d *dot) perform() {
	d.addHealth(-d.damage)
	d.attachOperator(newThreat(d.unit, d.performer, d.threat))
	d.triggerEvent(eventStats)
}
