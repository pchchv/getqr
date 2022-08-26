package reedsolomon

// Polynomial over GF(2^8).
type gfPoly struct {
	// The ith value is the coefficient of the ith degree of x.
	// term[0]*(x^0) + term[1]*(x^1) + term[2]*(x^2) ...
	term []gfElement
}
