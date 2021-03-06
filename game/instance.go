package game

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/Sirupsen/logrus"
)

type UserName string

type Instance struct {
	time *InstanceTime
	name map[ClientID]UserName
	uid  map[ClientID]UnitID

	world *World

	sf LevelFactory
	cf ClassFactory

	rand *rand.Rand
	r    InstanceInput
	w    InstanceOutputWriter
}

// NewInstance returns a Instance
func NewInstance(sf LevelFactory, cf ClassFactory, rand *rand.Rand, r InstanceInput, w InstanceOutputWriter) *Instance {
	time := new(InstanceTime)
	return &Instance{
		time: time,
		name: make(map[ClientID]UserName),
		uid:  make(map[ClientID]UnitID),

		world: NewWorld(rand, time, w, sf.New(1)), // TODO remove magic number

		cf: cf,
		sf: sf,

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
				i.sync(i.w)
				i.world.PerformRegenerationTick()
			}
			if i.time.IsPeriodicalTick() {
				i.world.PerformPeriodicalTick()
			}
			i.world.PerformGameTick()
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
			case InputProfile:
				i.profile(cid, input)
			case InputChat:
				i.chat(cid, input)
			case InputLevel:
				i.level(cid, input)
			case InputJoin:
				i.join(cid, input)
			case InputLeave:
				i.leave(cid, input)
			case InputAbility:
				i.ability(cid, input)
			case InputInterrupt:
				i.interrupt(cid, input)
			default:
				log.Fatal("Unknown input type")
			}
		}
	}
}

// generateName returns a random user name
func (i *Instance) generateName() UserName {
	return UserName(fmt.Sprintf("user%03d", rand.Intn(1000)))
}

// sync sends a OutputSync message
func (i *Instance) sync(w InstanceOutputWriter) {
	w.Write(OutputSync{
		InstanceTime: *i.time,
	})
}

// connect
func (i *Instance) connect(cid ClientID, input InputConnect) {
	name := i.generateName()
	i.name[cid] = name

	i.sync(i.w.BindClientID(cid))
	i.w.BindClientID(cid).Write(OutputMessage{
		Message: "Welcome to Crescent!",
	})
	i.w.BindClientID(cid).Write(OutputMessage{
		Message: "/profile <name> : change your name",
	})
	i.w.BindClientID(cid).Write(OutputMessage{
		Message: "/level 1 : reset level",
	})
	i.w.Write(OutputMessage{
		Message: fmt.Sprintf("%s has joined.", name),
	})
	i.world.SyncWorld(i.w.BindClientID(cid))
}

// disconnect
func (i *Instance) disconnect(cid ClientID, input InputDisconnect) {
	i.leave(cid, InputLeave{})

	name := i.name[cid]
	delete(i.name, cid)

	i.w.Write(OutputMessage{
		Message: fmt.Sprintf("%s has left.", name),
	})
}

// profile
func (i *Instance) profile(cid ClientID, input InputProfile) {
	// TODO validation
	before := i.name[cid]
	after := input.UserName

	i.name[cid] = after

	i.w.Write(OutputMessage{
		Message: fmt.Sprintf("%s has changed the name to %s.", before, after),
	})
}

// chat
func (i *Instance) chat(cid ClientID, input InputChat) {
	name := i.name[cid]
	message := input.Message

	i.w.Write(OutputChat{
		UserName: i.name[cid],
		Message:  input.Message,
	})
	log.WithFields(logrus.Fields{
		"type":    "chat",
		"name":    name,
		"message": message,
	}).Infof("%s: %s", name, message)
}

// level
func (i *Instance) level(cid ClientID, input InputLevel) {
	s := i.sf.New(input.LevelID)
	if s == nil {
		return
	}
	i.w.Write(OutputLevel{})
	i.uid = make(map[ClientID]UnitID)
	i.world = NewWorld(i.rand, i.time, i.w, s)
}

// join
func (i *Instance) join(cid ClientID, input InputJoin) {
	// TODO disable join when a game is in progress
	if _, ok := i.uid[cid]; ok {
		return
	}
	c := i.cf.New(input.ClassName)
	if c == nil {
		return
	}
	uid, err := i.world.Join(UnitGroupPlayer, UnitName(i.name[cid]), c)
	if err != nil {
		log.WithFields(logrus.Fields{
			"cid":  cid,
			"type": "join",
		}).Warn(err)
		return
	}
	i.uid[cid] = uid
	i.world.SyncUnit(i.w.BindClientID(cid), uid)
}

// leave
func (i *Instance) leave(cid ClientID, input InputLeave) {
	if _, ok := i.uid[cid]; !ok {
		return
	}
	if err := i.world.Leave(i.uid[cid]); err != nil {
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
	i.world.Ability(i.uid[cid], input.ObjectUnitID, input.AbilityName)
}

// interrupt
func (i *Instance) interrupt(cid ClientID, input InputInterrupt) {
	if _, ok := i.uid[cid]; !ok {
		return
	}
	i.world.Interrupt(i.uid[cid])
}
