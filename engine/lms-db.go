package engine

import (
	"lms-db/config"
	"lms-db/engine/meta"
	"lms-db/pool"
)

type LmsDB struct {
	bufferPool  *pool.ManagerPool
	MetaManager *meta.MetaManager
}

func NewLmsDB(config *config.Config) *LmsDB {
	return &LmsDB{
		bufferPool: pool.NewManagerPool(config),
	}
}
