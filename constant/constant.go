package constant

const LogicLogType = 1
const (
	_         uint8 = iota
	InsertLog       // insert type log
	DeleteLog       // delete type log
)

const PhysicalLogType = 2

const (
	_ uint8 = iota
	PageSplit
)

// meta log 数据类型
const (
	_                  uint8 = iota
	MetaCheckPointType       // meta check point type log
)

// KeyMaxLength key max key length
const KeyMaxLength = 255

// ValueMaxLength pass 3000 bytes storage to extend file
const ValueMaxLength = 3000

type PageId uint

type KeyType string
