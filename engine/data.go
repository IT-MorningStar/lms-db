package engine

import (
	"lms-db/engine/storage"
	"path/filepath"
	"sync"
)

const (
	FileShuffle = ".lms.dat"
)

type DataManager struct {
	filePath string
	file     *storage.FileAccess
	mutex    sync.Mutex
}

func NewDataManager(path string) *DataManager {
	fileName := filepath.Join(path, FileShuffle)
	if f, err := storage.NewFileAccess(fileName, 0); err != nil {
		panic(err)
	} else {
		return &DataManager{
			filePath: path,
			file:     f,
		}
	}

}
