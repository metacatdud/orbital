package cryptographer

import (
	"encoding/binary"
	"time"
)

type Timestamp int64

func Now() Timestamp {
	return Timestamp(time.Now().UnixMicro())
}

func (t Timestamp) Time() time.Time {
	return time.UnixMicro(int64(t))
}

func (t Timestamp) Bytes() []byte {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(t))
	return buf
}
