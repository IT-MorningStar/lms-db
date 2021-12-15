package pool

import (
	"lms-db/config"
	"lms-db/pool/optimize"
	"lms-db/pool/redo"
)

// ManagerPool buffer pool
type ManagerPool struct {
	pages        *PageManager             // page manager
	adapterIndex *optimize.AdapterManager // adapter index
	logManager   *redo.WALLogManager      // WAL write ahead log manager
	// todo: implement
	dataManager int // data manager
}

func NewManagerPool(config *config.Config) *ManagerPool {
	return &ManagerPool{
		pages:        NewPageManager(),
		adapterIndex: optimize.NewAdapterManager(),
		logManager:   redo.NewLogManager(config),
	}
}
