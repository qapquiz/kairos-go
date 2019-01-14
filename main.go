package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"kairos-go/packet"
	"kairos-go/remote"
	"log"
	"net/http"
	"os"

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

func determineListenPort() (string, error) {
	port := os.Getenv("PORT")
	if port == "" {
		return "", fmt.Errorf("$PORT not set")
	}
	return ":" + port, nil
}

func main() {
	http.HandleFunc("/ws", handleConnections) // ws://ip:port/ws

	go handleBinaryMessages()

	port, err := determineListenPort()
	if err != nil {
		log.Fatal("$PORT is not set")
	}

	log.Printf("http server started on %s", port)
	err = http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

	// log.Println("http server started on :8000")
	// err := http.ListenAndServe(":8000", nil)
	// if err != nil {
	// 	log.Fatal("ListenAndServe: ", err)
	// }
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
		messageType, data, err := ws.ReadMessage()
		log.Println("data: ", data)
		if err != nil {
			log.Println("error: %v", err)
			delete(clients, ws)
			break
		}

		if messageType == websocket.BinaryMessage {
			receiveBinayChannel <- BinaryMessage{
				Client: ws,
				Data:   data,
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
