package main

import (
	"image"
	"image/color"
)

// FrameBuffer stores the color for each pixel of the buffer
type FrameBuffer struct {
	pixels          []color.Color
	width           int
	height          int
	backgroundColor color.Color
	img             *image.RGBA
}

// NewFrameBuffer creates a new FrameBuffer of given size with
// background-color c
func NewFrameBuffer(width int, height int, c color.Color) FrameBuffer {
	upLeft := image.Point{0, 0}
	lowRight := image.Point{width, height}
	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})
	pixels := make([]color.Color, width*height)
	fb := FrameBuffer{pixels, width, height, c, img}
	fb.resetFrameBuffer()

	return fb
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

func (fb *FrameBuffer) resetFrameBuffer() {
	for index := range fb.pixels {
		fb.pixels[index] = fb.backgroundColor
	}
}

// ToRGBA converts the FrameBuffer into an Png-Image
// Buffer and returns an io.Reader
func (fb *FrameBuffer) ToRGBA() *image.RGBA {
	// Set color for each pixel.
	for x := 0; x < fb.width; x++ {
		for y := 0; y < fb.height; y++ {
			fb.img.Set(x, y, fb.ColorAt(x, y))
		}
	}

	return fb.img
	//outputFile, _ := os.Create("test.png")
	//png.Encode(outputFile, img)
	//outputFile.Close()

	// var buf bytes.Buffer

	// err := png.Encode(&buf, fb.img)
	// if err != nil {
	// 	fmt.Printf("Encode error")
	// }

	// return bytes.NewReader(buf.Bytes())
}
