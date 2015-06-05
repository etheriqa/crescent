package main

type disableType uint8

const (
	_ disableType = iota
	disableSilence
	disableStun
	disableTaunt
)

type disable struct {
	partialOperator
	disableType disableType
}

// onAttach removes duplicate disables and triggers eventDisable
func (d *disable) onAttach() {
	d.unit.addEventHandler(d, eventGameTick)
	for o := range d.unit.operators {
		if o == d {
			continue
		}
		if _, ok := o.(*disable); !ok {
			continue
		}
		if o.(*disable).disableType != d.disableType {
			continue
		}
		if o.(*disable).expirationTime >= d.expirationTime {
			d.unit.detachOperator(d)
			return
		}
		d.unit.detachOperator(o)
	}
	d.unit.publish(message{
		// todo pack message
		t: outDisableBegin,
	})
	d.unit.triggerEvent(eventDisable)
}

// onDetach removes the eventHandler
func (d *disable) onDetach() {
	d.unit.removeEventHandler(d, eventGameTick)
}

// handleEvent handles the event
func (d *disable) handleEvent(e event) {
	switch e {
	case eventGameTick:
		d.expire(d, message{
			// todo pack message
			t: outDisableEnd,
		})
	}
}
