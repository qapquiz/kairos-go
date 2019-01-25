package client

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

var packetMapper = map[Packet]func(remote *Remote, reader *packet.Reader){
	CSLogin: receiveLogin,
}

// PacketReceiver for receive message from websocket connection
type PacketReceiver struct {
	receive chan RemoteWithData
}

func newPacketReceiver() *PacketReceiver {
	return &PacketReceiver{
		receive: make(chan RemoteWithData, 256),
	}
}

func readPacketID(r *packet.Reader) Packet {
	return Packet(r.ReadUInt16())
}

func (pr *PacketReceiver) run() {
	for {
		select {
		case remoteWithData := <-pr.receive:

			packetReader := &packet.Reader{
				BytesReader: bytes.NewReader(remoteWithData.data),
			}
			packet := readPacketID(packetReader)

			packetFunc, ok := packetMapper[packet]
			if !ok {
				log.Println("Not found this packet: ", uint16(packet))
				return
			}

			packetFunc(remoteWithData.remote, packetReader)
		}
	}
}

func receiveLogin(remote *Remote, reader *packet.Reader) {
	name := reader.ReadString()
	fmt.Println("name: ", name)

	remote.send <- []byte{}
}
