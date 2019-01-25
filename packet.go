package main

import (
	"bytes"
	"fmt"
	"kairos-go/packet"
	"log"
)

// Packet is represent each packet
type Packet uint16

const (
	// CSLogin Client send login to channel
	CSLogin Packet = 10001
)

var packetMapper = map[Packet]func(c *Client, r *packet.Reader){
	CSLogin: receiveLogin,
}

// PacketReceiver for receive message from websocket connection
type PacketReceiver struct {
	receive chan ClientWithData
}

func newPacketReceiver() *PacketReceiver {
	return &PacketReceiver{
		receive: make(chan ClientWithData, 256),
	}
}

func readPacketID(r *packet.Reader) Packet {
	return Packet(r.ReadUInt16())
}

func (pr *PacketReceiver) run() {
	for {
		select {
		case clientWithData := <-pr.receive:

			packetReader := &packet.Reader{
				BytesReader: bytes.NewReader(clientWithData.data),
			}
			packet := readPacketID(packetReader)

			packetFunc, ok := packetMapper[packet]
			if !ok {
				log.Println("Not found this packet: ", uint16(packet))
				return
			}

			packetFunc(clientWithData.client, packetReader)
		}
	}
}

func receiveLogin(c *Client, r *packet.Reader) {
	name := r.ReadString()
	fmt.Println("name: ", name)

	c.send <- []byte{}
}
