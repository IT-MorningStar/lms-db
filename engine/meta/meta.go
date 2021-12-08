package meta

import (
	"encoding/binary"
	"errors"
	"fmt"
	"hash/crc32"
	"io/ioutil"
	"lms-db/constant"
	"lms-db/engine/storage"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

// 缓冲占比多大时进行同步磁盘
const SyncFileRate float64 = 0.9

// 多少秒之后检查一次meta 是否需要切换文件
const FileCheckSeconds float64 = 10

// 日志文件最多多少行
const FileMaxLineLength int64 = 1024

const MaxFileMerge int64 = 20

const FileShuffle string = ".meta-log"

type MetaManager struct {
	file         *storage.FileAccess
	filePath     string
	num          int64 // logfile 当前记录数
	buffChan     chan Meta
	mutex        sync.Mutex
	fileSequence int64
}

func NewMetaManager(fp string, num int64, fileSequence int64) (*MetaManager, error) {
	p := filepath.Join(fp, fmt.Sprintf("%d%s", fileSequence, FileShuffle))
	if a, err := storage.NewFileAccess(p, 0); err != nil {
		return nil, err
	} else {
		return &MetaManager{
			file:         a,
			filePath:     fp,
			num:          num,
			buffChan:     make(chan Meta, 1024),
			fileSequence: fileSequence,
		}, nil
	}
}

func (mm *MetaManager) consumerBuffMsg() {
	for {
		select {
		case r := <-mm.buffChan:
			fmt.Println(r)
			// 检查文件是否需要创建新的文件
		case <-time.Tick(10 * time.Second):
			func() {
				mm.mutex.Lock()
				defer mm.mutex.Unlock()
				len := len(mm.buffChan)
				if float64(mm.num)*SyncFileRate <= float64(len) {
					mm.openNewMetaLogFile(mm.filePath)
				}
			}()
		}
	}
}

// 默认都是强制刷盘
func (mm *MetaManager) WriteFile(meta Meta) {
	mm.mutex.Lock()
	defer mm.mutex.Unlock()
	len := len(mm.buffChan)
	if float64(mm.num)*SyncFileRate <= float64(len) {
		mm.WriteSyncFile(meta)
	} else {
		mm.num += 1
		mm.buffChan <- meta
	}
}

func (mm *MetaManager) WriteSyncFile(meta Meta) {
	mm.mutex.Lock()
	defer mm.mutex.Unlock()
	if mm.num > FileMaxLineLength {
		mm.openNewMetaLogFile(mm.filePath)
	} else {
		mm.num += 1
		mm.file.WriteAppendEnd(meta.Encode())
	}
}

// 该方法不加锁，调用该方法的地方加锁
func (mm *MetaManager) openNewMetaLogFile(fileName string) {
	if a, err := storage.NewFileAccess(fileName, 0); err != nil {
		panic(err)
	} else {
		mm.file.Close()
		mm.file = a
		mm.num = 0
	}
}

// todo 需要思考如何加锁
func (mm *MetaManager) ClearMetaLogFile() {
	if fs, err := ioutil.ReadDir(filepath.Join(mm.filePath)); err != nil {
		panic(errors.New(fmt.Sprintf("clear meta log file fail,error: %v", err)))
	} else {
		var SecondMax, FirstMax int64

		for _, item := range fs {
			s := strings.Split(item.Name(), ".")
			se, suffix := s[0], s[1]
			if fmt.Sprintf(".%s", suffix) == FileShuffle {
				if serial, err := strconv.ParseInt(se, 10, 64); err != nil {
					panic(errors.New(fmt.Sprintf("clear meta log file fail,error: %v", err)))
				} else {
					if serial > FirstMax {
						FirstMax, SecondMax = serial, FirstMax
					} else {
						if serial > SecondMax {
							SecondMax = serial
						}
					}
				}
			}
		}

	}
}

type Meta struct {
	metaType uint8
	crc32    uint32
	size     int64
	data     []byte
}

type IMetaStruct interface {
	Encode() []byte
}

func (m *Meta) Encode() []byte {
	result := make([]byte, 1+4+8, 1+4+8+len(m.data))
	result[0] = m.metaType
	binary.BigEndian.PutUint32(result[1:5], crc32.ChecksumIEEE(m.data))
	binary.BigEndian.PutUint64(result[5:13], uint64(m.size))
	result = append(result, m.data...)
	return result
}

func (m *Meta) DataDecode() IMetaStruct {
	var result IMetaStruct
	switch m.metaType {
	case constant.MetaCheckPointType:
		result = RedoMetaLogDecode(m.data)
	}
	return result
}

func DecodeMeta(data []byte) (*Meta, error) {
	var result Meta
	if len(data) < 1+4+8 {
		return nil, errors.New("data decode fail,data length is too short")
	} else {
		result.metaType = data[0]
		result.crc32 = binary.BigEndian.Uint32(data[1:5])
		result.size = int64(binary.BigEndian.Uint64(data[5:13]))
		c32 := crc32.ChecksumIEEE(data[13:])
		if c32 != result.crc32 {
			return nil, errors.New("data check sum fail,Illegal data ")
		}
		result.data = data[13:]
		return &result, nil
	}
}
