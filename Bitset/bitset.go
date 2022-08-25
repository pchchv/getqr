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

// Ensures the Bitset can store an additional numBits. The underlying array is expanded if necessary.
// To prevent frequent reallocation, expanding the underlying array at least doubles its capacity.
func (b *Bitset) ensureCapacity(numBits int) {
	numBits += b.numBits

	newNumBytes := numBits / 8
	if numBits%8 != 0 {
		newNumBytes++
	}

	if len(b.bits) >= newNumBytes {
		return
	}

	b.bits = append(b.bits, make([]byte, newNumBytes+2*len(b.bits))...)
}
