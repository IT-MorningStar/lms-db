package config

import "sync"

type (
	// Store parameters
	Store struct {
		WorkSpace string `yaml:"WorkSpace"` // work directory
		DataSpace string `yaml:"DataSpace"` // data directory
		LogSpace  string `yaml:"LogSpace"`  // log directory
	}
	// Optimize  parameters
	Optimize struct {
		AdapterIndexNum int `yaml:"AdapterIndexNum"` // adapter index number 200
	}

	PageParams struct {
		PageSize           int `yaml:"pageSize"`           // page size 4k
		PageMaxNum         int `yaml:"PageMaxNum"`         // max page num   500
		IdleMaxNum         int `yaml:"IdleMaxNum"`         // max idle page num  100
		IdleMinNum         int `yaml:"idleMinNum"`         // min idle page num  20
		FlushDiskSecond    int `yaml:"FlushDiskSecond"`    // how often seconds to flush 1s
		CheckPointUpdate   int `yaml:"CheckPointUpdate"`   // how often check point update 10s
		FlushMaxDirtyPage  int `yaml:"FlushMaxDirtyPage"`  // max dirty page flush count/10s 10
		MetaLogFileMaxLine int `yaml:"MetaLogFileMaxLine"` // meta log max line count,create new file exceeded
	}

	RedoLog struct {
		FlushDiskSecond    int   `yaml:"FlushDiskSecond"`    // how many seconds to flush 1
		Sync               bool  `yaml:"Sync"`               // sync file true
		RedoLogFileMaxSize int64 `yaml:"RedoLogFileMaxSize"` //
	}
)

type Config struct {
	store       *Store
	performance *Optimize
	redoLog     *RedoLog
	pageParams  *PageParams
	mutex       sync.RWMutex
}

func NewConfig() *Config {
	return &Config{
		store: &Store{
			WorkSpace: "/Users/lichangxiao/dinfull/Project/Golang/lms-db/LMS-DB/work",
			DataSpace: "/Users/lichangxiao/dinfull/Project/Golang/lms-db/LMS-DB/data",
			LogSpace:  "/Users/lichangxiao/dinfull/Project/Golang/lms-db/LMS-DB/log",
		},
		performance: &Optimize{AdapterIndexNum: 200},
		redoLog: &RedoLog{
			FlushDiskSecond:    1,
			Sync:               true,
			RedoLogFileMaxSize: 2 * 1024 * 1024,
		},
		pageParams: &PageParams{
			PageSize:          4 * 1024,
			PageMaxNum:        500,
			IdleMaxNum:        100,
			IdleMinNum:        20,
			FlushDiskSecond:   1,
			CheckPointUpdate:  10,
			FlushMaxDirtyPage: 10,
		},
	}
}

func (c *Config) GetStoreConfig() *Store {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.store
}

func (c *Config) GetPerformanceConfig() *Optimize {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.performance
}

func (c *Config) GetRedoLogConfig() *RedoLog {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.redoLog
}

func (c *Config) GetPageParams() *PageParams {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.pageParams
}
