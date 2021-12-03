package redo

import "encoding/binary"

var _ WALLogStructure = (*LogicLogStructure)(nil)

// 逻辑日志为了主从同步使用
type LogicLogStructure struct {
	LogType uint8  // buffer pool log type
	KLength uint8  // key length
	Key     string // key name
	VLength uint16 // value length
	Value   string // value
}

func NewLogStructure(lt uint8, kl uint8, k string, vl uint16, v string) LogicLogStructure {
	return LogicLogStructure{
		LogType: lt,
		KLength: kl,
		Key:     k,
		VLength: vl,
		Value:   v,
	}
}

// Encode encoding log structure struct
func (ls *LogicLogStructure) Encode() (result []byte, err error) {
	result = append(result, ls.LogType)
	result = append(result, ls.KLength)
	result = append(result, []byte(ls.Key)...)
	var vl = make([]byte, 2)
	binary.BigEndian.PutUint16(vl, ls.VLength)
	result = append(result, vl...)
	result = append(result, []byte(ls.Value)...)
	return
}

func (ls *LogicLogStructure) Decode([]byte) (err error) {
	return nil
}
