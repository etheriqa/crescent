package main

type Cooldown struct {
	PartialHandler
	*ability
}

// NewCooldown returns a Cooldown handler
func NewCooldown(subject *Unit, ability *ability) *Cooldown {
	return &Cooldown{
		PartialHandler: MakePartialHandler(MakeSubject(subject), ability.cooldownDuration),
		ability:        ability,
	}
}

// OnAttach adds the EventHandler
func (c *Cooldown) OnAttach() {
	c.Subject().AddEventHandler(c, EventDead)
	c.Subject().AddEventHandler(c, EventGameTick)
	c.Publish(message{
	// TODO pack message
	})
}

// OnDetach removes the EventHandler
func (c *Cooldown) OnDetach() {
	c.Subject().RemoveEventHandler(c, EventDead)
	c.Subject().RemoveEventHandler(c, EventGameTick)
}

// HandleEvent handles the Event
func (c *Cooldown) HandleEvent(e Event) {
	switch e {
	case EventDead:
		c.Stop(c)
	case EventGameTick:
		if c.IsExpired() {
			c.Up()
		}
	}
}

// Up ends the cooldown time
func (c *Cooldown) Up() {
	c.Stop(c)
	c.Publish(message{
	// TODO pack message
	})
}
