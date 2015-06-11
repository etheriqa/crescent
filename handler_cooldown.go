package main

type Cooldown struct {
	PartialHandler
	*ability
}

// OnAttach adds the EventHandler
func (c *Cooldown) OnAttach() {
	c.Object().AddEventHandler(c, EventDead)
	c.Object().AddEventHandler(c, EventGameTick)
	c.Publish(message{
	// TODO pack message
	})
}

// OnDetach removes the EventHandler
func (c *Cooldown) OnDetach() {
	c.Object().RemoveEventHandler(c, EventDead)
	c.Object().RemoveEventHandler(c, EventGameTick)
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
