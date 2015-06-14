package main

import (
	"fmt"

	"github.com/Sirupsen/logrus"
)

type Instance struct {
	name map[ClientID]ClientName

	r InstanceInput
	w InstanceOutputWriter
}

// NewInstance returns a Instance
func NewInstance(r InstanceInput, w InstanceOutputWriter) *Instance {
	return &Instance{
		name: make(map[ClientID]ClientName),

		r: r,
		w: w,
	}
}

// Run starts the instance routine
func (i *Instance) Run() {
	for {
		select {
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
			case InputLeave:
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
