package bitset

import (
	"math/rand"
	"testing"
)

func TestNewBitset(t *testing.T) {
	tests := [][]bool{
		{},
		{b1},
		{b0},
		{b1, b0},
		{b1, b0, b1},
		{b0, b0, b1},
	}
	for _, v := range tests {
		result := New(v...)
		if !equal(result.Bits(), v) {
			t.Errorf("%s", result.String())
			t.Errorf("%v => %v, want %v", v, result.Bits(), v)
		}
	}
}

func TestAppend(t *testing.T) {
	randomBools := make([]bool, 128)
	rng := rand.New(rand.NewSource(1))
	for i := 0; i < len(randomBools); i++ {
		randomBools[i] = rng.Intn(2) == 1
	}
	for i := 0; i < len(randomBools)-1; i++ {
		a := New(randomBools[0:i]...)
		b := New(randomBools[i:]...)
		a.Append(b)
		if !equal(a.Bits(), randomBools) {
			t.Errorf("got %v, want %v", a.Bits(), randomBools)
		}
	}
}

func equal(a []bool, b []bool) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestAppendByte(t *testing.T) {
	tests := []struct {
		initial  *Bitset
		value    byte
		numBits  int
		expected *Bitset
	}{
		{
			New(),
			0x01,
			1,
			New(b1),
		},
		{
			New(b1),
			0x01,
			1,
			New(b1, b1),
		},
		{
			New(b0),
			0x01,
			1,
			New(b0, b1),
		},
		{
			New(b1, b0, b1, b0, b1, b0, b1),
			0xAA, // 0b10101010
			2,
			New(b1, b0, b1, b0, b1, b0, b1, b1, b0),
		},
		{
			New(b1, b0, b1, b0, b1, b0, b1),
			0xAA, // 0b10101010
			8,
			New(b1, b0, b1, b0, b1, b0, b1, b1, b0, b1, b0, b1, b0, b1, b0),
		},
	}
	for _, test := range tests {
		test.initial.AppendByte(test.value, test.numBits)
		if !equal(test.initial.Bits(), test.expected.Bits()) {
			t.Errorf("Got %v, expected %v", test.initial.Bits(),
				test.expected.Bits())
		}
	}
}

func TestAppendUint32(t *testing.T) {
	tests := []struct {
		initial  *Bitset
		value    uint32
		numBits  int
		expected *Bitset
	}{
		{
			New(),
			0xAAAAAAAF,
			4,
			New(b1, b1, b1, b1),
		},
		{
			New(),
			0xFFFFFFFF,
			32,
			New(b1, b1, b1, b1, b1, b1, b1, b1, b1, b1, b1, b1, b1, b1, b1, b1, b1,
				b1, b1, b1, b1, b1, b1, b1, b1, b1, b1, b1, b1, b1, b1, b1),
		},
		{
			New(),
			0x0,
			32,
			New(b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0,
				b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0),
		},
		{
			New(),
			0xAAAAAAAA,
			32,
			New(b1, b0, b1, b0, b1, b0, b1, b0, b1, b0, b1, b0, b1, b0, b1, b0, b1,
				b0, b1, b0, b1, b0, b1, b0, b1, b0, b1, b0, b1, b0, b1, b0),
		},
		{
			New(),
			0xAAAAAAAA,
			31,
			New(b0, b1, b0, b1, b0, b1, b0, b1, b0, b1, b0, b1, b0, b1, b0, b1,
				b0, b1, b0, b1, b0, b1, b0, b1, b0, b1, b0, b1, b0, b1, b0),
		},
	}
	for _, test := range tests {
		test.initial.AppendUint32(test.value, test.numBits)
		if !equal(test.initial.Bits(), test.expected.Bits()) {
			t.Errorf("Got %v, expected %v", test.initial.Bits(),
				test.expected.Bits())
		}
	}
}
