package bitset

import "log"

// Bitset stores an array of bits.
type Bitset struct {
	numBits int    // The number of bits stored
	bits    []byte // Storage for individual bits
}

func New(v ...bool) *Bitset {
	b := &Bitset{numBits: 0, bits: make([]byte, 0)}
	b.AppendBools(v...)
	return b
}

// Ensures the Bitset can store an additional numBits. The underlying array is expanded if necessary
// To prevent frequent reallocation, expanding the underlying array at least doubles its capacity
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

// AppendBools appends bits to the Bitset.
func (b *Bitset) AppendBools(bits ...bool) {
	b.ensureCapacity(len(bits))

	for _, v := range bits {
		if v {
			b.bits[b.numBits/8] |= 0x80 >> uint(b.numBits%8)
		}
		b.numBits++
	}
}

// Len returns the length of the Bitset in bits
func (b *Bitset) Len() int {
	return b.numBits
}

// Returns the value of the bit at index
func (b *Bitset) At(index int) bool {
	if index >= b.numBits {
		log.Panicf("Index %d out of range", index)
	}
	return (b.bits[index/8] & (0x80 >> byte(index%8))) != 0
}

// Append bits copied from other. The new length is b.Len() + other.Len().
func (b *Bitset) Append(other *Bitset) {
	b.ensureCapacity(other.numBits)
	for i := 0; i < other.numBits; i++ {
		if other.At(i) {
			b.bits[b.numBits/8] |= 0x80 >> uint(b.numBits%8)
		}
		b.numBits++
	}
}

// Append the numBits least significant bits from value
func (b *Bitset) AppendUint32(value uint32, numBits int) {
	b.ensureCapacity(numBits)
	if numBits > 32 {
		log.Panicf("numBits %d out of range 0-32", numBits)
	}
	for i := numBits - 1; i >= 0; i-- {
		if value&(1<<uint(i)) != 0 {
			b.bits[b.numBits/8] |= 0x80 >> uint(b.numBits%8)
		}
		b.numBits++
	}
}

// Append the numBits least significant bits from value
func (b *Bitset) AppendByte(value byte, numBits int) {
	b.ensureCapacity(numBits)
	if numBits > 8 {
		log.Panicf("numBits %d out of range 0-8", numBits)
	}
	for i := numBits - 1; i >= 0; i-- {
		if value&(1<<uint(i)) != 0 {
			b.bits[b.numBits/8] |= 0x80 >> uint(b.numBits%8)
		}
		b.numBits++
	}
}

// Appends num bits of value value.
func (b *Bitset) AppendNumBools(num int, value bool) {
	for i := 0; i < num; i++ {
		b.AppendBools(value)
	}
}

// Returns a byte consisting of upto 8 bits starting at index
func (b *Bitset) ByteAt(index int) byte {
	if index < 0 || index >= b.numBits {
		log.Panicf("Index %d out of range", index)
	}
	var result byte
	for i := index; i < index+8 && i < b.numBits; i++ {
		result <<= 1
		if b.At(i) {
			result |= 1
		}
	}
	return result
}

// Appends a list of whole bytes
func (b *Bitset) AppendBytes(data []byte) {
	for _, d := range data {
		b.AppendByte(d, 8)
	}
}

// Returns a copy
func Clone(from *Bitset) *Bitset {
	return &Bitset{numBits: from.numBits, bits: from.bits[:]}
}
