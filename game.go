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
		inc:   inc,
		out:   out,
		names: make(map[uint64]string),
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
			case incChat:
				g.chat(&m)
			default:
				g.terminate(&m)
			}
		case <-tick:
			g.out <- message{
				t: "test",
				d: map[string]interface{}{
					"server_time": time.Now().String(),
				},
			}
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

func (g *game) chat(m *message) {
	name := g.names[m.cid] // xxx
	g.out <- message{
		t: incChat,
		d: map[string]interface{}{
			"name":    name,
			"message": m.d["message"].(string),
		},
	}
}

func (g *game) terminate(m *message) {
	g.out <- message{
		cid: m.cid,
		t:   gameTerminate,
	}
}
