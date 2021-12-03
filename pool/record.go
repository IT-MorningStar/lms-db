package pool

import "encoding/binary"

// Record
//   | StructureType|  KLength |  Key   |  VLength ｜ Value ｜
//   |      1B      |     1B   |   nB   |     4B   ｜  nB   ｜
type Record struct {
	StructureType uint8  // string、int、float.......
	KLength       uint8  // key length
	Key           string // key name
	VLength       uint32 // value length
	Value         []byte // value
}

func (r *Record) Encode() (result []byte, err error) {
	l := 1 + 4 + int32(r.KLength) + 4 + int32(r.VLength)
	result = make([]byte, l, l)
	result[0] = r.StructureType
	result[1] = r.KLength
	copy(result[2:2+r.KLength], r.Key)
	binary.BigEndian.PutUint32(result[2+r.KLength:6+r.KLength], r.VLength)
	copy(result[6+r.KLength:], r.Value)
	return
}

func NewRecord(StructureType uint8, KLength uint8, Key string, VLength uint32, Value []byte) Record {
	return Record{
		StructureType: StructureType,
		KLength:       KLength,
		Key:           Key,
		VLength:       VLength,
		Value:         Value,
	}
}
