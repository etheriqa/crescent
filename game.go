package main

import (
	"time"
)

type game struct {
	names map[uint64]string
	inc   chan message
	out   chan message
}

func newGame(inc chan message, out chan message) *game {
	return &game{
		names: make(map[uint64]string),
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
				g.terminate(&m)
			}
		case <-tick:
			g.tick()
		}
	}
}

func (g *game) register(m *message) {
	name := m.d["name"].(string)
	g.names[m.cid] = name
	g.out <- message{
		t: outConnect,
		d: map[string]interface{}{
			"name": name,
		},
	}
}

func (g *game) unregister(m *message) {
	name := g.names[m.cid] // xxx
	delete(g.names, m.cid)
	g.out <- message{
		t: outDisconnect,
		d: map[string]interface{}{
			"name": name,
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
	name := g.names[m.cid] // xxx
	body := m.d["body"].(string)
	g.out <- message{
		t: incChat,
		d: map[string]interface{}{
			"name": name,
			"body": body,
		},
	}
	log.Infof("@%s: %s", name, body)
}

func (g *game) terminate(m *message) {
	g.out <- message{
		cid: m.cid,
		t:   gameTerminate,
	}
}
