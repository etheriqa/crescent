package main

type cooldown struct {
	partialOperator
	ability ability
}

func (c *cooldown) onAttach() {
	c.unit.addEventHandler(c, eventDead)
	c.unit.addEventHandler(c, eventGameTick)
	c.unit.publish(message{
		// TODO pack message
		t: outCooldown,
	})
}

// onDetach removes the eventHandler
func (c *cooldown) onDetach() {
	c.unit.removeEventHandler(c, eventDead)
	c.unit.removeEventHandler(c, eventGameTick)
}

// handleEvent handles the event
func (c *cooldown) handleEvent(e event) {
	switch e {
	case eventDead:
		c.terminate(c)
	case eventGameTick:
		c.expire(c, message{
			// TODO pack message
			t: outCooldown,
		})
	}
}
