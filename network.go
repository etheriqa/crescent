package main

import (
	"net/http"
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/websocket"
)

type ClientID uint64
type ClientName string

type Network struct {
	mu   *sync.RWMutex
	seq  ClientID
	id   map[ClientID]*Client
	name map[ClientName]bool

	ichan chan Input
	ochan chan Output
}

type Client struct {
	id   ClientID
	name ClientName
	ws   *websocket.Conn
	buf  chan []byte
}

// NewNetwork returns a Network
func NewNetwork(i chan Input, o chan Output) *Network {
	return &Network{
		mu:   new(sync.RWMutex),
		seq:  0,
		id:   make(map[ClientID]*Client),
		name: make(map[ClientName]bool),

		ichan: i,
		ochan: o,
	}
}

// Run starts the network routine
func (n *Network) Run(addr string) {
	go n.dispatcher()
	http.HandleFunc("/", n.wsHandler)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.WithField("err", err).Fatal("Failed http.ListenAndServe()")
	}
}

// sync calls the callback with writing lock
func (n *Network) sync(callback func()) {
	n.mu.Lock()
	callback()
	n.mu.Unlock()
}

// sync calls the callback with reading lock
func (n *Network) syncR(callback func()) {
	n.mu.RLock()
	callback()
	n.mu.RUnlock()
}

// wsHandler handles a HTTP request
func (n *Network) wsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method Not Allowed", 405)
		return
	}
	name := ClientName(r.FormValue("name"))
	if name == "" {
		http.Error(w, "Bad Request", 400)
		return
	}
	var c *Client
	var ok bool
	n.sync(func() {
		if n.name[name] {
			http.Error(w, "Bad Request", 400)
			return
		}
		upgrader := websocket.Upgrader{
			CheckOrigin: func(*http.Request) bool { return true }, // FIXME
		}
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.WithField("err", err).Warn("Failed websocket.Upgrade()")
			http.Error(w, "Internal Server Error", 500)
			return
		}
		n.seq++
		c = NewClient(n.seq, name, ws)
		ok = true
	})
	if !ok {
		return
	}
	n.sync(func() { n.register(c) })
	go n.receiver(c)
	n.sender(c)
}

// register registers the Client
func (n *Network) register(c *Client) {
	n.id[c.id] = c
	n.name[c.name] = true
	n.ichan <- Input{
		ClientID: c.id,
		Input: InputConnect{
			ClientName: c.name,
		},
	}
	log.WithFields(logrus.Fields{
		"id":   c.id,
		"name": c.name,
		"addr": c.ws.RemoteAddr(),
	}).Info("Client has been registered")
}

// unregister unregisters the Client
func (n *Network) unregister(c *Client) {
	if _, ok := n.id[c.id]; !ok {
		return
	}
	delete(n.id, c.id)
	delete(n.name, c.name)
	close(c.buf)
	c.ws.Close()
	n.ichan <- Input{
		ClientID: c.id,
		Input:    InputDisconnect{},
	}
	log.WithFields(logrus.Fields{
		"id":   c.id,
		"name": c.name,
		"addr": c.ws.RemoteAddr(),
	}).Info("Client has been unregistered")
}

// receiver reads frames from the WebSocket connection and writes messages to the Instance
func (n *Network) receiver(c *Client) {
	defer n.sync(func() { n.unregister(c) })
	for {
		_, p, err := c.ws.ReadMessage()
		if err != nil {
			log.WithFields(logrus.Fields{
				"id":   c.id,
				"name": c.name,
				"err":  err,
			}).Warn("Failed websocket.ReadMessage()")
			return
		}
		in, err := DecodeInputFrame(p)
		if err != nil {
			log.WithFields(logrus.Fields{
				"id":   c.id,
				"name": c.name,
				"p":    string(p),
				"err":  err,
			}).Warn("Failed DecodeInputFrame()")
			return
		}
		n.ichan <- Input{
			ClientID: c.id,
			Input:    in,
		}
	}
}

// sender reads frames from the write buffer and writes frames to the WebSocket connection
func (n *Network) sender(c *Client) {
	defer n.sync(func() { n.unregister(c) })
	for {
		select {
		case p, ok := <-c.buf:
			if !ok {
				err := c.ws.WriteMessage(websocket.CloseMessage, []byte{})
				if err != nil {
					log.WithFields(logrus.Fields{
						"id":   c.id,
						"name": c.name,
						"err":  err,
					}).Warn("Failed websocket.WriteMessage()")
				}
				return
			}
			err := c.ws.WriteMessage(websocket.TextMessage, p)
			if err != nil {
				log.WithFields(logrus.Fields{
					"id":   c.id,
					"name": c.name,
					"err":  err,
				}).Warn("Failed websocket.WriteMessage()")
				return
			}
		}
	}
}

// dispatcher reads messages from the Instance and writes frames to the write buffer
func (n *Network) dispatcher() {
	for {
		select {
		case out, ok := <-n.ochan:
			if !ok {
				log.Fatal("Cannot read from the output channel")
			}
			p, err := EncodeOutputFrame(out.Output)
			if err != nil {
				log.WithFields(logrus.Fields{
					"Output": out.Output,
					"err":    err,
				}).Fatal("Failed EncodeOutputFrame")
			}
			if c, ok := n.id[out.ClientID]; ok {
				c.buf <- p
			} else {
				n.syncR(func() {
					for _, c := range n.id {
						c.buf <- p
					}
				})
			}
		}
	}
}

// NewClient returns a Client
func NewClient(id ClientID, name ClientName, ws *websocket.Conn) *Client {
	return &Client{
		id:   id,
		name: name,
		ws:   ws,
		buf:  make(chan []byte, 1024),
	}
}
