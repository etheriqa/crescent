package main

const (
	netRegister   = "netRegister"
	netUnregister = "netUnregister"
	outConnect    = "connect"
	outDisconnect = "disconnect"
	gameTerminate = "terminate"
)

type message struct {
	cid uint64
	t   string
	d   map[string]interface{}
}
