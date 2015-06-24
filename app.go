package crescent

import (
	"github.com/Sirupsen/logrus"
)

type App struct {
	Addr   string
	Origin string
}

var log = logrus.New()

// Run the server application
func (a App) Run() {
	log.WithFields(logrus.Fields{
		"addr":   a.Addr,
		"origin": a.Origin,
	}).Info("Start up")
	i := MakeInstanceInput(1024)
	o := MakeInstanceOutput(1024)
	instance := NewInstance(i, o)
	network := NewNetwork(a.Origin, o, i)
	go instance.Run()
	network.Run(a.Addr)
}
