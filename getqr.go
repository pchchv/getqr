package getqr

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"io/ioutil"
	"log"
	"os"

	bitset "github.com/pchchv/getqr/bitset"
	reedsolomon "github.com/pchchv/getqr/reedsolomon"
)

type QRCode struct {
	Content         string        // Original content encoded
	Level           RecoveryLevel // QR Code type
	VersionNumber   int
	BackgroundColor color.Color // User settable drawing options
	ForegroundColor color.Color
	DisableBorder   bool // Disable the QR Code border
	Border          bool // QR Code border. True — borders are enabled
	encoder         *dataEncoder
	version         qrCodeVersion
	data            *bitset.Bitset
	symbol          *symbol
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
// The terminator bits are thus added after the QR Code version is chosen, rather than at the data encoding stage
func (q *QRCode) addTerminatorBits(numTerminatorBits int) {
	q.data.AppendNumBools(numTerminatorBits, false)
}

// Pads the encoded data upto the full length required
func (q *QRCode) addPadding() {
	numDataBits := q.version.numDataBits()
	if q.data.Len() == numDataBits {
		return
	}
	// Pad to the nearest codeword boundary
	q.data.AppendNumBools(q.version.numBitsToPadToCodeword(q.data.Len()), false)
	// Pad codewords 0b11101100 and 0b00010001
	padding := [2]*bitset.Bitset{
		bitset.New(true, true, true, false, true, true, false, false),
		bitset.New(false, false, false, true, false, false, false, true),
	}
	// Insert pad codewords alternately
	i := 0
	for numDataBits-q.data.Len() >= 8 {
		q.data.Append(padding[i])
		i = 1 - i // Alternate between 0 and 1
	}
	if q.data.Len() != numDataBits {
		log.Panicf("BUG: got len %d, expected %d", q.data.Len(), numDataBits)
	}
}

// Takes the completed (terminated & padded) encoded data, splits the data into blocks (as specified by the QR Code version),
// applies error correction to each block, then interleaves the blocks together
// The QR Code's final data sequence is returned
func (q *QRCode) encodeBlocks() *bitset.Bitset {
	// Split into blocks
	type dataBlock struct {
		data          *bitset.Bitset
		ecStartOffset int
	}
	block := make([]dataBlock, q.version.numBlocks())
	start := 0
	end := 0
	blockID := 0
	for _, b := range q.version.block {
		for j := 0; j < b.numBlocks; j++ {
			start = end
			end = start + b.numDataCodewords*8
			// Apply error correction to each block.
			numErrorCodewords := b.numCodewords - b.numDataCodewords
			block[blockID].data = reedsolomon.Encode(q.data.Substr(start, end), numErrorCodewords)
			block[blockID].ecStartOffset = end - start
			blockID++
		}
	}
	// Interleave the blocks
	result := bitset.New()
	// Combine data blocks
	working := true
	for i := 0; working; i += 8 {
		working = false
		for j, b := range block {
			if i >= block[j].ecStartOffset {
				continue
			}
			result.Append(b.data.Substr(i, i+8))
			working = true
		}
	}
	// Combine error correction blocks
	working = true
	for i := 0; working; i += 8 {
		working = false
		for j, b := range block {
			offset := i + block[j].ecStartOffset
			if offset >= block[j].data.Len() {
				continue
			}
			result.Append(b.data.Substr(offset, offset+8))
			working = true
		}
	}
	// Append remainder bits.
	result.AppendNumBools(q.version.numRemainderBits, false)
	return result
}

// Returns the QR Code as a 2D array of 1-bit pixels bitmap[y][x] is true if the pixel at (x, y) is set
// The bitmap includes the required "quiet zone" around the QR Code to aid decoding.
func (q *QRCode) Bitmap() [][]bool {
	// Build QR code
	q.encode()
	return q.symbol.bitmap()
}

// Completes the steps required to encode the QR Code. These include adding the terminator bits and padding,
// splitting the data into blocks and applying the error correction, and selecting the best data mask
func (q *QRCode) encode() {
	numTerminatorBits := q.version.numTerminatorBitsRequired(q.data.Len())
	q.addTerminatorBits(numTerminatorBits)
	q.addPadding()
	encoded := q.encodeBlocks()
	const numMasks int = 8
	penalty := 0
	for mask := 0; mask < numMasks; mask++ {
		var s *symbol
		var err error
		s, err = buildRegularSymbol(q.version, mask, encoded, !q.DisableBorder)
		if err != nil {
			log.Panic(err.Error())
		}
		numEmptyModules := s.numEmptyModules()
		if numEmptyModules != 0 {
			log.Panicf("bug: numEmptyModules is %d (expected 0) (version=%d)",
				numEmptyModules, q.VersionNumber)
		}
		p := s.penaltyScore()
		if q.symbol == nil || p < penalty {
			q.symbol = s
			q.mask = mask
			penalty = p
		}
	}
}

// Returns the QR Code as an image.Image
// A positive size sets a fixed image width and height (e.g. 256 yields an 256x256px image)
// Depending on the amount of data encoded, fixed size images can have different amounts of padding (white space around the QR Code)
// As an alternative, a variable sized image can be generated instead: A negative size causes a variable sized image to be returned
// The image returned is the minimum size required for the QR Code. Choose a larger negative number to increase the scale of the image
// e.g. a size of -5 causes each module (QR Code "pixel") to be 5px in size
func (q *QRCode) Image(size int) image.Image {
	// Build QR code
	q.encode()
	// Minimum pixels (both width and height) required
	realSize := q.symbol.size
	// Variable size support
	if size < 0 {
		size = size * -1 * realSize
	}
	// Actual pixels available to draw the symbol. Automatically increase the image size if it's not large enough
	if size < realSize {
		size = realSize
	}
	// Output image
	rect := image.Rectangle{Min: image.Point{0, 0}, Max: image.Point{size, size}}
	// Saves a few bytes to have them in this order
	p := color.Palette([]color.Color{q.BackgroundColor, q.ForegroundColor})
	img := image.NewPaletted(rect, p)
	fgClr := uint8(img.Palette.Index(q.ForegroundColor))
	// QR code bitmap
	bitmap := q.symbol.bitmap()
	// Map each image pixel to the nearest QR code module
	modulesPerPixel := float64(realSize) / float64(size)
	for y := 0; y < size; y++ {
		y2 := int(float64(y) * modulesPerPixel)
		for x := 0; x < size; x++ {
			x2 := int(float64(x) * modulesPerPixel)
			v := bitmap[y2][x2]
			if v {
				pos := img.PixOffset(x, y)
				img.Pix[pos] = fgClr
			}
		}
	}
	return img
}

// Returns the QR Code as a PNG image
// Size is both the image width and height in pixels
// If size is too small then a larger image is silently returned
// Negative values for size cause a variable sized image to be returned:
// See the documentation for Image()
func (q *QRCode) PNG(size int) ([]byte, error) {
	img := q.Image(size)
	encoder := png.Encoder{CompressionLevel: png.BestCompression}
	var b bytes.Buffer
	err := encoder.Encode(&b, img)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

// Writes the QR Code as a PNG image to io.Writer size is both the image width and height in pixels
// If size is too small then a larger image is silently written
// Negative values for size cause a variable sized image to be written: See the documentation for Image()
func (q *QRCode) Write(size int, out io.Writer) error {
	png, err := q.PNG(size)
	if err != nil {
		return err
	}
	_, err = out.Write(png)
	return err
}

// Writes the QR Code as a PNG image to the specified file size is both the image width and height in pixels
// If size is too small then a larger image is silently written
// Negative values for size cause a variable sized image to be written: See the documentation for Image()
func (q *QRCode) WriteFile(size int, filename string) error {
	png, err := q.PNG(size)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, png, os.FileMode(0644))
}

// Produces a multi-line string that forms a QR-code image
func (q *QRCode) ToString(inverseColor bool) string {
	bits := q.Bitmap()
	var buf bytes.Buffer
	for y := range bits {
		for x := range bits[y] {
			if bits[y][x] != inverseColor {
				buf.WriteString("  ")
			} else {
				buf.WriteString("██")
			}
		}
		buf.WriteString("\n")
	}
	return buf.String()
}

// Produces a multi-line string that forms a QR-code image, a factor two smaller in x and y then ToString
func (q *QRCode) ToSmallString(inverseColor bool) string {
	bits := q.Bitmap()
	var buf bytes.Buffer
	// If there is an odd number of rows, the last one needs special treatment
	for y := 0; y < len(bits)-1; y += 2 {
		for x := range bits[y] {
			if bits[y][x] == bits[y+1][x] {
				if bits[y][x] != inverseColor {
					buf.WriteString(" ")
				} else {
					buf.WriteString("█")
				}
			} else {
				if bits[y][x] != inverseColor {
					buf.WriteString("▄")
				} else {
					buf.WriteString("▀")
				}
			}
		}
		buf.WriteString("\n")
	}
	// Special treatment for the last row if odd
	if len(bits)%2 == 1 {
		y := len(bits) - 1
		for x := range bits[y] {
			if bits[y][x] != inverseColor {
				buf.WriteString(" ")
			} else {
				buf.WriteString("▀")
			}
		}
		buf.WriteString("\n")
	}
	return buf.String()
}

// Encode a QR Code and return a raw PNG image
// Size is both the image width and height in pixels
// If size is too small then a larger image is silently returned
// Negative values for size cause a variable sized image to be returned:
// See the documentation for Image()
// To serve over HTTP, remember to send a Content-Type: image/png header
func Encode(content string, level RecoveryLevel, size int) ([]byte, error) {
	q, err := New(content, level)
	if err != nil {
		return nil, err
	}
	return q.PNG(size)
}

// Returns the maximum of a and b
func max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}
