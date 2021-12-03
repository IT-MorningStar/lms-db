package optimize

import (
	"lms-db/constant"
	"sync"
)

// AdapterManager 自适应hash索引
type AdapterManager struct {
	index map[constant.KeyType]*AdapterIndex
}

func NewAdapterManager() *AdapterManager {
	return &AdapterManager{
		index: make(map[constant.KeyType]*AdapterIndex),
	}
}

// AdapterIndex 自适应hash索引
type AdapterIndex struct {
	Id     constant.PageId // page id
	Mutex  sync.RWMutex
	Offset uint16 // []byte Offset
}

// NewAdapterIndex 创建一个自适应hash索引
func NewAdapterIndex(id constant.PageId, offset uint16) *AdapterIndex {
	return &AdapterIndex{
		Id:     id,
		Offset: offset,
	}
}

func (a *AdapterIndex) SetOffset(offset uint16) {
	a.Mutex.Lock()
	defer a.Mutex.Unlock()
	a.Offset = offset
}

func (a *AdapterIndex) GetOffset() uint16 {
	return a.Offset
}
