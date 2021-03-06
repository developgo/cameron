package cameron

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"image"
	"image/color"
)

// Identicon returns an identicon avatar based on the data with the length and
// the blockLength. Same parameters, same result.
func Identicon(data []byte, length, blockLength int) image.Image {
	b := sha256.Sum256(data)

	img := image.NewPaletted(
		image.Rect(0, 0, length, length),
		color.Palette{
			color.NRGBA{
				R: 0xff ^ b[0],
				G: 0xff ^ b[1],
				B: 0xff ^ b[2],
				A: 0xff,
			},
			color.NRGBA{
				R: b[0],
				G: b[1],
				B: b[2],
				A: 0xff,
			},
		},
	)

	if blockLength > length {
		blockLength = length
	}

	columnsCount := length / blockLength
	padding := blockLength / 2
	if length%blockLength != 0 {
		padding = (length - blockLength*columnsCount) / 2
	} else if columnsCount > 1 {
		columnsCount--
	} else {
		padding = 0
	}

	var (
		filled   = columnsCount == 1
		sequence = binary.BigEndian.Uint64(b[:])
		pixels   = bytes.Repeat([]byte{1}, blockLength)
	)

	for i, ri, ci := 0, 0, 0; i < columnsCount*(columnsCount+1)/2; i++ {
		if filled || sequence>>uint(i%64)&1 == 1 {
			for i := 0; i < blockLength; i++ {
				x := padding + ri*blockLength
				y := padding + ci*blockLength + i
				copy(img.Pix[img.PixOffset(x, y):], pixels)

				x = padding + (columnsCount-1-ri)*blockLength
				copy(img.Pix[img.PixOffset(x, y):], pixels)
			}
		}

		ci++
		if ci == columnsCount {
			ci = 0
			ri++
		}
	}

	return img
}
