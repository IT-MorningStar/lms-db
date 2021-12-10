package ds

const (
	_ uint8 = iota
	StringDataType
	FloatDataType
	IntDataType
)

type DataStructure interface {
	Encode()
	Decode()
}
