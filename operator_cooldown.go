package main

type cooldown struct {
	partialOperator
	*ability
}

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
	c.expire(c, message{
		// TODO pack message
		t: outCooldown,
	})
}

// handleEvent handles the event
func (c *cooldown) handleEvent(e event) {
	switch e {
	case eventGameTick:
		if c.isExpired() {
			c.detachOperator(c)
		}
	}
}
