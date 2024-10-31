package cryptographer

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
)

// writeBytesToBuffer writes length-prefixed bytes to the buffer.
func writeBytesToBuffer(buf *bytes.Buffer, b []byte) {
	_ = binary.Write(buf, binary.BigEndian, uint32(len(b)))
	buf.Write(b)
}

func writeTLVtoBuffer(buf *bytes.Buffer, fieldType byte, data []byte) error {
	// Write Type
	if err := buf.WriteByte(fieldType); err != nil {
		return err
	}

	// Write Length
	if len(data) > math.MaxUint32 {
		return fmt.Errorf("%w:[%d]", ErrTLVLenExceed, len(data))
	}

	length := uint32(len(data))
	lengthBytes := make([]byte, 4) // 4 bytes hold for uint32
	binary.LittleEndian.PutUint32(lengthBytes, length)
	if _, err := buf.Write(lengthBytes); err != nil {
		return err
	}

	// Write Value
	if _, err := buf.Write(data); err != nil {
		return err
	}

	return nil
}
