package main

type gameTime int64

type game struct {
	time  gameTime
	names map[string]*unit
	uid   uidType
	uids  map[uidType]*unit
	seats map[uint8]*unit
	inc   chan message
	out   chan message
}

func newGame(inc chan message, out chan message) *game {
	return &game{
		time:  0,
		names: make(map[string]*unit),
		uid:   0,
		uids:  make(map[uidType]*unit),
		seats: make(map[uint8]*unit),
		inc:   inc,
		out:   out,
	}
}

func (g *game) nextUID() uidType {
	g.uid++
	return g.uid
}

func (g *game) run() {
	for {
		select {
		case m := <-g.inc:
			switch m.t {
			case sysTick:
				g.tick(&m)
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
			case incActivate:
				g.activate(&m)
			case incInterrupt:
				g.interrupt(&m)
			case incChat:
				g.chat(&m)
			default:
				g.terminate(m.name)
			}
		}
	}
}

func (g *game) tick(m *message) {
	g.time++
	for _, u := range g.uids {
		u.gameTick(g.out)
	}
	if int64(g.time)*int64(gameTick)%int64(statsTick) != 0 {
		return
	}
	for _, u := range g.uids {
		u.statsTick(g.out)
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
	if u, _ := g.names[m.name]; u != nil {
		delete(g.uids, u.id)
		delete(g.seats, u.seat)
	}
	delete(g.names, m.name)
	g.out <- message{
		t: outDisconnect,
		d: map[string]interface{}{
			"name": m.name,
		},
	}
}

func (g *game) stage(m *message) {
	// todo change stage
}

func (g *game) seat(m *message) {
	seat := uint8(m.d["seat"].(float64))
	unitName := m.d["unit"].(string)
	if _, ok := g.seats[seat]; ok {
		return
	}
	var u *unit
	switch unitName {
	// todo
	default:
		g.terminate(m.name)
		return
	}
	g.names[m.name] = u
	g.uids[u.id] = u
	g.seats[seat] = u
	g.out <- message{
		t: outSeat,
		d: map[string]interface{}{
			"name": m.name,
			"uid":  u.id,
			"seat": seat,
		},
	}
}

func (g *game) leave(m *message) {
	u, ok := g.names[m.name]
	if !ok || u == nil {
		return
	}
	delete(g.uids, u.id)
	delete(g.seats, u.seat)
	g.names[m.name] = nil
	g.out <- message{
		t: outLeave,
		d: map[string]interface{}{
			"seat": u.seat,
		},
	}
}

func (g *game) activate(m *message) {
}

func (g *game) interrupt(m *message) {
}

func (g *game) chat(m *message) {
	body := m.d["body"].(string)
	g.out <- message{
		t: outChat,
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
