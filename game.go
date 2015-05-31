package main

import (
	"time"
)

type game struct {
	inc  chan message
	out  chan message
	cids map[uint64]bool
}

type message struct {
	cid  uint64
	code string
	data map[string]interface{}
}

func newGame(inc chan message, out chan message) *game {
	return &game{
		inc:  inc,
		out:  out,
		cids: make(map[uint64]bool),
	}
}

func (g *game) run() {
	// todo mock
	tick := time.Tick(time.Second)
	for {
		select {
		case m := <-g.inc:
			g.cids[m.cid] = true
		case <-tick:
			for cid := range g.cids {
				g.out <- message{
					cid:  cid,
					code: "test",
					data: map[string]interface{}{
						"server_time": time.Now().String(),
					},
				}
			}
		}
	}
}
