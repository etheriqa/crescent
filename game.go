package main

import (
	"time"
)

type gameTime int64
type gameDuration int64

const (
	millisecond = gameDuration(time.Millisecond / gameTick)
	second      = gameDuration(time.Second / gameTick)
)

type game struct {
	time  gameTime
	names map[string]*unit
	uid   unitID
	uids  map[unitID]*unit
	seats map[uint8]*unit
	inc   chan message
	out   chan message
}

func newGame(inc chan message, out chan message) *game {
	return &game{
		time:  0,
		names: make(map[string]*unit),
		uid:   0,
		uids:  make(map[unitID]*unit),
		seats: make(map[uint8]*unit),
		inc:   inc,
		out:   out,
	}
}

// now returns the game time
func (g *game) now() gameTime {
	return g.time
}

// after returns the game time after the duration
func (g *game) after(d gameDuration) gameTime {
	return g.time + gameTime(d)
}

// friends returns given unit's friend units **including itself**
func (g *game) friends(u *unit) (us []*unit) {
	for _, unit := range g.uids {
		if unit.group == u.group {
			us = append(us, unit)
		}
	}
	return
}

// enemies returns given unit's enemy units
func (g *game) enemies(u *unit) (us []*unit) {
	for _, unit := range g.uids {
		if unit.group != u.group {
			us = append(us, unit)
		}
	}
	return
}

func (g *game) publish(m message) {
	g.out <- m
}

func (g *game) nextUnitID() unitID {
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
		u.gameTick()
	}
	if int64(g.time)*int64(gameTick)%int64(xotTick) != 0 {
		return
	}
	for _, u := range g.uids {
		u.xotTick()
	}
}

func (g *game) register(m *message) {
	g.names[m.name] = nil
	g.publish(message{
		t: outConnect,
		d: map[string]interface{}{
			"name": m.name,
		},
	})
}

func (g *game) unregister(m *message) {
	if u, _ := g.names[m.name]; u != nil {
		delete(g.uids, u.id)
		delete(g.seats, u.seat)
	}
	delete(g.names, m.name)
	g.publish(message{
		t: outDisconnect,
		d: map[string]interface{}{
			"name": m.name,
		},
	})
}

func (g *game) stage(m *message) {
	g.uids = make(map[unitID]*unit)
	g.seats = make(map[uint8]*unit)
	g.publish(message{
	// TODO pack message
	})
}

func (g *game) seat(m *message) {
	seat := uint8(m.d["seat"].(float64))
	unitName := m.d["unit"].(string)
	if _, ok := g.seats[seat]; ok {
		return
	}
	var c *class
	switch unitName {
	case "Assassin":
		c = newClassAssassin()
	case "Disabler":
		c = newClassDisabler()
	case "Healer":
		c = newClassHealer()
	case "Mage":
		c = newClassMage()
	case "Tank":
		c = newClassTank()
	default:
		return
	}
	u := newUnit(g, c)
	g.names[m.name] = u
	g.uids[u.id] = u
	g.seats[seat] = u
	g.publish(message{
		t: outSeat,
		d: map[string]interface{}{
			"name": m.name,
			"uid":  u.id,
			"seat": seat,
		},
	})
}

func (g *game) leave(m *message) {
	u, ok := g.names[m.name]
	if !ok || u == nil {
		return
	}
	delete(g.uids, u.id)
	delete(g.seats, u.seat)
	g.names[m.name] = nil
	g.publish(message{
		t: outLeave,
		d: map[string]interface{}{
			"seat": u.seat,
		},
	})
}

func (g *game) activate(m *message) {
	key := m.d["key"].(string)
	unit := g.names[m.name]
	var a *ability
	switch key {
	case "q":
		a = unit.class.abilities[0]
	case "w":
		a = unit.class.abilities[1]
	case "e":
		a = unit.class.abilities[2]
	case "r":
		a = unit.class.abilities[3]
	default:
		return
	}
	target := g.uids[m.d["uid"].(unitID)]
	unit.attachHandler(NewActivating(unit, target, a))
}

func (g *game) interrupt(m *message) {
	unit := g.names[m.name]
	for o := range unit.handlers {
		switch o.(type) {
		case *Activating:
			unit.detachHandler(o)
		}
	}
}

func (g *game) chat(m *message) {
	body := m.d["body"].(string)
	g.publish(message{
		t: outChat,
		d: map[string]interface{}{
			"name": m.name,
			"body": body,
		},
	})
	log.Infof("@%s: %s", m.name, body)
}

func (g *game) terminate(name string) {
	g.publish(message{
		name: name,
		t:    gameTerminate,
	})
}
