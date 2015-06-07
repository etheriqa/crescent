package main

type cooldown struct {
	partialOperator
	ability ability
}

func (c *cooldown) onAttach() {
	c.addEventHandler(c, eventDead)
	c.addEventHandler(c, eventGameTick)
	c.publish(message{
		// TODO pack message
		t: outCooldown,
	})
}

// onDetach removes the eventHandler
func (c *cooldown) onDetach() {
	c.removeEventHandler(c, eventDead)
	c.removeEventHandler(c, eventGameTick)
}

// handleEvent handles the event
func (c *cooldown) handleEvent(e event) {
	switch e {
	case eventDead:
		c.detachOperator(c)
	case eventGameTick:
		c.expire(c, message{
			// TODO pack message
			t: outCooldown,
		})
	}
}
