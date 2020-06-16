package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
)

type FrameBuffer struct {
	pixels []color.Color
	width  int
	height int
}

func NewFrameBuffer(width int, height int, c color.Color) FrameBuffer {
	pixels := make([]color.Color, width*height)
	for index := range pixels {
		pixels[index] = c
	}
	return FrameBuffer{pixels, width, height}
}

func (fb *FrameBuffer) SetColorAt(x int, y int, c color.Color) {
	if x >= 0 && y >= 0 && x < fb.width && y < fb.height {
		fb.pixels[y*fb.width+x] = c
	}
}

func (fb FrameBuffer) ColorAt(x int, y int) color.Color {
	return fb.pixels[y*fb.width+x]
}

func (fb *FrameBuffer) Fill(rect Rect, c color.Color) {
	for i := rect.min.y; i < rect.max.y; i++ {
		for j := rect.min.x; j < rect.max.x; j++ {
			fb.SetColorAt(int(i), int(j), c)
		}
	}
}

func (fb FrameBuffer) ToTexture() io.Reader {
	upLeft := image.Point{0, 0}
	lowRight := image.Point{fb.width, fb.height}

	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})

	// Set color for each pixel.
	for x := 0; x < fb.width; x++ {
		for y := 0; y < fb.height; y++ {
			switch {
			// case x < fb.width/2 && y < fb.height/2: // upper left quadrant
			// 	img.Set(x, y, cyan)
			// case x >= fb.width/2 && y >= fb.height/2: // lower right quadrant
			// 	img.Set(x, y, color.White)
			default:
				img.Set(x, y, fb.ColorAt(x, y))
				//fmt.Printf("%v\n", fb.ColorAt(x, y))
				// Use zero value.
			}
		}
	}

	var buf bytes.Buffer

	//f, _ := os.Create("image.png")
	//png.Encode(f, img)

	err := png.Encode(&buf, img)
	if err != nil {
		fmt.Printf("Encode error")
	}

	// temp := bytes.NewReader(buf.Bytes())
	// temp.Seek(0, 0)
	// fmt.Printf("%v\n", temp.Len())

	return bytes.NewReader(buf.Bytes())
}
