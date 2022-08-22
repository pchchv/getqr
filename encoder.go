package getqr

// A dataEncoder encodes data for a particular QR Code version
type dataEncoder struct {
	numericModeIndicator      *bitset // Mode indicator bit sequences.
	alphanumericModeIndicator *bitset
	byteModeIndicator         *bitset
	numericCharCountBits      int // Character count lengths
	alphanumericCharCountBits int
	byteCharCountBits         int
	data                      []byte    // The raw input data.
	actual                    []segment // The data classified into unoptimised segmentss
	optimised                 []segment // The data classified into optimised segments.
}

// A segment encoding mode
type dataMode uint8

// segment is a single segment of data.
type segment struct {
	dataMode dataMode // Data Mode (e.g. numeric)
	data     []byte   // segment data (e.g. "abc")
}
