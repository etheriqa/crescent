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

// onAttach removes duplicate disables and triggers EventDisableInterrupt
func (d *disable) onAttach() {
	d.AddEventHandler(d, EventDead)
	d.AddEventHandler(d, EventGameTick)
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
	d.TriggerEvent(EventDisableInterrupt)
}

// onDetach removes the EventHandler
func (d *disable) onDetach() {
	d.RemoveEventHandler(d, EventDead)
	d.RemoveEventHandler(d, EventGameTick)
}

// HandleEvent handles the event
func (d *disable) HandleEvent(e Event) {
	switch e {
	case EventDead:
		d.detachOperator(d)
	case EventGameTick:
		d.expire(d, message{
			// TODO pack message
			t: outDisableEnd,
		})
	}
}
