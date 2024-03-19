package hash

// primeRK is the prime base
const primeRK = 16777619

type Hash struct {
	pow uint32
	h   uint32
}

// Initialize a hash and calculates pow
func Init(blockSize int) *Hash {
	hash := &Hash{pow: 1, h: 0}
	var mul uint32 = primeRK

	for i := blockSize; i > 0; i >>= 1 {
		if i&1 != 0 {
			hash.pow *= mul
		}
		mul *= mul
	}

	return hash
}

// Calculates the Rabin-Karp hash for a block
func (hash *Hash) HashBlock(sep []byte) *Hash {
	hash.h = 0
	for i := 0; i < len(sep); i++ {
		hash.h = hash.h*primeRK + uint32(sep[i])
	}

	return hash
}

// Calculates rolling Rabin-Karp hash
func (hash *Hash) HashRolling(sep []byte) *Hash {
	hash.h = hash.h*primeRK + uint32(sep[len(sep)-1]) - hash.pow*uint32(sep[0])

	return hash
}

// Gets only the hash
func (h *Hash) GetHash() uint32 {
	return h.h
}
