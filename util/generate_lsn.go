package util

import (
	"errors"
	"sync"
	"time"
)

//模型发号器
const (
	workerBits  uint8 = 10                      //分布式机器bit占位数
	numberBits  uint8 = 12                      //计数器的所bit占位数
	workerMax   int64 = -1 ^ (-1 << workerBits) //最大分布式机器个数
	numberMax   int64 = -1 ^ (-1 << numberBits) //每毫秒最大计数
	timeShift   uint8 = workerBits + numberBits
	workerShift uint8 = numberBits
	startTime   int64 = 1595983389000
)

type GeneratorLsn struct {
	mutex     sync.Mutex
	timestamp int64
	workerId  int64
	number    int64
}

func NewGeneratorLsn(workerId int64) (*GeneratorLsn, error) {
	if workerId < 0 || workerId > workerMax {
		return nil, errors.New("Illegal Worker ID ")
	}
	// 生成一个新节点
	return &GeneratorLsn{
		timestamp: 0,
		workerId:  workerId,
		number:    0,
	}, nil
}

func (w *GeneratorLsn) GetId() int64 {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	now := time.Now().UnixNano() / 1e6
	if w.timestamp == now {
		w.number++
		if w.number > numberMax {
			for now <= w.timestamp {
				now = time.Now().UnixNano() / 1e6
			}
		}
	} else {
		w.number = 0
		w.timestamp = now
	}
	ID := int64((now-startTime)<<timeShift | (w.workerId << workerShift) | (w.number))
	return ID
}
