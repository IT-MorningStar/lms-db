package pool

import (
	"encoding/binary"
	"lms-db/constant"
	"lms-db/strategy"
	"sync"
)

type PageManager struct {
	strategy strategy.LruStrategy
}

// NewPageManager new a  PageManager
func NewPageManager() *PageManager {
	return &PageManager{}
}

// Page 4K ，与操作系统的page cache刚好对应上 Data最大可用字节，4060字节
//  lsn  | Crc32 |  Id  | Prev | Next  |  Num ｜ leaf  +  预留   +  idleSpace |  Data ｜
//  8B   |  4B   |  4B  |  8B  |  8B   |  2B  ｜  1b   +   3b   +   12b      ｜  nB   ｜
type Page struct {
	Lsn    int64           // log sequence number
	Crc32  uint32          // crc32 check sum
	Id     constant.PageId // page id
	Prev   int64           // previous node position
	Next   int64           // next node position
	PrevId constant.PageId // previous page id
	NextId constant.PageId // next page id
	Num    uint16          // current node keys num
	Leaf   bool            // is leaf node，highest 1 bit
	Idle   uint16          // high 3 bit reserve ，low 12 bit storage idle space
	Dirty  bool            // dirty pages
	Data   []byte          // data spaces
	Mutex  sync.RWMutex
}

// ReadByKey  load data by key in page
func (p *Page) ReadByKey(key string) (result Record, off int16) {
	p.Mutex.RLocker()
	defer p.Mutex.RUnlock()
	var n = 0
	for i := 0; i < len(p.Data); i += n {
		var st uint8
		var kl uint8
		var k string
		var vl uint32
		var v []byte
		st = p.Data[n]
		kl = p.Data[n+1]
		k = string(p.Data[n+2 : int(kl)+n+2])
		binary.BigEndian.PutUint32(p.Data[int(kl)+n+2:int(kl)+n+2+4], vl)
		v = p.Data[int(kl)+n+2+4 : int(kl)+n+2+4+int(vl)]
		if k == key {
			result = Record{
				KLength:       kl,
				StructureType: st,
				Key:           k,
				VLength:       vl,
				Value:         v,
			}
			return
		}
		n = n + 1 + 1 + int(kl) + 4 + int(vl)
	}

	return Record{}, -1
}

// ReadKeyByOffset 根据key和offset从页中读出来
func (p *Page) ReadKeyByOffset(key string, offset uint16) (result Record) {
	p.Mutex.RLocker()
	defer p.Mutex.RUnlock()
	var n = int(offset)
	var st uint8
	var kl uint8
	var k string
	var vl uint32
	var v []byte
	st = p.Data[n]
	kl = p.Data[n+1]
	k = string(p.Data[n+2 : int(kl)+n+2])
	binary.BigEndian.PutUint32(p.Data[int(kl)+n+2:int(kl)+n+2+4], vl)
	v = p.Data[int(kl)+n+2+4 : int(kl)+n+2+4+int(vl)]
	if k == key {
		result = Record{
			KLength:       kl,
			StructureType: st,
			Key:           k,
			VLength:       vl,
			Value:         v,
		}
		return
	}
	return
}

// Write todo 设置key 到page 中，复杂！复杂！复杂！！！
func (p *Page) Write(key string, data []byte) (id constant.PageId, off uint16) {
	p.Mutex.Lock()
	defer p.Mutex.Unlock()
	return 0, 0
}
