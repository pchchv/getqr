package getqr

import (
	"log"

	bitset "github.com/pchchv/getqr/bitset"
)

const (
	dataEncoderType1To9 dataEncoderType = iota
	dataEncoderType10To26
	dataEncoderType27To40
)

// A dataEncoder encodes data for a particular QR Code version
type dataEncoder struct {
	minVersion                int            //Minimum version supported
	maxVersion                int            // Maximum version supported.
	numericModeIndicator      *bitset.Bitset // Mode indicator bit sequences.
	alphanumericModeIndicator *bitset.Bitset
	byteModeIndicator         *bitset.Bitset
	numericCharCountBits      int // Character count lengths
	alphanumericCharCountBits int
	byteCharCountBits         int
	data                      []byte    // The raw input data.
	actual                    []segment // The data classified into unoptimised segmentss
	optimised                 []segment // The data classified into optimised segments.
}

type dataEncoderType uint8

// A segment encoding mode
type dataMode uint8

// segment is a single segment of data.
type segment struct {
	dataMode dataMode // Data Mode (e.g. numeric)
	data     []byte   // segment data (e.g. "abc")
}

func newDataEncoder(t dataEncoderType) *dataEncoder {
	d := &dataEncoder{}
	switch t {
	case dataEncoderType1To9:
		d = &dataEncoder{
			minVersion:                1,
			maxVersion:                9,
			numericModeIndicator:      bitset.New(false, false, false, true),
			alphanumericModeIndicator: bitset.New(false, false, true, false),
			byteModeIndicator:         bitset.New(false, true, false, false),
			numericCharCountBits:      10,
			alphanumericCharCountBits: 9,
			byteCharCountBits:         8,
		}
	case dataEncoderType10To26:
		d = &dataEncoder{
			minVersion:                10,
			maxVersion:                26,
			numericModeIndicator:      bitset.New(false, false, false, true),
			alphanumericModeIndicator: bitset.New(false, false, true, false),
			byteModeIndicator:         bitset.New(false, true, false, false),
			numericCharCountBits:      12,
			alphanumericCharCountBits: 11,
			byteCharCountBits:         16,
		}
	case dataEncoderType27To40:
		d = &dataEncoder{
			minVersion:                27,
			maxVersion:                40,
			numericModeIndicator:      bitset.New(false, false, false, true),
			alphanumericModeIndicator: bitset.New(false, false, true, false),
			byteModeIndicator:         bitset.New(false, true, false, false),
			numericCharCountBits:      14,
			alphanumericCharCountBits: 13,
			byteCharCountBits:         16,
		}
	default:
		log.Panic("Unknown dataEncoderType")
	}
	return d
}
