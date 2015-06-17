package main

import (
	"flag"

	"github.com/Sirupsen/logrus"
)

var addr = flag.String("addr", ":25200", "service address")
var origin = flag.String("origin", "", "acceptable origin header")
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
		"addr":   *addr,
		"origin": *origin,
		"debug":  *debug,
	}).Info("Start up")
	i := MakeInstanceInput(1024)
	o := MakeInstanceOutput(1024)
	instance := NewInstance(i, o)
	network := NewNetwork(*origin, o, i)
	go instance.Run()
	network.Run(*addr)
}
