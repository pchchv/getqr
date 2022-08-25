package getqr

import (
	"log"

	bitset "github.com/pchchv/getqr/bitset"
)

const (
	dataEncoderType1To9 dataEncoderType = iota
	dataEncoderType10To26
	dataEncoderType27To40
	// Each dataMode is a subset of the subsequent dataMode:
	// dataModeNone < dataModeNumeric < dataModeAlphanumeric < dataModeByte
	// This ordering is important for determining which data modes a character can be encoded with.
	// E.g. 'E' can be encoded in both dataModeAlphanumeric and dataModeByte.
	dataModeNone dataMode = 1 << iota
	dataModeNumeric
	dataModeAlphanumeric
	dataModeByte
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

// Classifies the raw data into unoptimised segments.
// e.g. "123ZZ#!#!" => [numeric, 3, "123"] [alphanumeric, 2, "ZZ"] [byte, 4, "#!#!"].
// Returns the highest data mode needed to encode the data.
// e.g. for a mixed numeric/alphanumeric input, the highest is alphanumeric.
// dataModeNone < dataModeNumeric < dataModeAlphanumeric < dataModeByte
func (d *dataEncoder) classifyDataModes() dataMode {
	var start int
	mode := dataModeNone
	highestRequiredMode := mode
	for i, v := range d.data {
		newMode := dataModeNone
		switch {
		case v >= 0x30 && v <= 0x39:
			newMode = dataModeNumeric
		case v == 0x20 || v == 0x24 || v == 0x25 || v == 0x2a || v == 0x2b || v ==
			0x2d || v == 0x2e || v == 0x2f || v == 0x3a || (v >= 0x41 && v <= 0x5a):
			newMode = dataModeAlphanumeric
		default:
			newMode = dataModeByte
		}
		if newMode != mode {
			if i > 0 {
				d.actual = append(d.actual, segment{dataMode: mode, data: d.data[start:i]})

				start = i
			}
			mode = newMode
		}
		if newMode > highestRequiredMode {
			highestRequiredMode = newMode
		}
	}
	d.actual = append(d.actual, segment{dataMode: mode, data: d.data[start:len(d.data)]})
	return highestRequiredMode
}

// Returns the segment header bits for a segment of type dataMode
func (d *dataEncoder) modeIndicator(dataMode dataMode) *bitset.Bitset {
	switch dataMode {
	case dataModeNumeric:
		return d.numericModeIndicator
	case dataModeAlphanumeric:
		return d.alphanumericModeIndicator
	case dataModeByte:
		return d.byteModeIndicator
	default:
		log.Panic("Unknown data mode")
	}
	return nil
}
