package main

type cooldown struct {
	partialOperator
}

func (c *cooldown) onAttach() {
	c.unit.addEventHandler(c, eventGameTick)
	c.unit.publish(message{
		// todo pack message
		t: outCooldown,
	})
}

// onDetach removes the eventHandler
func (c *cooldown) onDetach() {
	c.unit.removeEventHandler(c, eventGameTick)
}

// handleEvent handles the event
func (c *cooldown) handleEvent(e event) {
	switch e {
	case eventGameTick:
		c.expire(c, message{
			// todo pack message
			t: outCooldown,
		})
	}
}
