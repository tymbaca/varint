package varint

// encode signed value to bits
func encodeZigZag(n int64) uint64 {
	return uint64((n << 1) ^ (n >> 63))
}

// decode signed value from bits
func decodeZigZag(n uint64) int64 {
	return int64((n >> 1) ^ -(n & 1))
}
