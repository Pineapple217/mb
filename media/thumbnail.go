package media

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	_ "image/png"
	"io"
	"math"
	"runtime"
	"sync"
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

// resize resizes an image to the specified width and height
// func resize(img image.Image, width, height int) image.Image {
// 	// Calculate new dimensions
// 	bounds := img.Bounds()
// 	newWidth := bounds.Max.X
// 	newHeight := bounds.Max.Y
// 	if newWidth > width {
// 		newHeight = newHeight * width / newWidth
// 		newWidth = width
// 	}
// 	if newHeight > height {
// 		newWidth = newWidth * height / newHeight
// 		newHeight = height
// 	}

// 	// Create a new image with the new dimensions
// 	newImage := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))

// 	// Resize the image using nearest-neighbor interpolation
// 	for y := 0; y < newHeight; y++ {
// 		for x := 0; x < newWidth; x++ {
// 			newImage.Set(x, y, img.At(x*bounds.Max.X/newWidth, y*bounds.Max.Y/newHeight))
// 		}
// 	}

// 	return newImage
// }

func resize(img image.Image, width, height int) *image.RGBA {
	if width <= 0 || height <= 0 || img.Bounds().Empty() {
		return image.NewRGBA(image.Rect(0, 0, 0, 0))
	}

	src := AsShallowRGBA(img)

	var dst *image.RGBA

	dst = resampleHorizontal(src, width)
	dst = resampleVertical(dst, height)

	return dst
}

func filter(x float64) float64 {
	x = math.Abs(x)
	if x == 0 {
		return 1.0
	} else if x < 3.0 {
		return (3.0 * math.Sin(math.Pi*x) * math.Sin(math.Pi*(x/3.0))) / (math.Pi * math.Pi * x * x)
	}
	return 0.0
}

func resampleHorizontal(src *image.RGBA, width int) *image.RGBA {
	srcWidth, srcHeight := src.Bounds().Dx(), src.Bounds().Dy()
	srcStride := src.Stride

	delta := float64(srcWidth) / float64(width)
	// Scale must be at least 1. Special case for image size reduction filter radius.
	scale := math.Max(delta, 1.0)

	dst := image.NewRGBA(image.Rect(0, 0, width, srcHeight))
	dstStride := dst.Stride

	filterRadius := math.Ceil(scale * 3.0)

	Line(srcHeight, func(start, end int) {
		for y := start; y < end; y++ {
			for x := 0; x < width; x++ {
				// value of x from src
				ix := (float64(x)+0.5)*delta - 0.5
				istart, iend := int(ix-filterRadius+0.5), int(ix+filterRadius)

				if istart < 0 {
					istart = 0
				}
				if iend >= srcWidth {
					iend = srcWidth - 1
				}

				var r, g, b, a float64
				var sum float64
				for kx := istart; kx <= iend; kx++ {

					srcPos := y*srcStride + kx*4
					// normalize the sample position to be evaluated by the filter
					normPos := (float64(kx) - ix) / scale
					fValue := filter(normPos)

					r += float64(src.Pix[srcPos+0]) * fValue
					g += float64(src.Pix[srcPos+1]) * fValue
					b += float64(src.Pix[srcPos+2]) * fValue
					a += float64(src.Pix[srcPos+3]) * fValue
					sum += fValue
				}

				dstPos := y*dstStride + x*4
				dst.Pix[dstPos+0] = uint8(Clamp((r/sum)+0.5, 0, 255))
				dst.Pix[dstPos+1] = uint8(Clamp((g/sum)+0.5, 0, 255))
				dst.Pix[dstPos+2] = uint8(Clamp((b/sum)+0.5, 0, 255))
				dst.Pix[dstPos+3] = uint8(Clamp((a/sum)+0.5, 0, 255))
			}
		}
	})

	return dst
}

func resampleVertical(src *image.RGBA, height int) *image.RGBA {
	srcWidth, srcHeight := src.Bounds().Dx(), src.Bounds().Dy()
	srcStride := src.Stride

	delta := float64(srcHeight) / float64(height)
	scale := math.Max(delta, 1.0)

	dst := image.NewRGBA(image.Rect(0, 0, srcWidth, height))
	dstStride := dst.Stride

	filterRadius := math.Ceil(scale * 3.0)

	Line(height, func(start, end int) {
		for y := start; y < end; y++ {
			iy := (float64(y)+0.5)*delta - 0.5

			istart, iend := int(iy-filterRadius+0.5), int(iy+filterRadius)

			if istart < 0 {
				istart = 0
			}
			if iend >= srcHeight {
				iend = srcHeight - 1
			}

			for x := 0; x < srcWidth; x++ {
				var r, g, b, a float64
				var sum float64
				for ky := istart; ky <= iend; ky++ {

					srcPos := ky*srcStride + x*4
					normPos := (float64(ky) - iy) / scale
					fValue := filter(normPos)

					r += float64(src.Pix[srcPos+0]) * fValue
					g += float64(src.Pix[srcPos+1]) * fValue
					b += float64(src.Pix[srcPos+2]) * fValue
					a += float64(src.Pix[srcPos+3]) * fValue
					sum += fValue
				}

				dstPos := y*dstStride + x*4
				dst.Pix[dstPos+0] = uint8(Clamp((r/sum)+0.5, 0, 255))
				dst.Pix[dstPos+1] = uint8(Clamp((g/sum)+0.5, 0, 255))
				dst.Pix[dstPos+2] = uint8(Clamp((b/sum)+0.5, 0, 255))
				dst.Pix[dstPos+3] = uint8(Clamp((a/sum)+0.5, 0, 255))
			}
		}
	})

	return dst
}

func AsShallowRGBA(src image.Image) *image.RGBA {
	if rgba, ok := src.(*image.RGBA); ok {
		return rgba
	}
	return AsRGBA(src)
}

func AsRGBA(src image.Image) *image.RGBA {
	bounds := src.Bounds()
	img := image.NewRGBA(bounds)
	draw.Draw(img, bounds, src, bounds.Min, draw.Src)
	return img
}

func Line(length int, fn func(start, end int)) {
	procs := runtime.GOMAXPROCS(0)
	counter := length
	partSize := length / procs
	if procs <= 1 || partSize <= procs {
		fn(0, length)
	} else {
		var wg sync.WaitGroup
		for counter > 0 {
			start := counter - partSize
			end := counter
			if start < 0 {
				start = 0
			}
			counter -= partSize
			wg.Add(1)
			go func() {
				defer wg.Done()
				fn(start, end)
			}()
		}

		wg.Wait()
	}
}

func Clamp(value, min, max float64) float64 {
	if value > max {
		return max
	}
	if value < min {
		return min
	}
	return value
}

// MIT License

// Copyright (c) 2021 Anthony Najjar Simon

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
