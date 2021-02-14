package fractal

import (
	"bytes"
	"encoding/binary"
)

func Int64ToBytes(i64 []int64) []byte {
	var bb bytes.Buffer
	if err := binary.Write(&bb, binary.LittleEndian, i64); err != nil {
		panic(err)
	}
	return bb.Bytes()
}

func BytesToInt64(b []byte) []int64 {
	bb := bytes.NewBuffer(b)
	i64 := make([]int64, len(b)/8)
	if err := binary.Read(bb, binary.LittleEndian, i64); err != nil {
		panic(err)
	}
	return i64
}
