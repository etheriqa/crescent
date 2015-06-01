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
	tick := time.Tick(time.Second)
	for {
		select {
		case m := <-g.inc:
			switch m.t {
			case netRegister:
				g.register(&m)
			case netUnregister:
				g.unregister(&m)
			default:
				g.terminate(&m)
			}
		case <-tick:
			for cid := range g.names {
				g.out <- message{
					cid: cid,
					t:   "test",
					d: map[string]interface{}{
						"server_time": time.Now().String(),
					},
				}
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
	name := g.names[m.cid]
	delete(g.names, m.cid)
	g.out <- message{
		t: outDisconnect,
		d: map[string]interface{}{
			"name": name,
		},
	}
}

func (g *game) terminate(m *message) {
	g.out <- message{
		cid: m.cid,
		t:   gameTerminate,
	}
}
