package main

const (
	netRegister   = "netRegister"
	netUnregister = "netUnregister"
	incStage      = "stage"
	incSeat       = "seat"
	incLeave      = "leave"
	incAbility    = "ability"
	incInterrupt  = "interrupt"
	incChat       = "chat"
	outConnect    = "connect"
	outDisconnect = "disconnect"
	outStage      = "stage"
	outSeat       = "seat"
	outLeave      = "leave"
	outEvent      = "event"
	outChat       = "chat"
	gameTerminate = "gameTerminate"
)

type message struct {
	name string
	t    string
	d    map[string]interface{}
}
