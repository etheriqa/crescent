package main

import (
	"flag"

	"github.com/Sirupsen/logrus"

	"github.com/etheriqa/crescent"
)

var addr = flag.String("addr", ":25200", "service address")
var origin = flag.String("origin", "", "acceptable origin header")

func init() {
	flag.Parse()
}

func main() {
	crescent.Logger().WithFields(logrus.Fields{
		"addr":   *addr,
		"origin": *origin,
	}).Info("Start up")
	i := crescent.MakeInstanceInput(1024)
	o := crescent.MakeInstanceOutput(1024)
	instance := crescent.NewInstance(i, o)
	network := crescent.NewNetwork(*origin, o, i)
	go instance.Run()
	network.Run(*addr)
}
