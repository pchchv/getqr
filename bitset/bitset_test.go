package bitset

import "testing"

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
