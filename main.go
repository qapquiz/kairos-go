package main

import (
	"bytes"
	"encoding/binary"
	"kairos-go/packet"
	"kairos-go/remote"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]remote.Remote) // connected clients
var broadcast = make(chan BinaryMessage)              // broadcast channel
var receiveBinaryChannel = make(chan BinaryMessage)

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

func determineListenPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		return ":80"
	}

	return ":" + port
}

func main() {
	runServer(determineListenPort())
}

func runServer(port string) {
	http.HandleFunc("/ws", handleConnections) // ws://ip:port/ws

	go handleBinaryMessages() // thread::spawn

	log.Printf("http server started on %s", port)
	err := http.ListenAndServe(port, nil)
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
		messageType, data, err := ws.ReadMessage()
		log.Println("data: ", data)
		if err != nil {
			log.Printf("error: %v", err)
			delete(clients, ws)
			break
		}

		switch messageType {
		case websocket.TextMessage:
			log.Println("Data from socket(TextMessage)", string(data))
		case websocket.BinaryMessage:
			receiveBinaryChannel <- BinaryMessage{
				Client: ws,
				Data:   data,
			}
		case websocket.CloseMessage:
		case websocket.PingMessage:
		case websocket.PongMessage:
		}
	}
}

func handleBinaryMessages() {
	for {
		binaryMessage := <-receiveBinaryChannel

		reader := createReaderFromData(binaryMessage.Data)

		var packetID uint16
		err := binary.Read(reader, binary.LittleEndian, &packetID)
		if err != nil {
			log.Println("binary.Read failed: ", err)
		}

		remoteClient := clients[binaryMessage.Client]

		packet.ReceiveMessage(packetID, reader, remoteClient, clients)
	}
}

func createReaderFromData(data []byte) *bytes.Reader {
	return bytes.NewReader(data)
}
