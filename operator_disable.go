package main

type disableType uint8

const (
	_ disableType = iota
	disableTypeSilence
	disableTypeStun
	disableTypeTaunt
)

type disable struct {
	partialOperator
	disableType disableType
}

// newDisable returns a disable operator
func newDisable(receiver *unit, disableType disableType, duration gameDuration) *disable {
	return &disable{
		partialOperator: partialOperator{
			unit:           receiver,
			expirationTime: receiver.after(duration),
		},
		disableType: disableType,
	}
}

// onAttach removes duplicate disables and triggers eventDisableInterrupt
func (d *disable) onAttach() {
	d.addEventHandler(d, eventDead)
	d.addEventHandler(d, eventGameTick)
	for o := range d.operators {
		switch o := o.(type) {
		case *disable:
			if o == d || o.disableType != d.disableType {
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
		t: outDisableBegin,
	})
	d.triggerEvent(eventDisableInterrupt)
}

// onDetach removes the eventHandler
func (d *disable) onDetach() {
	d.removeEventHandler(d, eventDead)
	d.removeEventHandler(d, eventGameTick)
}

// handleEvent handles the event
func (d *disable) handleEvent(e event) {
	switch e {
	case eventDead:
		d.detachOperator(d)
	case eventGameTick:
		d.expire(d, message{
			// TODO pack message
			t: outDisableEnd,
		})
	}
}
