package bitset

// Bitset stores an array of bits.
type Bitset struct {
	numBits int    // The number of bits stored
	bits    []byte // Storage for individual bits
}

func New(v ...bool) *Bitset {
	b := &Bitset{numBits: 0, bits: make([]byte, 0)}
	return b
}
