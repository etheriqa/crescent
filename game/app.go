package game

import (
	"math/rand"

	"github.com/Sirupsen/logrus"
)

var log = logrus.New()

type App struct {
	Addr   string
	Origin string
	Seed   int64
	LevelFactory
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
		"seed":   a.Seed,
	}).Info("Start up")
	i := MakeInstanceInput(1024)
	o := MakeInstanceOutput(1024)
	rand := rand.New(rand.NewSource(a.Seed))
	instance := NewInstance(a.LevelFactory, a.ClassFactory, rand, i, o)
	network := NewNetwork(a.Origin, o, i)
	go instance.Run()
	network.Run(a.Addr)
}
