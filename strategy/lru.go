package strategy

import (
	"container/list"
	"sync"
)

// 内存策略 lru
type LruStrategy struct {
	capacity int
	lruMap   map[string]*list.Element
	lruList  *list.List
	mutex    sync.RWMutex
}

type LruStruct struct {
	Key   string
	Value interface{}
}

func NewLruStrategy(capacity int) *LruStrategy {
	return &LruStrategy{
		capacity: capacity,
		lruMap:   make(map[string]*list.Element),
		lruList:  list.New(),
	}
}

type CallBackFunc func(r *LruStruct) bool

func (ls *LruStrategy) Get(key string) (*LruStruct, bool) {
	if ls.capacity <= 0 && len(key) == 0 && len(ls.lruMap) == 0 {
		return nil, false
	}
	return ls.get(string(key))
}

func (ls *LruStrategy) get(key string) (*LruStruct, bool) {
	ls.mutex.RLock()
	defer ls.mutex.RUnlock()
	if el, ok := ls.lruMap[key]; ok {
		ls.lruList.MoveToFront(el)
		return el.Value.(*LruStruct), true
	} else {
		return nil, false
	}
}

func (ls *LruStrategy) Set(key string, value interface{}, h CallBackFunc) {
	if ls.capacity > 0 && len(key) != 0 {
		ls.set(key, value, h)
	}
}

// 这里可能死锁
func (ls *LruStrategy) set(key string, value interface{}, fn CallBackFunc) {
	ls.mutex.Lock()
	defer ls.mutex.Unlock()
	if el, ok := ls.lruMap[key]; ok {
		ls.lruList.MoveToFront(el)
	} else {
		ls.lruMap[key] = ls.lruList.PushFront(&LruStruct{Key: key, Value: value})
		if ls.lruList.Len() > ls.capacity {
			element := ls.lruList.Back()
			ele := element.Value.(*LruStruct)
			ls.lruList.Remove(element)
			delete(ls.lruMap, ele.Key)
			if fn != nil {
				fn(ele)
			}
		}
	}
}
