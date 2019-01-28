package packet

import (
	"bytes"
	"encoding/binary"
	"log"
)

// Writer store byte buffer for writing.
type Writer struct {
	ByteBuffer *bytes.Buffer
}

// NewWriter create new instance of writer.
func NewWriter(packetID uint16) *Writer {
	writer := &Writer{
		ByteBuffer: new(bytes.Buffer),
	}
	writer.WriteUInt16(packetID)
	return writer
}

// GetData perform getting data as byte array.
func (packetWriter *Writer) GetData() []byte {
	return packetWriter.ByteBuffer.Bytes()
}

// WriteUInt8 perform writing uint8 data to byte buffer.
func (packetWriter *Writer) WriteUInt8(data uint8) {
	packetWriter.write(data)
}

// WriteUInt16 perform writing uint16 data to byte buffer.
func (packetWriter *Writer) WriteUInt16(data uint16) {
	packetWriter.write(data)
}

// WriteUInt32 perform writing uint32 data to byte buffer.
func (packetWriter *Writer) WriteUInt32(data uint32) {
	packetWriter.write(data)
}

// WriteUInt64 perform writing uint64 data to byte buffer.
func (packetWriter *Writer) WriteUInt64(data uint64) {
	packetWriter.write(data)
}

// WriteInt8 perform writing int8 data to byte buffer.
func (packetWriter *Writer) WriteInt8(data int8) {
	packetWriter.write(data)
}

// WriteInt16 perform writing int16 data to byte buffer.
func (packetWriter *Writer) WriteInt16(data int16) {
	packetWriter.write(data)
}

// WriteInt32 perform writing int32 data to byte buffer.
func (packetWriter *Writer) WriteInt32(data int32) {
	packetWriter.write(data)
}

// WriteInt64 perform writing int64 data to byte buffer.
func (packetWriter *Writer) WriteInt64(data int64) {
	packetWriter.write(data)
}

// WriteFloat32 perform writing float32 data to byte buffer.
func (packetWriter *Writer) WriteFloat32(data float32) {
	packetWriter.write(data)
}

// WriteFloat64 perform writing float64 data to byte buffer.
func (packetWriter *Writer) WriteFloat64(data float64) {
	packetWriter.write(data)
}

// WriteString perform writing string data to byte buffer.
func (packetWriter *Writer) WriteString(data string) {
	packetWriter.write(data)
}

// WriteBoolean perform writing boolean data to byte buffer.
func (packetWriter *Writer) WriteBoolean(data bool) {
	packetWriter.write(data)
}

func (packetWriter *Writer) write(data interface{}) {
	switch data.(type) {
	case string:
		packetWriter.ByteBuffer.Write([]byte(data.(string)))
		packetWriter.ByteBuffer.WriteByte(uint8(0))
	default:
		if err := binary.Write(packetWriter.ByteBuffer, binary.LittleEndian, data); err != nil {
			log.Fatal("binary.Write failed: ", data, err)
		}
	}
}
