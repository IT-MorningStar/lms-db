package optimize

import (
	"lms-db/constant"
	"lms-db/strategy"
)

// AdapterManager 自适应hash索引
type AdapterManager struct {
	strategy strategy.LruStrategy
}

func NewAdapterManager() *AdapterManager {
	return &AdapterManager{}
}

func (a *AdapterManager) SetIndex(key string, value *AdapterIndex) {
	a.strategy.Set(key, value, nil)
}

func (a *AdapterManager) GetIndex(key string) (*AdapterIndex, bool) {

	if v, ok := a.strategy.Get(key); ok {
		return v.Value.(*AdapterIndex), true
	} else {
		return nil, false
	}
}

// AdapterIndex 自适应hash索引
type AdapterIndex struct {
	Id     constant.PageId // page id
	Offset uint16          // []byte Offset
}

// NewAdapterIndex 创建一个自适应hash索引
func NewAdapterIndex(id constant.PageId, offset uint16) *AdapterIndex {
	return &AdapterIndex{
		Id:     id,
		Offset: offset,
	}
}

func (a *AdapterIndex) SetOffset(offset uint16) {
	a.Offset = offset
}

func (a *AdapterIndex) GetOffset() uint16 {
	return a.Offset
}
