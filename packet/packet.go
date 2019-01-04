package packet

import (
	"bytes"
	"kairos-go/packet_reader"
	"kairos-go/packet_writer"
	"kairos-go/remote"
	"log"

	"github.com/gorilla/websocket"
)

type PacketName uint16

const (
	CSLogin PacketName = 10001

	SCLoggedIn PacketName = 20001
)

type Packet struct{}

func ReceiveMessage(packetID uint16, reader *bytes.Reader, remoteClient remote.Remote, clients map[*websocket.Conn]remote.Remote) {
	switch PacketName(packetID) {
	case CSLogin:
		log.Println("CSLogin")

		packetReader := packet_reader.PacketReader{
			BytesReader: reader,
		}

		name := packetReader.ReadString()
		number := packetReader.ReadInt8()

		log.Println(name)
		log.Println(number)

		remoteClient.Send(sendReceiveLoggedIn(), clients)
	}
}

func sendReceiveLoggedIn() []byte {
	var data = []interface{}{
		uint16(SCLoggedIn),
		"armariya",
	}

	return packet_writer.WritePacket(data)
}
