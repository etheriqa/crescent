package main

import (
	"time"
)

type GameTime int64
type GameDuration int64

const (
	Millisecond = GameDuration(time.Millisecond / GameTick)
	Second      = GameDuration(time.Second / GameTick)
)

type GameClock interface {
	Now() GameTime
	After(GameDuration) GameTime
}

type Game struct {
	HandlerContainer
	time  GameTime
	names map[string]*Unit
	uid   unitID
	uids  map[unitID]*Unit
	seats map[uint8]*Unit
	inc   chan message
	out   chan message
}

func NewGame(inc chan message, out chan message) *Game {
	return &Game{
		HandlerContainer: NewHandlerContainer(),
		time:             0,
		names:            make(map[string]*Unit),
		uid:              0,
		uids:             make(map[unitID]*Unit),
		seats:            make(map[uint8]*Unit),
		inc:              inc,
		out:              out,
	}
}

// Now returns the game time
func (g *Game) Now() GameTime {
	return g.time
}

// After returns the game time after the game duration
func (g *Game) After(d GameDuration) GameTime {
	return g.time + GameTime(d)
}

// Friends returns given unit's friend units **including itself**
func (g *Game) Friends(u *Unit) (us []*Unit) {
	for _, unit := range g.uids {
		if unit.group == u.group {
			us = append(us, unit)
		}
	}
	return
}

// Enemies returns given unit's enemy units
func (g *Game) Enemies(u *Unit) (us []*Unit) {
	for _, unit := range g.uids {
		if unit.group != u.group {
			us = append(us, unit)
		}
	}
	return
}

func (g *Game) Publish(m message) {
	g.out <- m
}

func (g *Game) nextUnitID() unitID {
	g.uid++
	return g.uid
}

func (g *Game) Run() {
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

func (g *Game) tick(m *message) {
	g.time++
	for _, u := range g.uids {
		u.GameTick()
	}
	if int64(g.time)*int64(GameTick)%int64(TickerTick) != 0 {
		return
	}
	for _, u := range g.uids {
		u.TickerTick()
	}
}

func (g *Game) register(m *message) {
	g.names[m.name] = nil
	g.Publish(message{
		t: outConnect,
		d: map[string]interface{}{
			"name": m.name,
		},
	})
}

func (g *Game) unregister(m *message) {
	if u, _ := g.names[m.name]; u != nil {
		delete(g.uids, u.id)
		delete(g.seats, u.seat)
	}
	delete(g.names, m.name)
	g.Publish(message{
		t: outDisconnect,
		d: map[string]interface{}{
			"name": m.name,
		},
	})
}

func (g *Game) stage(m *message) {
	g.uids = make(map[unitID]*Unit)
	g.seats = make(map[uint8]*Unit)
	g.Publish(message{
	// TODO pack message
	})
}

func (g *Game) seat(m *message) {
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
	u := NewUnit(g, c)
	g.names[m.name] = u
	g.uids[u.id] = u
	g.seats[seat] = u
	g.Publish(message{
		t: outSeat,
		d: map[string]interface{}{
			"name": m.name,
			"uid":  u.id,
			"seat": seat,
		},
	})
}

func (g *Game) leave(m *message) {
	u, ok := g.names[m.name]
	if !ok || u == nil {
		return
	}
	delete(g.uids, u.id)
	delete(g.seats, u.seat)
	g.names[m.name] = nil
	g.Publish(message{
		t: outLeave,
		d: map[string]interface{}{
			"seat": u.seat,
		},
	})
}

func (g *Game) activate(m *message) {
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
	g.AttachHandler(NewActivating(MakeUnitPair(unit, target), a))
}

func (g *Game) interrupt(m *message) {
	unit := g.names[m.name]
	g.ForSubjectHandler(unit, func(ha Handler) {
		switch ha.(type) {
		case *Activating:
			g.DetachHandler(ha)
		}
	})
}

func (g *Game) chat(m *message) {
	body := m.d["body"].(string)
	g.Publish(message{
		t: outChat,
		d: map[string]interface{}{
			"name": m.name,
			"body": body,
		},
	})
	log.Infof("@%s: %s", m.name, body)
}

func (g *Game) terminate(name string) {
	g.Publish(message{
		name: name,
		t:    gameTerminate,
	})
}
