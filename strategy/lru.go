package strategy

import (
	"container/list"
	"errors"
	"sync"
)

// 内存策略 lru
type LruStrategy struct {
	capacity int
	LruMap   map[string]*list.Element
	LruList  *list.List
	mutex    sync.Mutex
}

type LruStruct struct {
	key   string
	value interface{}
}

type CallBackFunc func(r LruStruct) bool

func (ls *LruStrategy) Get(key []byte) (LruStruct, bool) {
	if ls.capacity <= 0 && key == nil && len(ls.LruMap) == 0 {
		return LruStruct{}, false
	}
	return ls.get(string(key))
}

func (ls *LruStrategy) get(key string) (LruStruct, bool) {
	ls.mutex.Lock()
	defer ls.mutex.Unlock()
	if el, ok := ls.LruMap[key]; ok {
		ls.LruList.MoveToFront(el)
		return el.Value.(LruStruct), true
	} else {
		return LruStruct{}, false
	}
}

func (ls *LruStrategy) Set(key []byte, value interface{}, h CallBackFunc) {
	if ls.capacity > 0 && key != nil {
		ls.set(string(key), value, h)
	}
}

func (ls *LruStrategy) set(key string, value interface{}, fn CallBackFunc) {
	ls.mutex.Lock()
	defer ls.mutex.Unlock()
	if el, ok := ls.LruMap[key]; ok {
		ls.LruList.MoveToFront(el)
	} else {
		ls.LruList.PushFront(LruStruct{key: key, value: value})
		if ls.LruList.Len() > ls.capacity {
			element := ls.LruList.Back()
			ele := element.Value.(LruStruct)
			if fn(ele) {
				ls.LruList.Remove(element)
				delete(ls.LruMap, ele.key)
			} else {
				panic(errors.New("callback function execute fail"))
			}
		}
	}
}
