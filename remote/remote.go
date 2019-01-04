package remote

import (
	"log"

	"github.com/gorilla/websocket"
)

type Remote struct {
	State  RemoteState
	Socket *websocket.Conn
}

type RemoteState uint8

const (
	Connected    RemoteState = 1
	Disconnected RemoteState = 2
)

func (remote *Remote) OnConnected() {

}

func (remote *Remote) OnDisconnected() {

}

func (remote *Remote) Send(data []byte, clients map[*websocket.Conn]Remote) {
	err := remote.Socket.WriteMessage(websocket.BinaryMessage, data)
	if err != nil {
		log.Println("Send failed: ", err)

		thisRemote := clients[remote.Socket]
		thisRemote.OnDisconnected()

		delete(clients, remote.Socket)
	}
}
