package reedsolomon

import (
	"log"
)

// Returns the Reed-Solomon generator polynomial with degree.
// The generator polynomial is calculated as:
// (x + a^0)(x + a^1)...(x + a^degree-1)
func rsGeneratorPoly(degree int) gfPoly {
	if degree < 2 {
		log.Panic("degree < 2")
	}
	generator := gfPoly{term: []gfElement{1}}
	for i := 0; i < degree; i++ {
		nextPoly := gfPoly{term: []gfElement{gfExpTable[i], 1}}
		generator = gfPolyMultiply(generator, nextPoly)
	}
	return generator
}
