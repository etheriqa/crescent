package main

import (
	"flag"

	"github.com/etheriqa/crescent"
)

func main() {
	var addr = flag.String("addr", ":25200", "service address")
	var origin = flag.String("origin", "", "acceptable origin header")
	flag.Parse()
	app := crescent.App{
		Addr:   *addr,
		Origin: *origin,
	}
	app.Run()
}
