package redo

import (
	"encoding/binary"
	"hash/crc32"
	"lms-db/constant"
)

// 数据落盘的物理日志
// PhysicsRedoLog 结构
//  lsn  |  size |   items   |
//   8B  |   8B  |   size-B  |
type PhysicsRedoLog struct {
	lsn   int64  // log sequence number，log redo
	crc32 uint32 // log check sum
	size  int64  // data bytes size
	items []PhysicsRedoLogItem
	data  []byte
}

// 从rodo log文件读出来的构造方法
func NewPhysicsRedoLogItems(lsn int64, c32 uint32, size int64, data []byte) PhysicsRedoLog {

	return PhysicsRedoLog{}
}

// 为落盘准备的构造方法
func NewPhysicsRedoLogData(lsn int64, data []byte) PhysicsRedoLog {
	return PhysicsRedoLog{
		lsn:   lsn,
		crc32: crc32.ChecksumIEEE(data),
		size:  int64(len(data)),
		data:  data,
	}
}

// PhysicsRedoLogItem 结构
//   id  |  offset |  size  |   data    |
//   4B  |   8B    |   8B   |   size-B  |
type PhysicsRedoLogItem struct {
	id     constant.PageId
	offset int64
	size   int64
	data   []byte
}

func NewPhysicsRedoLogItem(id constant.PageId, offset, size int64, data []byte) PhysicsRedoLogItem {
	return PhysicsRedoLogItem{id: id, offset: offset, size: size, data: data}
}

var _ WALLogStructure = (*PhysicsRedoLog)(nil)

func (p *PhysicsRedoLogItem) Encode() ([]byte, error) {
	result := make([]byte, 4+8+8, 4+8+8+len(p.data))
	binary.BigEndian.PutUint32(result[0:4], uint32(p.id))
	binary.BigEndian.PutUint64(result[4:12], uint64(p.offset))
	binary.BigEndian.PutUint64(result[12:18], uint64(p.size))
	result = append(result, p.data...)
	return result, nil
}

func (p *PhysicsRedoLog) Encode() ([]byte, error) {
	result := make([]byte, 8+4+8, 8+4+8)
	binary.BigEndian.PutUint64(result[0:8], uint64(p.lsn))
	binary.BigEndian.PutUint32(result[8:12], uint32(p.crc32))
	binary.BigEndian.PutUint64(result[12:20], uint64(p.size))
	for _, item := range p.items {
		if i, err := (&item).Encode(); err != nil {
			return []byte{}, err
		} else {
			result = append(result, i...)
		}
	}
	return result, nil
}

func (p *PhysicsRedoLog) Decode([]byte) (err error) {
	return nil
}
