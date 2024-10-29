package cryptographer

import (
	"encoding/binary"
	"time"
)

type Timestamp int64

func Now() Timestamp {
	return Timestamp(time.Now().UnixNano() / 1000)
}

func (t Timestamp) Time() time.Time {
	return time.Unix(int64(t), 0)
}

func (t Timestamp) Bytes() []byte {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(t))
	return buf
}
