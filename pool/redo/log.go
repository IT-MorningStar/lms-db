package redo

import (
	"lms-db/config"
	"lms-db/engine/storage"
	"path/filepath"
	"sync"
)

// 重做日志的后缀
const RedoFileName = "lms-db.redo"

type WALLogManager struct {
	bufferSize int                  // RedoLog max size
	sync       bool                 // if ture sync disk，else tick flush disk
	flushSec   int                  // how many seconds to  flush second
	bufferPool chan WALLogStructure // buffer pool
	mutex      sync.RWMutex
	logFile    *storage.FileAccess
	logPath    string // log file path
	start      int64  // current start offset
	end        int64  // current end offset
	maxSize    int64  // log file max size
}

type WALLogStructure interface {
	Encode() (result []byte, err error)
	Decode([]byte) (err error)
}

func NewLogManager(config *config.Config) *WALLogManager {
	return newLogManager(
		1024,
		config.GetRedoLogConfig().Sync,
		config.GetRedoLogConfig().FlushDiskSecond,
		filepath.Join(config.GetStoreConfig().LogSpace, RedoFileName),
		0,
		0,
		config.GetRedoLogConfig().RedoLogFileMaxSize,
	)
}

// NewLogManager creates a new WALLogManager
func newLogManager(bs int, sync bool, flushSec int, logPath string, start, end, maxSize int64) *WALLogManager {
	if lf, err := storage.NewFileAccess(logPath, 0); err == nil {
		lm := &WALLogManager{
			bufferSize: bs,
			sync:       sync,
			flushSec:   flushSec,
			bufferPool: make(chan WALLogStructure, bs),
			logPath:    logPath,
			logFile:    lf,
			start:      start,
			end:        end,
			maxSize:    maxSize,
		}
		return lm
	} else {
		return &WALLogManager{
			bufferSize: bs,
			sync:       sync,
			flushSec:   flushSec,
			bufferPool: make(chan WALLogStructure, bs),
			logPath:    logPath,
			start:      start,
			end:        end,
			maxSize:    maxSize,
		}
	}
}

// todo 1、判断日志是否超长。
//      2、超长等待---调用redo更新
//      3、不超长循环写入。
func (lm *WALLogManager) Write(ls WALLogStructure, offset int64) (int64, error) {
	if lm.sync {
		lm.writeDisk(ls)
	} else {

	}
	return 0, nil
}

// todo check 写文件和更新meta数据是用一个互斥锁，回写
func (lm *WALLogManager) writeDisk(ls WALLogStructure) (int64, error) {
	lm.mutex.Lock()
	defer lm.mutex.Unlock()
	if d, err := ls.Encode(); err != nil {
		return -1, err
	} else {
		// 写入文件的大小
		l := len(d)
		if lm.end >= lm.start {
			tailIdle := lm.maxSize - lm.end
			if tailIdle >= int64(l) {
				// 空闲空间可以存储
				if m, err := lm.logFile.WriteAt(d, lm.end); err != nil {
					return -1, err
				} else {
					// todo 写入meta文件区
					lm.end += int64(m)

				}
			} else {
				// 判断文件头是不是还可以写
				ds := int64(l) - tailIdle
				if lm.start > ds {
					// 还可以写
					if _, err := lm.logFile.WriteAt(d[0:tailIdle], lm.end); err != nil {
						return -1, err
					} else {
						if _, err := lm.logFile.WriteAt(d[tailIdle:], 0); err != nil {
							return -1, err
						} else {
							// todo 写入meta文件区
							lm.end = ds
						}
					}
				} else if lm.start == ds {
					// todo 正好写满，赶紧记录然后刷盘
					if _, err := lm.logFile.WriteAt(d[0:tailIdle], lm.end); err != nil {
						return -1, err
					} else {
						if _, err := lm.logFile.WriteAt(d[tailIdle:], 0); err != nil {
							return -1, err
						} else {
							// todo 写入meta文件区
							lm.end = ds
							// todo 需要从start 开始读日志进行重做，给新的日志文件留空间，日志刚好写满

						}
					}
				} else {
					// todo 需要从start 开始读日志进行重做，给新的日志文件留空间

				}
			}
		} else {
			idle := lm.start - lm.end
			if idle > int64(l) {
				// 空闲空间可以存储
				if _, err := lm.logFile.WriteAt(d, lm.end); err != nil {
					return -1, err
				} else {
					// todo 写入meta文件区
					lm.end += int64(l)

				}
			} else if idle == int64(l) {
				if _, err := lm.logFile.WriteAt(d, lm.end); err != nil {
					return -1, err
				} else {
					// todo 写入meta文件区
					lm.end += int64(l)
					// todo 需要从start 开始读日志进行重做，给新的日志文件留空间，日志刚好写满

				}

			} else {
				// todo 需要从start 开始读日志进行重做，给新的日志文件留空间

			}
		}

		// todo
		// 1、写入之前判断是否超出日志最大长度，超出则等待，执行当前所有页进行刷盘

	}
	return 0, nil
}
