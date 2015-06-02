package main

import (
	"time"
)

type game struct {
	names map[string]interface{}
	inc   chan message
	out   chan message
}

func newGame(inc chan message, out chan message) *game {
	return &game{
		names: make(map[string]interface{}),
		inc:   inc,
		out:   out,
	}
}

func (g *game) run() {
	tick := time.Tick(10 * time.Second)
	for {
		select {
		case m := <-g.inc:
			switch m.t {
			case netRegister:
				g.register(&m)
			case netUnregister:
				g.unregister(&m)
			case incStage:
				g.stage(&m)
			case incSeat:
				g.seat(&m)
			case incLeave:
				g.leave(&m)
			case incAbility:
				g.ability(&m)
			case incInterrupt:
				g.interrupt(&m)
			case incChat:
				g.chat(&m)
			default:
				g.terminate(m.name)
			}
		case <-tick:
			g.tick()
		}
	}
}

func (g *game) register(m *message) {
	g.names[m.name] = nil
	g.out <- message{
		t: outConnect,
		d: map[string]interface{}{
			"name": m.name,
		},
	}
}

func (g *game) unregister(m *message) {
	delete(g.names, m.name)
	g.out <- message{
		t: outDisconnect,
		d: map[string]interface{}{
			"name": m.name,
		},
	}
}

func (g *game) stage(m *message) {
}

func (g *game) seat(m *message) {
}

func (g *game) leave(m *message) {
}

func (g *game) ability(m *message) {
}

func (g *game) interrupt(m *message) {
}

func (g *game) tick() {
	// todo mock
	g.out <- message{
		t: "tick",
		d: map[string]interface{}{
			"server_time": time.Now().String(),
		},
	}
}

func (g *game) chat(m *message) {
	body := m.d["body"].(string)
	g.out <- message{
		t: incChat,
		d: map[string]interface{}{
			"name": m.name,
			"body": body,
		},
	}
	log.Infof("@%s: %s", m.name, body)
}

func (g *game) terminate(name string) {
	g.out <- message{
		name: name,
		t:    gameTerminate,
	}
}
