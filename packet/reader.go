package packet

import (
	"bytes"
	"encoding/binary"
	"log"
)

type Reader struct {
	BytesReader *bytes.Reader
}

// Use for ReadString from []bytes read until found charCode 0 (it's mean that string end)
func (packetReader *Reader) ReadString() string {
	message := make([]byte, 0)

	for {
		if charCode, err := packetReader.BytesReader.ReadByte(); err == nil {
			if uint8(charCode) == 0 {
				break
			}

			message = append(message, charCode)
		}
	}

	return string(message[:])
}

func (packetReader *Reader) ReadUInt8() uint8 {
	number, err := packetReader.BytesReader.ReadByte()
	if err != nil {
		log.Println("ReadUInt8 Error!")
	}

	return uint8(number)
}

func (packetReader *Reader) ReadUInt16() uint16 {
	var number uint16
	err := binary.Read(packetReader.BytesReader, binary.LittleEndian, &number)
	if err != nil {
		log.Println("binary.Read uint16 failed: ", err)
	}

	return number
}

func (packetReader *Reader) ReadUInt32() uint32 {
	var number uint32
	err := binary.Read(packetReader.BytesReader, binary.LittleEndian, &number)
	if err != nil {
		log.Println("binary.Read uint32 failed: ", err)
	}

	return number
}

func (packetReader *Reader) ReadUInt64() uint64 {
	var number uint64
	err := binary.Read(packetReader.BytesReader, binary.LittleEndian, &number)
	if err != nil {
		log.Println("binary.Read uint64 failed: ", err)
	}

	return number
}

func (packetReader *Reader) ReadInt8() int8 {
	number, err := packetReader.BytesReader.ReadByte()
	if err != nil {
		log.Println("ReadInt8 Error!")
	}

	return int8(number)
}

func (packetReader *Reader) ReadInt16() int16 {
	var number int16
	err := binary.Read(packetReader.BytesReader, binary.LittleEndian, &number)
	if err != nil {
		log.Println("binary.Read int16 failed: ", err)
	}

	return number
}

func (packetReader *Reader) ReadInt32() int32 {
	var number int32
	err := binary.Read(packetReader.BytesReader, binary.LittleEndian, &number)
	if err != nil {
		log.Println("binary.Read int32 failed: ", err)
	}

	return number
}

func (packetReader *Reader) ReadInt64() int64 {
	var number int64
	err := binary.Read(packetReader.BytesReader, binary.LittleEndian, &number)
	if err != nil {
		log.Println("binary.Read int64 failed: ", err)
	}

	return number
}

func (packetReader *Reader) ReadFloat32() float32 {
	var number float32
	err := binary.Read(packetReader.BytesReader, binary.LittleEndian, &number)
	if err != nil {
		log.Println("binary.Read float32 failed: ", err)
	}

	return number
}

func (packetReader *Reader) ReadFloat64() float64 {
	var number float64
	err := binary.Read(packetReader.BytesReader, binary.LittleEndian, &number)
	if err != nil {
		log.Println("binary.Read float64 failed: ", err)
	}

	return number
}

func (packetReader *Reader) ReadBoolean() bool {
	number, err := packetReader.BytesReader.ReadByte()
	if err != nil {
		log.Println("ReadBoolean Error!")
	}

	return uint8(number) == 1
}
