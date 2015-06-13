package main

import (
	"flag"

	"github.com/Sirupsen/logrus"
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
	i := make(chan Input, 1024)
	o := make(chan Output, 1024)
	network := NewNetwork(i, o)
	instance := NewInstance(i, o)
	go network.Run(*addr)
	instance.Run()
}
