package main

import "fmt"

// PacketReceiver for receive message from websocket connection
type PacketReceiver struct {
	receive chan []byte
}

func newPacketReceiver() *PacketReceiver {
	return &PacketReceiver{
		receive: make(chan []byte, 256),
	}
}

func (pr *PacketReceiver) run() {
	for {
		select {
		case message := <-pr.receive:
			fmt.Println("message: ", message)
		}
	}
}
