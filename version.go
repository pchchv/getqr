package getqr

type RecoveryLevel int

type qrCodeVersion struct {
	version          int           // Version number (1-40)
	level            RecoveryLevel // Error recovery level
	dataEncoderType  dataEncoderType
	block            []block // The encoded data can be broken into blocks. They contain data and error recovery bytes. Larger QR codes contain more blocks.
	numRemainderBits int     // Number of bits required to pad the combined data and error correction bit stream up to the symbol's full capacity.
}

type block struct {
	numBlocks        int
	numCodewords     int // Total codewords (numErrorCodewords+numDataCodewords)
	numDataCodewords int // Number of data codewords
}
