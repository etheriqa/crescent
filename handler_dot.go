package main

type dotType string

type DoT struct {
	*PartialHandler
	*Damage
	ability *ability
}

// NewDoT returns a DoT handler
func NewDoT(d *Damage, a *ability, duration GameDuration) *DoT {
	return &DoT{
		PartialHandler: NewPartialHandler(d.subject, d.object, duration),
		Damage:         d,
		ability:        a,
	}
}

// OnAttach removes duplicate DoTs
func (d *DoT) OnAttach() {
	d.Object().AddEventHandler(d, EventDead)
	d.Object().AddEventHandler(d, EventGameTick)
	d.Object().AddEventHandler(d, EventXoT)
	ok := d.Container().EveryObjectHandler(d.Object(), func(ha Handler) bool {
		switch ha := ha.(type) {
		case *DoT:
			if ha == d || ha.Subject() != d.Subject() || ha.ability != d.ability {
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
}

// OnDetach removes the EventHandlers
func (d *DoT) OnDetach() {
	d.Object().RemoveEventHandler(d, EventDead)
	d.Object().RemoveEventHandler(d, EventGameTick)
	d.Object().RemoveEventHandler(d, EventXoT)
}

// HandleEvent handles the event
func (d *DoT) HandleEvent(e Event) {
	switch e {
	case EventDead:
		d.Stop(d)
	case EventGameTick:
		if d.IsExpired() {
			d.Up()
		}
	case EventXoT:
		d.Perform()
	}
}

// Up ends the DoT
func (d *DoT) Up() {
	d.Stop(d)
	d.Publish(message{
	// TODO pack message
	})
}
