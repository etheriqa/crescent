package main

type Cooldown struct {
	partialHandler
	*ability
}

// OnAttach adds the EventHandler
func (c *Cooldown) OnAttach() {
	c.AddEventHandler(c, EventDead)
	c.AddEventHandler(c, EventGameTick)
	c.publish(message{
		// TODO pack message
		t: outCooldown,
	})
}

// OnDetach removes the EventHandler
func (c *Cooldown) OnDetach() {
	c.RemoveEventHandler(c, EventDead)
	c.RemoveEventHandler(c, EventGameTick)
}

// HandleEvent handles the event
func (c *Cooldown) HandleEvent(e Event) {
	switch e {
	case EventDead:
		c.detachHandler(c)
	case EventGameTick:
		c.up()
	}
}

// up ends the cooldown time
func (c *Cooldown) up() {
	c.expire(c, message{
		// TODO pack message
		t: outCooldown,
	})
}
