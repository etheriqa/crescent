package main

type disableType uint8

const (
	_ disableType = iota
	disableTypeSilence
	disableTypeStun
	disableTypeTaunt
)

type Disable struct {
	partialHandler
	disableType disableType
}

// NewDisable returns a disable handler
func NewDisable(receiver *unit, disableType disableType, duration gameDuration) *Disable {
	return &Disable{
		partialHandler: partialHandler{
			unit:           receiver,
			expirationTime: receiver.after(duration),
		},
		disableType: disableType,
	}
}

// OnAttach removes duplicate disables and triggers EventDisableInterrupt
func (d *Disable) OnAttach() {
	d.AddEventHandler(d, EventDead)
	d.AddEventHandler(d, EventGameTick)
	for ha := range d.handlers {
		switch ha := ha.(type) {
		case *Disable:
			if ha == d || ha.disableType != d.disableType {
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
		t: outDisableBegin,
	})
	d.TriggerEvent(EventDisableInterrupt)
}

// OnDetach removes the EventHandler
func (d *Disable) OnDetach() {
	d.RemoveEventHandler(d, EventDead)
	d.RemoveEventHandler(d, EventGameTick)
}

// HandleEvent handles the event
func (d *Disable) HandleEvent(e Event) {
	switch e {
	case EventDead:
		d.detachHandler(d)
	case EventGameTick:
		d.expire(d, message{
			// TODO pack message
			t: outDisableEnd,
		})
	}
}
