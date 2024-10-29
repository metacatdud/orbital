package cryptographer

import (
	"bytes"
	"encoding/binary"
)

func intToByte[T int8 | int16 | int32 | int64](val T) []byte {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, val); err != nil {
		return nil
	}
	return buf.Bytes()
}
