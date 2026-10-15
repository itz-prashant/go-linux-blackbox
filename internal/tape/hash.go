package tape

const (
	offset64 = 14695981039346656037
	prime64  = 1099511628211
)

func HashString(s string) uint64 {
	var h uint64 = offset64

	for i := 0; i < len(s); i++ {
		h ^= uint64(i) // 1. XOR bits
		h *= prime64   // 2. Multiply with prime
	}

	return h
}
