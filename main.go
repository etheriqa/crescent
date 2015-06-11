package main

import (
	"flag"
	"time"

	"github.com/Sirupsen/logrus"
)

const (
	GameTick   = time.Millisecond
	TickerTick = time.Second
)

var addr = flag.String("addr", ":25200", "service address")
var debug = flag.Bool("debug", false, "debug mode")
var log = logrus.New()

func init() {
	flag.Parse()
	if *debug {
		log.Level = logrus.DebugLevel
	}
}

func main() {
	log.WithFields(logrus.Fields{
		"addr":  *addr,
		"debug": *debug,
	}).Info("Start up")
	i := make(chan message)
	o := make(chan message)
	n := NewNetwork(i, o)
	g := NewGame(i, o)
	go n.Run(*addr)
	go g.Run()
	t := time.Tick(GameTick)
	for {
		select {
		case <-t:
			i <- message{
				t: "sysTick",
			}
		}
	}
}
