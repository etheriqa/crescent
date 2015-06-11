package main

type cooldown struct {
	partialOperator
	*ability
}

// onAttach adds the EventHandler
func (c *cooldown) onAttach() {
	c.AddEventHandler(c, EventDead)
	c.AddEventHandler(c, EventGameTick)
	c.publish(message{
		// TODO pack message
		t: outCooldown,
	})
}

// onDetach removes the EventHandler
func (c *cooldown) onDetach() {
	c.RemoveEventHandler(c, EventDead)
	c.RemoveEventHandler(c, EventGameTick)
}

// HandleEvent handles the event
func (c *cooldown) HandleEvent(e Event) {
	switch e {
	case EventDead:
		c.detachOperator(c)
	case EventGameTick:
		c.up()
	}
}

// up ends the cooldown time
func (c *cooldown) up() {
	c.expire(c, message{
		// TODO pack message
		t: outCooldown,
	})
}
