package main

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/websocket"
)

var wsUpgrader = websocket.Upgrader{
	CheckOrigin: func(*http.Request) bool { return true }, // fixme
}

type network struct {
	rw    *sync.RWMutex
	cid   uint64
	cs    map[uint64]*connection
	names map[string]bool
	inc   chan message
	out   chan message
}

type connection struct {
	id   uint64
	name string
	ws   *websocket.Conn
	buf  chan []byte
}

type frame struct {
	Type string           `json:"type"`
	Data *json.RawMessage `json:"data"`
}

// newNetwork initializes a network
func newNetwork(inc chan message, out chan message) *network {
	return &network{
		rw:    new(sync.RWMutex),
		cid:   0,
		cs:    make(map[uint64]*connection),
		names: make(map[string]bool),
		inc:   inc,
		out:   out,
	}
}

// nextCID generates a connection ID
func (n *network) nextCID() uint64 {
	n.cid++
	return n.cid
}

// run executes the network routine
func (n *network) run(addr string) {
	go n.dispatcher()
	http.HandleFunc("/", n.wsHandler)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.WithField("err", err).Fatal("Failed http.ListenAndServe()")
	}
}

// wsHandler handles a WebSocket connection
func (n *network) wsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method Not Allowed", 405)
		return
	}
	name := r.FormValue("name")
	n.rw.Lock()
	if _, ok := n.names[name]; name == "" || ok {
		http.Error(w, "Bad Request", 400)
		n.rw.Unlock()
		return
	}
	ws, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.WithField("err", err).Warn("Failed websocket.Upgrade()")
		http.Error(w, "Internal Server Error", 500)
		n.rw.Unlock()
		return
	}
	c := newConnection(n.nextCID(), name, ws)
	n.register(c)
	n.rw.Unlock()
	go n.receiver(c)
	n.sender(c)
}

// register registers a connection
func (n *network) register(c *connection) {
	n.cs[c.id] = c
	n.names[c.name] = true
	n.inc <- message{
		cid: c.id,
		t:   netRegister,
		d: map[string]interface{}{
			"name": c.name,
		},
	}
	log.WithFields(logrus.Fields{
		"cid":  c.id,
		"name": c.name,
		"addr": c.ws.RemoteAddr(),
	}).Info("Client has been registered")
}

// unregister closes the connection and unregisters it
func (n *network) unregister(c *connection) {
	if _, ok := n.cs[c.id]; !ok {
		return
	}
	delete(n.cs, c.id)
	delete(n.names, c.name)
	close(c.buf)
	c.ws.Close()
	n.inc <- message{
		cid: c.id,
		t:   netUnregister,
	}
	log.WithFields(logrus.Fields{
		"cid":  c.id,
		"name": c.name,
		"addr": c.ws.RemoteAddr(),
	}).Info("Client has been unregistered")
}

// receiver reads frames from the connection then writes messages to the game routine
func (n *network) receiver(c *connection) {
	defer func() {
		n.rw.Lock()
		n.unregister(c)
		n.rw.Unlock()
	}()
	for {
		_, p, err := c.ws.ReadMessage()
		if err != nil {
			log.WithFields(logrus.Fields{
				"cid": c.id,
				"err": err,
			}).Warn("Failed websocket.ReadMessage()")
			return
		}
		log.WithFields(logrus.Fields{
			"cid": c.id,
			"p":   string(p),
		}).Debug("websocket.ReadMessage()")
		m, err := decodeFrame(p)
		if err != nil {
			log.WithFields(logrus.Fields{
				"cid": c.id,
				"p":   string(p),
				"err": err,
			}).Warn("Failed decodeFrame()")
			return
		}
		m.cid = c.id
		n.inc <- m
	}
}

// sender reads frames from the write buffer then writes frames to the connection
func (n *network) sender(c *connection) {
	defer func() {
		n.rw.Lock()
		n.unregister(c)
		n.rw.Unlock()
	}()
	for {
		select {
		case p, ok := <-c.buf:
			if !ok {
				c.ws.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			err := c.ws.WriteMessage(websocket.TextMessage, p)
			if err != nil {
				log.WithFields(logrus.Fields{
					"cid": c.id,
					"err": err,
				}).Warn("Failed websocket.WriteMessage()")
				return
			}
			log.WithFields(logrus.Fields{
				"cid": c.id,
				"p":   string(p),
			}).Debug("websocket.WriteMessage()")
		}
	}
}

// dispatcher reads messages from the game routine then writes frames to the write buffer
func (n *network) dispatcher() {
	for {
		select {
		case m, ok := <-n.out:
			if !ok {
				log.Fatal("Cannot read the outgoing message channel")
			}
			p, err := encodeFrame(m)
			if err != nil {
				log.WithFields(logrus.Fields{
					"m":   m,
					"err": err,
				}).Fatal("Failed encodeFrame()")
			}
			switch m.t {
			case gameTerminate:
				n.rw.Lock()
				if c, ok := n.cs[m.cid]; ok {
					n.unregister(c)
				}
				n.rw.Unlock()
			default:
				n.rw.RLock()
				for _, c := range n.cs {
					c.buf <- p
				}
				n.rw.RUnlock()
			}
		}
	}
}

// newConnection initializes a connection
func newConnection(id uint64, name string, ws *websocket.Conn) *connection {
	return &connection{
		id:   id,
		name: name,
		ws:   ws,
		buf:  make(chan []byte, 1024),
	}
}

// decodeFrame validates JSON schema and converts a JSON text into a incoming message
func decodeFrame(p []byte) (message, error) {
	var f frame
	if err := json.Unmarshal(p, &f); err != nil {
		return message{}, err
	}
	var d map[string]interface{}
	if err := json.Unmarshal(*f.Data, &d); err != nil {
		return message{}, err
	}
	// todo validate d
	return message{
		t: f.Type,
		d: d,
	}, nil
}

// encodeFrame converts a outgoing message into a JSON text and validates JSON schema
func encodeFrame(m message) ([]byte, error) {
	// todo validate m.data
	d := new(json.RawMessage)
	var err error
	*d, err = json.Marshal(m.d)
	if err != nil {
		return nil, err
	}
	return json.Marshal(frame{
		Type: m.t,
		Data: d,
	})
}
