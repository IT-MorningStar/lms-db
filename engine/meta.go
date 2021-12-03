package engine

import (
	"encoding/binary"
	"hash/crc32"
	"lms-db/engine/storage"
)

// todo 元数据相关

type MetaManager struct {
	file     *storage.FileAccess
	fileName string
	num      int64 // log 记录数
	buffChan chan Meta
}

func NewMetaManager() {

}

// write file sync
func (mm *MetaManager) WriteFile() {

}

func (mm *MetaManager) WriteSyncFile(meta Meta) {
	mm.buffChan <- meta
}

type Meta struct {
	metaType uint8
	crc32    uint32
	size     int64
	data     []byte
}

func (m *Meta) Encode() []byte {
	result := make([]byte, 1+4+8, 1+4+8+len(m.data))
	result[0] = m.metaType
	binary.BigEndian.PutUint32(result[1:5], crc32.ChecksumIEEE(m.data))
	binary.BigEndian.PutUint64(result[5:13], uint64(m.size))
	result = append(result, m.data...)
	return result
}
