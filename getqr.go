package getqr

import "image/color"

type QRCode struct {
	Content         string      // Original content encoded
	BackgroundColor color.Color // User settable drawing options
	ForegroundColor color.Color
	Border          bool // QR Code border. True — borders are enabled
	encoder         *dataEncoder
	data            *bitset
	synbol          *symbol
	mask            int
}
