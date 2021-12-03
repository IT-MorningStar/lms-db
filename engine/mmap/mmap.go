package mmap

import (
	"golang.org/x/exp/mmap"
	"sync"
)

type MMap struct {
	rFile *mmap.ReaderAt
	mutex sync.RWMutex
}

func (m *MMap) Read(p []byte) (n int, err error) {
	m.mutex.RLocker()
	defer m.mutex.RUnlock()
	return m.rFile.ReadAt(p, 0)
}

func (m *MMap) ReadAt(p []byte, offset int64) (n int, err error) {
	m.mutex.RLocker()
	defer m.mutex.RUnlock()
	return m.rFile.ReadAt(p, offset)
}

func NewMMap(path string) (*MMap, error) {
	if file, err := mmap.Open(path); err != nil {
		return &MMap{}, err
	} else {
		return &MMap{rFile: file}, nil
	}
}

func (m *MMap) Close() error {
	return m.rFile.Close()
}
