package bitset

import (
	"fmt"
	"log"
)

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

// Appends bits to the Bitset
func (b *Bitset) AppendBools(bits ...bool) {
	b.ensureCapacity(len(bits))
	for _, v := range bits {
		if v {
			b.bits[b.numBits/8] |= 0x80 >> uint(b.numBits%8)
		}
		b.numBits++
	}
}

// Returns the length of the Bitset in bits
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

// Append bits copied from other. The new length is b.Len() + other.Len()
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

// Appends num bits of value value
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

// Returns a substring, consisting of the bits from indexes start to end
func (b *Bitset) Substr(start int, end int) *Bitset {
	if start > end || end > b.numBits {
		log.Panicf("Out of range start=%d end=%d numBits=%d", start, end, b.numBits)
	}
	result := New()
	result.ensureCapacity(end - start)
	for i := start; i < end; i++ {
		if b.At(i) {
			result.bits[result.numBits/8] |= 0x80 >> uint(result.numBits%8)
		}
		result.numBits++
	}
	return result
}

// Returns a human readable representation of the Bitset's contents
func (b *Bitset) String() string {
	var bitString string
	for i := 0; i < b.numBits; i++ {
		if (i % 8) == 0 {
			bitString += " "
		}
		if (b.bits[i/8] & (0x80 >> byte(i%8))) != 0 {
			bitString += "1"
		} else {
			bitString += "0"
		}
	}
	return fmt.Sprintf("numBits=%d, bits=%s", b.numBits, bitString)
}

// Returns the contents of the Bitset
func (b *Bitset) Bits() []bool {
	result := make([]bool, b.numBits)
	var i int
	for i = 0; i < b.numBits; i++ {
		result[i] = (b.bits[i/8] & (0x80 >> byte(i%8))) != 0
	}
	return result
}

// Returns a copy
func Clone(from *Bitset) *Bitset {
	return &Bitset{numBits: from.numBits, bits: from.bits[:]}
}

// Constructs and returns a Bitset from a string
// The string consists of '1', '0' or ' ' characters, e.g. "1010 0101". The '1' and '0'
// characters represent true/false bits respectively, and ' ' characters are ignored
// The function panics if the input string contains other characters
func NewFromBase2String(b2string string) *Bitset {
	b := &Bitset{numBits: 0, bits: make([]byte, 0)}
	for _, c := range b2string {
		switch c {
		case '1':
			b.AppendBools(true)
		case '0':
			b.AppendBools(false)
		case ' ':
		default:
			log.Panicf("Invalid char %c in NewFromBase2String", c)
		}
	}
	return b
}
