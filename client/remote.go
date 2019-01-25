package client

import (
	"bytes"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to pper with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// RemoteWithData bundle remote with data to send into channel
type RemoteWithData struct {
	remote *Remote
	data   []byte
}

// Remote is a middleman between the websocket connection and the hub.
type Remote struct {
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Packet receiver
	packetReceiver *PacketReceiver

	// Buffered channel of outbound messages.
	send chan []byte
}

// red messages from the websocket connection to the packet
func (r *Remote) read() {
	defer func() {
		r.hub.unregister <- r
		r.conn.Close()
	}()

	r.conn.SetReadLimit(maxMessageSize)
	r.conn.SetReadDeadline(time.Now().Add(pongWait))
	r.conn.SetPongHandler(func(string) error {
		r.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	go r.packetReceiver.run()

	for {
		messageType, message, err := r.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		switch messageType {
		case websocket.BinaryMessage:
			r.packetReceiver.receive <- RemoteWithData{
				remote: r,
				data:   message,
			}
			bytes.NewReader(message)
		}

	}
}

func (r *Remote) write() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		r.conn.Close()
	}()

	for {
		select {
		case message, ok := <-r.send:
			r.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				r.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := r.conn.WriteMessage(websocket.BinaryMessage, message); err != nil {
				log.Println(err)
				return
			}
		case <-ticker.C:
			r.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := r.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// ServeWs handles websocket requests from the peer.
func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	remote := &Remote{
		hub:            hub,
		conn:           conn,
		packetReceiver: newPacketReceiver(),
		send:           make(chan []byte, 256),
	}

	remote.hub.register <- remote

	go remote.read()
	go remote.write()
}
