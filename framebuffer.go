package main

import (
	"github.com/hajimehoshi/ebiten"
	"image/color"
	"math"
)

// FrameBuffer stores the color for each pixel of the buffer
type FrameBuffer struct {
	pixels          []color.Color
	width           int
	height          int
	backgroundColor color.Color
	img             *ebiten.Image
}

// NewFrameBuffer creates a new FrameBuffer of given size with
// background-color c
func NewFrameBuffer(width int, height int, c color.Color) FrameBuffer {
	img, _ := ebiten.NewImage(width, height, ebiten.FilterDefault)

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

func (fb *FrameBuffer) drawColumn(sourceX int, source Texture, atPoint Vector, height float64, windowHeight int) {
	start := int(atPoint.y)
	end := int(math.Ceil(atPoint.y + height))

	stepY := float64(source.image.Bounds().Size().Y) / height
	for y := math.Max(0.0, float64(start)); y < math.Min(float64(windowHeight), float64(end)); y++ {
		sourceY := math.Max(0, y-atPoint.y) * stepY
		sourceColor := source.GetColorAt(sourceX, int(sourceY))
		fb.blendPixel(int(atPoint.x), int(y), sourceColor)
	}
}

func (fb *FrameBuffer) DrawLine(from, to Vector, color color.Color) {
	difference := SubstractVectors(to, from)
	var stepCount int
	sign := -1.0
	var step Vector
	if math.Abs(difference.x) > math.Abs(difference.y) {
		stepCount = int(math.Ceil(math.Abs(difference.x)))
		if difference.x > 0 {
			sign = 1.0
		}
		step = Vector{1, difference.y / difference.x}
	} else {
		stepCount = int(math.Ceil(math.Abs(difference.y)))
		if difference.y > 0 {
			sign = 1.0
		}
		step = Vector{difference.x / difference.y, 1}
	}
	step.Multiply(sign)

	point := from

	for i := 0; i < stepCount; i++ {
		fb.SetColorAt(int(point.x), int(point.y), color)
		point.Add(step)
	}
}

func (fb *FrameBuffer) resetFrameBuffer() {
	for index := range fb.pixels {
		fb.pixels[index] = fb.backgroundColor
	}
}

// ToRGBA converts the FrameBuffer into an Png-Image
// Buffer and returns an io.Reader
func (fb *FrameBuffer) ToImage() *ebiten.Image {
	// Set color for each pixel.
	for x := 0; x < fb.width; x++ {
		for y := 0; y < fb.height; y++ {
			fb.img.Set(x, y, fb.ColorAt(x, y))
		}
	}

	return fb.img
}
func (fb *FrameBuffer) blendPixel(x, y int, newColor color.Color) {
	oldR, oldG, oldB, _ := fb.ColorAt(x, y).RGBA()
	newR, newG, newB, newA := newColor.RGBA()
	inverseAlpha := 1.0 - float64(uint8(newA))/255.0
	fb.SetColorAt(x, y, color.RGBA{
		R: uint8(float64(uint8(oldR))*inverseAlpha) + uint8(newR),
		G: uint8(float64(uint8(oldG))*inverseAlpha) + uint8(newG),
		B: uint8(float64(uint8(oldB))*inverseAlpha) + uint8(newB),
		A: uint8(255),
	})
}

func (fb *FrameBuffer) tint(tintColor color.Color, opacity float64) {
	r, g, b, a := tintColor.RGBA()
	alpha := math.Min(1.0, math.Max(0.0, float64(uint8(a))/255*opacity))
	effectColor := color.RGBA{
		uint8(float64(uint8(r)) * alpha),
		uint8(float64(uint8(g)) * alpha),
		uint8(float64(uint8(b)) * alpha),
		uint8(uint8(alpha) * 255)}

	//fmt.Printf("%v %v %v %v %v %v %v \n", r, g, b, a, effectColor, opacity, alpha)

	for y := 0; y < fb.height; y++ {
		for x := 0; x < fb.width; x++ {
			fb.blendPixel(x, y, effectColor)
		}
	}
}
