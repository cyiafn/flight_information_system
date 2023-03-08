package bytes

import (
	"encoding/binary"
	"log"
	"math"
	"unsafe"
)

func ToBytes[T any](a T) []byte {
	switch v := any(a).(type) {
	case int32:
		return Int32ToBytes(v)
	case int:
		return Int32ToBytes(int32(v)) // Int adopts the bits of the OS, for the purposes of this DB, we assume all integers fit into a 32-bit space
	case string:
		return ToBytes(v)
	}
	return nil // this should never happen
}

func Int32ToBytes(a int32) []byte {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(a))
	return buf
}

func Int64ToBytes(a int64) []byte {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(a))
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

func StringToBytes(a string) []byte {
	return []byte(a)
}

func ToString(a []byte) string {
	return string(a)
}

func ByteArrToBytePtrArr(a []byte) []*byte {
	newArr := make([]*byte, len(a))
	for i, v := range a {
		v := v
		newArr[i] = &v
	}
	return newArr
}

func Equals(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func BytePtrArrToByteArr(a []*byte) []byte {
	newArr := make([]byte, len(a))
	for i, v := range a {
		v := v
		if v == nil {
			log.Fatal("byte conversion encountered a nil byte")
		}
		newArr[i] = *v
	}
	return newArr
}

// StructPtrToBytes takes in a pointer to a structure or a primitive and translates it into an array of bytes
func StructPtrToBytes[T any](a *T) []byte {
	return *(*[]byte)(unsafe.Pointer(a))
}

// ToStructPtr simply takes in an array of bytes and casts it to a pointer of a specific type.
func ToStructPtr[T any](a []byte) *T {
	return (*T)(unsafe.Pointer(&a))
}

func IsNil(a []*byte) bool {
	for _, v := range a {
		if v == nil {
			return true
		}
	}
	return false
}

func IsEmpty(a []byte) bool {
	for _, v := range a {
		if v != 0 {
			return false
		}
	}
	return true
}
