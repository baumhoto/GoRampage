package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
)

// FrameBuffer stores the color for each pixel of the buffer
type FrameBuffer struct {
	pixels []color.Color
	width  int
	height int
}

// NewFrameBuffer creates a new FrameBuffer of given size with
// background-color c
func NewFrameBuffer(width int, height int, c color.Color) FrameBuffer {
	pixels := make([]color.Color, width*height)
	for index := range pixels {
		pixels[index] = c
	}
	return FrameBuffer{pixels, width, height}
}

// SetColorAt sets the color for the pixel at locaton x, y to c
// in the Framebuffer
func (fb *FrameBuffer) SetColorAt(x int, y int, c color.Color) {
	if x >= 0 && y >= 0 && x < fb.width && y < fb.height {
		fb.pixels[y*fb.width+x] = c
	}
}

// ColorAt retrieves the Color for a pixel at x, y in the FrameBuffer
func (fb FrameBuffer) ColorAt(x int, y int) color.Color {
	c := fb.pixels[y*fb.width+x]
	return c
}

// Fill draws an Rect with Color c in the FrameBuffer
func (fb *FrameBuffer) Fill(rect Rect, c color.Color) {
	for i := rect.min.y; i < rect.max.y; i++ {
		for j := rect.min.x; j < rect.max.x; j++ {
			fb.SetColorAt(int(j), int(i), c)
		}
	}
}

// ToImageReader converts the FrameBuffer into an Png-Image
// Buffer and returns an io.Reader
func (fb FrameBuffer) ToImageReader() io.Reader {
	upLeft := image.Point{0, 0}
	lowRight := image.Point{fb.width, fb.height}

	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})

	// Set color for each pixel.
	for x := 0; x < fb.width; x++ {
		for y := 0; y < fb.height; y++ {
			img.Set(x, y, fb.ColorAt(x, y))
		}
	}

	//outputFile, _ := os.Create("test.png")
	//png.Encode(outputFile, img)
	//outputFile.Close()

	var buf bytes.Buffer

	err := png.Encode(&buf, img)
	if err != nil {
		fmt.Printf("Encode error")
	}

	return bytes.NewReader(buf.Bytes())
}
