package util

// Hash key for
func Hash(key []byte) uint {
	var nr, nr2 uint = 1, 4
	for i := len(key); i > 0; i-- {
		nr ^= ((nr & 63) + nr2) * (uint(key[i-1]) + (nr << 8))
		nr2 += 3
	}
	return nr
}

var generateLsnInstance *GeneratorLsn

// GenerateLSN  generates a log sequence number
func GenerateLSN() int64 {
	if generateLsnInstance == nil {
		if o, err := NewGeneratorLsn(1); err == nil {
			generateLsnInstance = o
		} else {
			panic(err)
		}
	}
	return generateLsnInstance.GetId()
}
