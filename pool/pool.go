package pool

import (
	"lms-db/pool/optimize"
	"lms-db/pool/redo"
)

// ManagerPool buffer pool
type ManagerPool struct {
	pages        *PageManager             // page manager
	adapterIndex *optimize.AdapterManager // adapter index
	logManager   *redo.WALLogManager      // insert buffer pool log
}
