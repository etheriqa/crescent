package main

type DisableType uint8

const (
	_ DisableType = iota
	DisableTypeSilence
	DisableTypeStun
	DisableTypeTaunt
)

type Disable struct {
	PartialHandler
	DisableType
}

// NewDisable returns a Disable handler
func NewDisable(object *Unit, dt DisableType, duration GameDuration) *Disable {
	return &Disable{
		PartialHandler: MakePartialHandler(MakeObject(object), duration),
		DisableType:    dt,
	}
}

// OnAttach removes duplicate disables and triggers EventDisableInterrupt
func (d *Disable) OnAttach() {
	d.Object().AddEventHandler(d, EventDead)
	d.Object().AddEventHandler(d, EventGameTick)
	ok := d.EveryObjectHandler(d.Object(), func(ha Handler) bool {
		switch ha := ha.(type) {
		case *Disable:
			if ha == d || ha.DisableType != d.DisableType {
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
