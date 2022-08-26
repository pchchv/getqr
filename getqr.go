package getqr

import (
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
