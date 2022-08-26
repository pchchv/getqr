package getqr

import (
	"errors"
	"image/color"

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
