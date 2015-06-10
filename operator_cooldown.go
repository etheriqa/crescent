package main

type cooldown struct {
	partialOperator
	*ability
}

// onAttach adds the eventHandler
func (c *cooldown) onAttach() {
	c.addEventHandler(c, eventGameTick)
	c.publish(message{
		// TODO pack message
		t: outCooldown,
	})
}

// onDetach removes the eventHandler
func (c *cooldown) onDetach() {
	c.removeEventHandler(c, eventGameTick)
}

// handleEvent handles the event
func (c *cooldown) handleEvent(e event) {
	switch e {
	case eventGameTick:
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
