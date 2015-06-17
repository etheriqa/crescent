package main

import (
	"fmt"
	"time"

	"github.com/Sirupsen/logrus"
)

type Instance struct {
	time *InstanceTime
	name map[ClientID]ClientName
	uid  map[ClientID]UnitID

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
		uid:  make(map[ClientID]UnitID),

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
			cid := input.ClientID
			switch input := input.Input.(type) {
			case InputConnect:
				i.connect(cid, input)
			case InputDisconnect:
				i.disconnect(cid, input)
			case InputChat:
				i.chat(cid, input)
			case InputStage:
				// TODO WIP
			case InputJoin:
				i.join(cid, input)
			case InputLeave:
				i.leave(cid, input)
			case InputAbility:
				i.ability(cid, input)
			case InputInterrupt:
				// TODO WIP
			default:
				log.Fatal("Unknown input type")
			}
		}
	}
}

// connect
func (i *Instance) connect(cid ClientID, input InputConnect) {
	i.name[cid] = input.ClientName

	i.w.BindClientID(cid).Write(OutputMessage{
		Message: "Welcome to Crescent!",
	})
	i.w.Write(OutputMessage{
		Message: fmt.Sprintf("%s has connected.", input.ClientName),
	})
}

// disconnect
func (i *Instance) disconnect(cid ClientID, input InputDisconnect) {
	i.leave(cid, InputLeave{})

	name := i.name[cid]
	delete(i.name, cid)

	i.w.Write(OutputMessage{
		Message: fmt.Sprintf("%s has disconnected.", name),
	})
}

// chat
func (i *Instance) chat(cid ClientID, input InputChat) {
	name := i.name[cid]
	message := input.Message

	i.w.Write(OutputChat{
		ClientName: i.name[cid],
		Message:    input.Message,
	})
	log.WithFields(logrus.Fields{
		"type":    "chat",
		"name":    name,
		"message": message,
	}).Infof("%s: %s", name, message)
}

// join
func (i *Instance) join(cid ClientID, input InputJoin) {
	// TODO disable join when a game is in progress
	if _, ok := i.uid[cid]; ok {
		return
	}
	// TODO refactor: make ClassFactory
	var class *Class
	switch input.ClassName {
	case "Assassin":
		class = NewClassAssassin()
	case "Disabler":
		class = NewClassDisabler()
	case "Healer":
		class = NewClassHealer()
	case "Mage":
		class = NewClassMage()
	case "Tank":
		class = NewClassTank()
	default:
		return
	}
	uid, err := i.g.Join(UnitGroupPlayer, UnitName(i.name[cid]), class)
	if err != nil {
		log.WithFields(logrus.Fields{
			"cid":  cid,
			"type": "join",
		}).Warn(err)
		return
	}
	i.uid[cid] = uid
}

// leave
func (i *Instance) leave(cid ClientID, input InputLeave) {
	if _, ok := i.uid[cid]; !ok {
		return
	}
	if err := i.g.Leave(i.uid[cid]); err != nil {
		log.WithFields(logrus.Fields{
			"cid":  cid,
			"type": "leave",
		}).Warn(err)
		return
	}
	delete(i.uid, cid)
}

// ability
func (i *Instance) ability(cid ClientID, input InputAbility) {
	if _, ok := i.uid[cid]; !ok {
		return
	}
	i.g.Ability(i.uid[cid], input.ObjectUnitID, input.AbilityName)
}
