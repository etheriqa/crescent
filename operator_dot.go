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
	d.unit.addEventHandler(d, eventGameTick)
	d.unit.addEventHandler(d, eventStatsTick)
	for o := range d.unit.operators {
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
			d.unit.detachOperator(d)
			return
		}
		d.unit.detachOperator(o)
	}
	d.unit.publish(message{
		// TODO pack message
		t: outDoTBegin,
	})
}

// onDetach removes the eventHandlers
func (d *dot) onDetach() {
	d.unit.removeEventHandler(d, eventGameTick)
	d.unit.removeEventHandler(d, eventStatsTick)
}

// handleEvent handles the event
func (d *dot) handleEvent(e event) {
	switch e {
	case eventGameTick:
		d.expire(d, message{
			// TODO pack message
			t: outDoTEnd,
		})
	case eventStatsTick:
		d.perform()
	}
}

// perform performs the DoT
func (d *dot) perform() {
	if d.unit.isDead() {
		return
	}
	d.unit.addHealth(-d.damage)
	d.unit.attachOperator(newThreat(d.unit, d.performer, d.threat))
}
