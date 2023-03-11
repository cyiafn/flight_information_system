package bytes

import (
	"encoding/binary"
	"math"
)

func Int32ToBytes(a int32) []byte {
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, uint32(a))
	return buf
}

func Int64ToBytes(a int64) []byte {
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, uint64(a))
	return buf
}

func ToInt32(a []byte) int32 {
	return int32(binary.LittleEndian.Uint32(a))
}

func ToInt64(a []byte) int64 {
	return int64(binary.LittleEndian.Uint64(a))
}

func Float64ToBytes(a float64) []byte {
	var buf [8]byte
	binary.LittleEndian.PutUint64(buf[:], math.Float64bits(a))
	return buf[:]
}

func ToFloat64(a []byte) float64 {
	bits := binary.LittleEndian.Uint64(a)
	float := math.Float64frombits(bits)
	return float
}
