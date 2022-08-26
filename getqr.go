package getqr

import (
	"errors"
	"fmt"
	"image/color"
	"log"

	bitset "github.com/pchchv/getqr/bitset"
)

type QRCode struct {
	Content         string        // Original content encoded
	Level           RecoveryLevel // QR Code type.
	VersionNumber   int
	BackgroundColor color.Color // User settable drawing options
	ForegroundColor color.Color
	Border          bool // QR Code border. True — borders are enabled
	encoder         *dataEncoder
	version         qrCodeVersion
	data            *bitset.Bitset
	synbol          *symbol
	mask            int
}

// Constructs a QR Code. An error occurs if the content is too long
func New(content string, level RecoveryLevel) (*QRCode, error) {
	var err error
	var encoder *dataEncoder
	var encoded *bitset.Bitset
	var chosenVersion *qrCodeVersion
	encoders := []dataEncoderType{dataEncoderType1To9, dataEncoderType10To26, dataEncoderType27To40}
	for _, t := range encoders {
		encoder = newDataEncoder(t)
		encoded, err = encoder.encode([]byte(content))
		if err != nil {
			continue
		}
		chosenVersion = chooseQRCodeVersion(level, encoder, encoded.Len())
		if chosenVersion != nil {
			break
		}
	}
	if err != nil {
		return nil, err
	} else if chosenVersion == nil {
		return nil, errors.New("content too long to encode")
	}
	q := &QRCode{
		Content:         content,
		Level:           level,
		VersionNumber:   chosenVersion.version,
		ForegroundColor: color.Black,
		BackgroundColor: color.White,
		encoder:         encoder,
		data:            encoded,
		version:         *chosenVersion,
	}
	return q, nil
}

// Constructs a QR code of a specific version. An error occurs in case of invalid version
func NewWithForcedVersion(content string, version int, level RecoveryLevel) (*QRCode, error) {
	var encoder *dataEncoder
	switch {
	case version >= 1 && version <= 9:
		encoder = newDataEncoder(dataEncoderType1To9)
	case version >= 10 && version <= 26:
		encoder = newDataEncoder(dataEncoderType10To26)
	case version >= 27 && version <= 40:
		encoder = newDataEncoder(dataEncoderType27To40)
	default:
		return nil, fmt.Errorf("Invalid version %d (expected 1-40 inclusive)", version)
	}
	var encoded *bitset.Bitset
	encoded, err := encoder.encode([]byte(content))
	if err != nil {
		return nil, err
	}
	chosenVersion := getQRCodeVersion(level, version)
	if chosenVersion == nil {
		return nil, errors.New("cannot find QR Code version")
	}
	if encoded.Len() > chosenVersion.numDataBits() {
		return nil, fmt.Errorf("Cannot encode QR code: content too large for fixed size QR Code version %d (encoded length is %d bits, maximum length is %d bits)",
			version,
			encoded.Len(),
			chosenVersion.numDataBits())
	}
	q := &QRCode{
		Content:         content,
		Level:           level,
		VersionNumber:   chosenVersion.version,
		ForegroundColor: color.Black,
		BackgroundColor: color.White,
		encoder:         encoder,
		data:            encoded,
		version:         *chosenVersion,
	}
	return q, nil
}

// Adds final terminator bits to the encoded data. The number of terminator bits required is determined when the QR Code version is chosen
// The terminator bits are thus added after the QR Code version is chosen, rather than at the data encoding stage.
func (q *QRCode) addTerminatorBits(numTerminatorBits int) {
	q.data.AppendNumBools(numTerminatorBits, false)
}

// Pads the encoded data upto the full length required.
func (q *QRCode) addPadding() {
	numDataBits := q.version.numDataBits()
	if q.data.Len() == numDataBits {
		return
	}
	// Pad to the nearest codeword boundary.
	q.data.AppendNumBools(q.version.numBitsToPadToCodeword(q.data.Len()), false)
	// Pad codewords 0b11101100 and 0b00010001.
	padding := [2]*bitset.Bitset{
		bitset.New(true, true, true, false, true, true, false, false),
		bitset.New(false, false, false, true, false, false, false, true),
	}
	// Insert pad codewords alternately.
	i := 0
	for numDataBits-q.data.Len() >= 8 {
		q.data.Append(padding[i])
		i = 1 - i // Alternate between 0 and 1.
	}
	if q.data.Len() != numDataBits {
		log.Panicf("BUG: got len %d, expected %d", q.data.Len(), numDataBits)
	}
}
