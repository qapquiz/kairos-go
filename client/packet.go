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

	SCLogin Packet = 20001
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
	age := reader.ReadInt8()
	fmt.Printf("\nLogin with name : %s and age : %v", name, age)

	var friends []string
	friends = append(friends, "ball", "P'Jo", "P'Gun", "P'Nan")

	packetWriter := packet.NewWriter(uint16(SCLogin))
	packetWriter.WriteString(name)
	packetWriter.WriteInt8(age)
	packetWriter.WriteBoolean(true)
	packetWriter.WriteUInt8(uint8(len(friends)))
	for _, v := range friends {
		packetWriter.WriteString(v)
	}
	remote.send <- packetWriter.GetData()
}
