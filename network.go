package crescent

import (
	"net/http"
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/websocket"
)

type ClientID uint64

type Network struct {
	mu  *sync.RWMutex
	seq ClientID
	id  map[ClientID]*Client

	upgrader websocket.Upgrader

	r InstanceOutput
	w InstanceInputWriter
}

type Client struct {
	id  ClientID
	ws  *websocket.Conn
	buf chan []byte
}

// NewNetwork returns a Network
func NewNetwork(origin string, r InstanceOutput, w InstanceInputWriter) *Network {
	return &Network{
		mu:  new(sync.RWMutex),
		seq: 0,
		id:  make(map[ClientID]*Client),

		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				if origin == "" || origin == r.Header.Get("Origin") {
					return true
				}
				log.WithFields(logrus.Fields{
					"Origin": r.Header.Get("Origin"),
					"Host":   r.Host,
				}).Warn("Unacceptable Origin header")
				return false
			},
		},

		r: r,
		w: w,
	}
}

// Run starts the network routine
func (n *Network) Run(addr string) {
	go n.dispatcher()
	http.HandleFunc("/", n.wsEffect)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.WithField("err", err).Fatal("Failed http.ListenAndServe()")
	}
}

// generateClientID generates a unique ClientID
func (n *Network) generateClientID() ClientID {
	n.seq++
	return n.seq
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

// wsEffect handles a HTTP request
func (n *Network) wsEffect(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method Not Allowed", 405)
		return
	}
	var c *Client
	var ok bool
	n.sync(func() {
		ws, err := n.upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.WithField("err", err).Warn("Failed websocket.Upgrade()")
			return
		}
		c = NewClient(n.generateClientID(), ws)
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
	n.w.Write(c.id, InputConnect{})
	log.WithFields(logrus.Fields{
		"id":   c.id,
		"addr": c.ws.RemoteAddr(),
	}).Info("Client has been registered")
}

// unregister unregisters the Client
func (n *Network) unregister(c *Client) {
	if _, ok := n.id[c.id]; !ok {
		return
	}
	delete(n.id, c.id)
	close(c.buf)
	c.ws.Close()
	n.w.Write(c.id, InputDisconnect{})
	log.WithFields(logrus.Fields{
		"id":   c.id,
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
				"id":  c.id,
				"err": err,
			}).Warn("Failed websocket.ReadMessage()")
			return
		}
		in, err := DecodeInputFrame(p)
		if err != nil {
			log.WithFields(logrus.Fields{
				"id":  c.id,
				"p":   string(p),
				"err": err,
			}).Warn("Failed DecodeInputFrame()")
			return
		}
		n.w.Write(c.id, in)
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
						"id":  c.id,
						"err": err,
					}).Warn("Failed websocket.WriteMessage()")
				}
				return
			}
			err := c.ws.WriteMessage(websocket.TextMessage, p)
			if err != nil {
				log.WithFields(logrus.Fields{
					"id":  c.id,
					"err": err,
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
		case out, ok := <-n.r:
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
func NewClient(id ClientID, ws *websocket.Conn) *Client {
	return &Client{
		id:  id,
		ws:  ws,
		buf: make(chan []byte, 1024),
	}
}
