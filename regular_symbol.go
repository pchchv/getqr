package getqr

import bitset "github.com/pchchv/getqr/bitset"

type regularSymbol struct {
	version qrCodeVersion
	mask    int
	data    *bitset.Bitset
	symbol  *symbol
	size    int
}
