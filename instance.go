package main

import (
	"fmt"
	"time"

	"github.com/Sirupsen/logrus"
)

type Instance struct {
	time *InstanceTime
	name map[ClientID]ClientName

	g *Game
	r InstanceInput
	w InstanceOutputWriter
}

// NewInstance returns a Instance
func NewInstance(r InstanceInput, w InstanceOutputWriter) *Instance {
	time := new(InstanceTime)
	return &Instance{
		time: time,
		name: make(map[ClientID]ClientName),

		g: NewGame(time, w),
		r: r,
		w: w,
	}
}

// Run starts the instance routine
func (i *Instance) Run() {
	t := time.Tick(RealGameTick)
	for {
		select {
		case <-t:
			*i.time = i.time.Add(GameTick)
			if i.time.IsRegenerationTick() {
				i.g.PerformRegenerationTick()
			}
			if i.time.IsPeriodicalTick() {
				i.g.PerformPeriodicalTick()
			}
			i.g.PerformGameTick()
		case input, ok := <-i.r:
			if !ok {
				log.Fatal("Cannot read from the input channel")
			}
			id := input.ClientID
			switch input := input.Input.(type) {
			case InputConnect:
				i.connect(id, input)
			case InputDisconnect:
				i.disconnect(id, input)
			case InputChat:
				i.chat(id, input)
			case InputStage:
			case InputJoin:
				i.join(id, input)
			case InputLeave:
				i.leave(id, input)
			case InputAbility:
			case InputInterrupt:
			default:
				log.Fatal("Unknown input type")
			}
		}
	}
}

// connect
func (i *Instance) connect(id ClientID, input InputConnect) {
	i.name[id] = input.ClientName

	i.w.BindClientID(id).Write(OutputMessage{
		Message: "Welcome to the Crescent!",
	})
	i.w.Write(OutputMessage{
		Message: fmt.Sprintf("%s has connected.", input.ClientName),
	})
}

// disconnect
func (i *Instance) disconnect(id ClientID, input InputDisconnect) {
	name := i.name[id]
	delete(i.name, id)

	i.w.Write(OutputMessage{
		Message: fmt.Sprintf("%s has disconnected.", name),
	})
}

// chat
func (i *Instance) chat(id ClientID, input InputChat) {
	name := i.name[id]
	message := input.Message

	i.w.Write(OutputChat{
		ClientName: i.name[id],
		Message:    input.Message,
	})
	log.WithFields(logrus.Fields{
		"type":    "chat",
		"name":    name,
		"message": message,
	}).Infof("%s: %s", name, message)
}

// join
func (i *Instance) join(id ClientID, input InputJoin) {
	// TODO WIP
	class := NewClassHealer()
	u, err := i.g.units.Join(0, UnitName(i.name[id]), class)
	if err != nil {
		return
	}
	i.w.Write(OutputUnitJoin{
		UnitID:       u.ID(),
		UnitGroup:    u.Group(),
		UnitPosition: u.Position(),
		UnitName:     u.Name(),
		ClassName:    u.ClassName(),
		Health:       u.Health(),
		HealthMax:    u.HealthMax(),
		Mana:         u.Mana(),
		ManaMax:      u.ManaMax(),
	})
}

// leave
func (i *Instance) leave(id ClientID, input InputLeave) {
	// TODO WIP
	i.g.units.Leave(0)
	i.w.Write(OutputUnitLeave{
		UnitID: 0,
	})
}
