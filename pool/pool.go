package pool

import (
	"lms-db/config"
	"lms-db/engine"
	"lms-db/index/bptree"
	"lms-db/pool/optimize"
	"lms-db/pool/redo"
)

// ManagerPool buffer pool
type ManagerPool struct {
	pagesManager *PageManager             // page manager
	adapterIndex *optimize.AdapterManager // adapter index
	logManager   *redo.WALLogManager      // WAL write ahead log manager
	dataManager  *engine.DataManager      // data manager
	index        *bptree.BPTree
}

func NewManagerPool(config *config.Config) *ManagerPool {
	return &ManagerPool{
		pagesManager: NewPageManager(),
		adapterIndex: optimize.NewAdapterManager(),
		logManager:   redo.NewLogManager(config),
		dataManager:  engine.NewDataManager(config.GetStoreConfig().DataSpace),
		index:        bptree.NewBPTree(bptree.NewRootNode()),
	}
}

// todo
func (mp *ManagerPool) Get(key string) Record {
	// 读自适应hash索引
	if di, ok := mp.adapterIndex.GetIndex(key); ok {
		// 读到了
		// 读page缓冲池
		if page, ok := mp.pagesManager.GetPage(di.Id); ok {
			var re Record
			var offset int16
			if re, offset = page.ReadByKey(key); offset != -1 {
				mp.adapterIndex.SetIndex(key, optimize.NewAdapterIndex(di.Id, offset))
			}
			return re
		} else {
			// 根据 page id 读 b+树，查找到page

		}
	} else {
		// 没读到
		// 根据 key 读 b+树，找到具体的节点

	}
	return Record{}
}

func (mq *ManagerPool) Set(key string, value interface{}) {

}
