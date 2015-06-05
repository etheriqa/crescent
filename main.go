package main

import (
	"flag"
	"time"

	"github.com/Sirupsen/logrus"
)

const (
	gameTick  = time.Millisecond
	statsTick = time.Second
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
	log.Info("Start up")
	i := make(chan message)
	o := make(chan message)
	n := newNetwork(i, o)
	g := newGame(i, o)
	go n.run(*addr)
	go g.run()
	t := time.Tick(gameTick)
	for {
		select {
		case <-t:
			i <- message{
				t: "sysTick",
			}
		}
	}
}
