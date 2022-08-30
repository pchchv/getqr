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
