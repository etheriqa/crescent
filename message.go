package main

const (
	netRegister       = "netRegister"
	netUnregister     = "netUnregister"
	incStage          = "stage"
	incSeat           = "seat"
	incLeave          = "leave"
	incActivate       = "activate"
	incInterrupt      = "interrupt"
	incChat           = "chat"
	outConnect        = "connect"
	outDisconnect     = "disconnect"
	outStage          = "stage"
	outSeat           = "seat"
	outLeave          = "leave"
	outEvent          = "event"
	outChat           = "chat"
	outActivate       = "activate"
	outInterrupt      = "interrupt"
	outModifierAttach = "modifierAttach"
	outModifierDetach = "modifierDetach"
	gameTerminate     = "gameTerminate"
)

type message struct {
	name string
	t    string
	d    map[string]interface{}
}
