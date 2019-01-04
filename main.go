package main

import (
	"bytes"
	"encoding/binary"
	"kairos-go/packet"
	"kairos-go/remote"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]remote.Remote) // connected clients
var broadcast = make(chan BinaryMessage)              // broadcast channel
var receiveBinayChannel = make(chan BinaryMessage)

// Configure the upgrader
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type BinaryMessage struct {
	Client *websocket.Conn
	Data   []byte
}

func main() {
	http.HandleFunc("/ws", handleConnections)

	go handleBinaryMessages()

	log.Println("http server started on :8000")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	defer ws.Close()

	clients[ws] = remote.Remote{
		State:  remote.Connected,
		Socket: ws,
	}

	remoteClient := clients[ws]
	remoteClient.OnConnected()

	for {
		messageType, p, err := ws.ReadMessage()
		log.Println("data: ", p)
		if err != nil {
			log.Println("error: %v", err)
			delete(clients, ws)
			break
		}

		if messageType == websocket.BinaryMessage {
			receiveBinayChannel <- BinaryMessage{
				Client: ws,
				Data:   p,
			}
		}
	}
}

func handleBinaryMessages() {
	for {
		binaryMessage := <-receiveBinayChannel

		reader := bytes.NewReader(binaryMessage.Data)

		var packetID uint16
		err := binary.Read(reader, binary.LittleEndian, &packetID)
		if err != nil {
			log.Println("binary.Read failed: ", err)
		}

		remoteClient := clients[binaryMessage.Client]

		packet.ReceiveMessage(packetID, reader, remoteClient, clients)
	}
}
