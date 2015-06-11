package main

type disableType uint8

const (
	_ disableType = iota
	disableTypeSilence
	disableTypeStun
	disableTypeTaunt
)

type Disable struct {
	*PartialHandler
	disableType
}

// NewDisable returns a Disable handler
func NewDisable(object *unit, disableType disableType, duration gameDuration) *Disable {
	return &Disable{
		PartialHandler: NewPartialHandler(nil, object, duration),
		disableType:    disableType,
	}
}

// OnAttach removes duplicate disables and triggers EventDisableInterrupt
func (d *Disable) OnAttach() {
	d.Object().AddEventHandler(d, EventDead)
	d.Object().AddEventHandler(d, EventGameTick)
	ok := d.Container().AllObjectHandler(d.Object(), func(ha Handler) bool {
		switch ha := ha.(type) {
		case *Disable:
			if ha == d || ha.disableType != d.disableType {
				return true
			}
			if ha.expirationTime > d.expirationTime {
				return false
			}
			ha.Stop(ha)
		}
		return true
	})
	if !ok {
		d.Stop(d)
		return
	}
	d.Publish(message{
	// TODO pack message
	})
	d.Object().TriggerEvent(EventDisableInterrupt)
}

// OnDetach removes the EventHandler
func (d *Disable) OnDetach() {
	d.Object().RemoveEventHandler(d, EventDead)
	d.Object().RemoveEventHandler(d, EventGameTick)
}

// HandleEvent handles the Event
func (d *Disable) HandleEvent(e Event) {
	switch e {
	case EventDead:
		d.Stop(d)
	case EventGameTick:
		if d.IsExpired() {
			d.Up()
		}
	}
}

// Up ends the disable
func (d *Disable) Up() {
	d.Stop(d)
	d.Publish(message{
	// TODO pack message
	})
}
