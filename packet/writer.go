package packet

import (
	"bytes"
	"encoding/binary"
	"log"
)

func WritePacket(data []interface{}) []byte {
	buffer := new(bytes.Buffer)

	for _, v := range data {
		switch v.(type) {
		case string:
			buffer.Write([]byte(v.(string)))
			buffer.WriteByte(uint8(0))
		default:
			err := binary.Write(buffer, binary.LittleEndian, v)
			if err != nil {
				log.Fatal("binary.Write failed: ", err)
			}
		}
	}

	//log.Println("Send data back: ", buffer.Bytes())

	return buffer.Bytes()
}
