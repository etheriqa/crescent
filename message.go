package main

const (
	netRegister   = "netRegister"
	netUnregister = "netUnregister"
	incChat       = "chat"
	outConnect    = "connect"
	outDisconnect = "disconnect"
	outChat       = "chat"
	gameTerminate = "terminate"
)

type message struct {
	cid uint64
	t   string
	d   map[string]interface{}
}
