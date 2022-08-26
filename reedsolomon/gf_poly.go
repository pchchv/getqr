package reedsolomon

import bitset "github.com/pchchv/getqr/bitset"

// Polynomial over GF(2^8).
type gfPoly struct {
	// The ith value is the coefficient of the ith degree of x.
	// term[0]*(x^0) + term[1]*(x^1) + term[2]*(x^2) ...
	term []gfElement
}

// Returns data as a polynomial over GF(2^8)
// Each data byte becomes the coefficient of an x term
// For an n byte input the polynomial is:
// data[n-1]*(x^n-1) + data[n-2]*(x^n-2) ... + data[0]*(x^0).
func newGFPolyFromData(data *bitset.Bitset) gfPoly {
	numTotalBytes := data.Len() / 8
	if data.Len()%8 != 0 {
		numTotalBytes++
	}
	result := gfPoly{term: make([]gfElement, numTotalBytes)}
	i := numTotalBytes - 1
	for j := 0; j < data.Len(); j += 8 {
		result.term[i] = gfElement(data.ByteAt(j))
		i--
	}
	return result
}

// Returns the number of
func (e gfPoly) numTerms() int {
	return len(e.term)
}

func (e gfPoly) normalised() gfPoly {
	numTerms := e.numTerms()
	maxNonzeroTerm := numTerms - 1
	for i := numTerms - 1; i >= 0; i-- {
		if e.term[i] != 0 {
			break
		}
		maxNonzeroTerm = i - 1
	}
	if maxNonzeroTerm < 0 {
		return gfPoly{}
	} else if maxNonzeroTerm < numTerms-1 {
		e.term = e.term[0 : maxNonzeroTerm+1]
	}
	return e
}
