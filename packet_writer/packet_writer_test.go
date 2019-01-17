package packetwriter

import (
	"reflect"
	"testing"
)

func TestWritePacket(t *testing.T) {
	var expectedResult = []byte{20, 116, 101, 115, 116, 0, 8, 2}

	var data = []interface{}{
		uint8(20),
		string("test"),
		uint16(520),
	}

	actualResult := WritePacket(data)

	if !compareBytes(expectedResult, actualResult) {
		t.Errorf("Expected %v, actual %v", expectedResult, actualResult)
	}
}

func compareBytes(firstBytes []byte, secondBytes []byte) bool {
	return reflect.DeepEqual(firstBytes, secondBytes)
}