package media

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	_ "image/png"
	"io"
)

const thumbnailSize = 150

func thumbnail(file io.Reader) []byte {

	// Decode the image
	img, _, err := image.Decode(file)
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}

	// Create thumbnail (resize the image)
	b := img.Bounds()
	originalWidth := b.Dx()
	originalHeight := b.Dy()

	var scaleFactor float64
	if originalWidth > originalHeight {
		scaleFactor = float64(thumbnailSize) / float64(originalWidth)
	} else {
		scaleFactor = float64(thumbnailSize) / float64(originalHeight)
	}

	// Calculate the new dimensions based on the scaling factor
	newWidth := int(float64(originalWidth) * scaleFactor)
	newHeight := int(float64(originalHeight) * scaleFactor)

	thumbnail := resize(img, newWidth, newHeight)

	var buf bytes.Buffer
	jpeg.Encode(&buf, thumbnail, &jpeg.Options{
		Quality: 90,
	})
	return buf.Bytes()
}
