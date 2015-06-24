package crescent

import (
	"github.com/Sirupsen/logrus"
)

var log = logrus.New()

type App struct {
	Addr   string
	Origin string
	StageFactory
	ClassFactory
}

// Logger returns the Logger
func Logger() *logrus.Logger {
	return log
}

// Run the server application
func (a App) Run() {
	log.WithFields(logrus.Fields{
		"addr":   a.Addr,
		"origin": a.Origin,
	}).Info("Start up")
	i := MakeInstanceInput(1024)
	o := MakeInstanceOutput(1024)
	instance := NewInstance(a.StageFactory, a.ClassFactory, i, o)
	network := NewNetwork(a.Origin, o, i)
	go instance.Run()
	network.Run(a.Addr)
}
