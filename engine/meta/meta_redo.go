package meta

import "encoding/binary"

var _ IMetaStruct = (*RedoMetaLog)(nil)

type RedoMetaLog struct {
	Lsn   int64
	Start int64
	End   int64
}

const RedoMetaLogBytesLen = 8 + 8 + 8

func NewRedoMetaLog(lsn, start, end int64) *RedoMetaLog {
	return &RedoMetaLog{
		lsn,
		start,
		end,
	}
}

func (r *RedoMetaLog) Encode() []byte {
	result := make([]byte, RedoMetaLogBytesLen, RedoMetaLogBytesLen)
	binary.BigEndian.PutUint64(result[0:8], uint64(r.Lsn))
	binary.BigEndian.PutUint64(result[8:16], uint64(r.Start))
	binary.BigEndian.PutUint64(result[16:24], uint64(r.End))
	return result
}

func RedoMetaLogDecode(data []byte) *RedoMetaLog {
	var result RedoMetaLog
	if len(data) != RedoMetaLogBytesLen {
		return &result
	} else {
		result.Lsn = int64(binary.BigEndian.Uint64(data[0:8]))
		result.Start = int64(binary.BigEndian.Uint64(data[8:16]))
		result.End = int64(binary.BigEndian.Uint64(data[16:24]))
	}
	return &result
}
