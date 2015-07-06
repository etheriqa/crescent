package main

import (
	"flag"
	"time"

	"github.com/etheriqa/crescent/class"
	"github.com/etheriqa/crescent/game"
	"github.com/etheriqa/crescent/stage"
)

func main() {
	var addr = flag.String("addr", ":25200", "service address")
	var origin = flag.String("origin", "", "acceptable origin header")
	flag.Parse()
	app := game.App{
		Addr:   *addr,
		Origin: *origin,
		Seed:   time.Now().UnixNano(),
		StageFactory: game.StageFactories{
			1: stage.NewStagePrototype,
		},
		ClassFactory: game.ClassFactories{
			"Tank":     class.NewClassTank,
			"Assassin": class.NewClassAssassin,
			"Disabler": class.NewClassDisabler,
			"Mage":     class.NewClassMage,
			"Healer":   class.NewClassHealer,
		},
	}
	app.Run()
}
