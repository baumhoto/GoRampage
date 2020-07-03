package render

import (
	"image/color"
	"math"

	_asset "github.com/baumhoto/GoRampage/engine/asset"
	_core "github.com/baumhoto/GoRampage/engine/core"
	"github.com/hajimehoshi/ebiten"
)

// FrameBuffer stores the color for each pixel of the buffer
type FrameBuffer struct {
	width  int
	height int
	img    *ebiten.Image
}

// NewFrameBuffer creates a new FrameBuffer of given size with
// background-color c
func NewFrameBuffer(width int, height int, c color.Color) FrameBuffer {
	img, _ := ebiten.NewImage(width, height, ebiten.FilterDefault)
	fb := FrameBuffer{width, height, img}
	return fb
}

// SetColorAt sets the color for the pixel at locaton x, y to c
// in the Framebuffer
func (fb *FrameBuffer) SetColorAt(x int, y int, c color.Color) {
	if x >= 0 && y >= 0 && x < fb.width && y < fb.height {
		fb.img.Set(x, y, c) // = c
	}
}

// ColorAt retrieves the Color for a pixel at x, y in the FrameBuffer
func (fb FrameBuffer) ColorAt(x int, y int) color.Color {
	c := fb.img.At(x, y)
	return c
}

// Fill draws an Rect with Color c in the FrameBuffer
func (fb *FrameBuffer) Fill(rect _core.Rect, c color.Color) {
	for i := rect.Min.Y; i < rect.Max.Y; i++ {
		for j := rect.Min.X; j < rect.Max.X; j++ {
			fb.SetColorAt(int(j), int(i), c)
		}
	}
}

func (fb *FrameBuffer) drawColumn(sourceX int, source _asset.Texture, atPoint _core.Vector, height float64, windowHeight int) {
	start := int(atPoint.Y)
	end := int(math.Ceil(atPoint.Y + height))

	stepY := float64(source.Image.Bounds().Size().Y) / height
	for y := math.Max(0.0, float64(start)); y < math.Min(float64(windowHeight), float64(end)); y++ {
		sourceY := math.Max(0, y-atPoint.Y) * stepY
		sourceColor := source.GetColorAt(sourceX, int(sourceY))
		fb.blendPixel(int(atPoint.X), int(y), sourceColor)
	}
}

func (fb *FrameBuffer) drawImage(source _asset.Texture, at, size _core.Vector) {
	start := at.X
	end := at.X + size.X
	stepX := float64(source.Width()) / size.X
	for x := math.Max(0.0, start); x < math.Min(float64(fb.width), end); x++ {
		sourceX := (x - at.X) * stepX
		outPosition := _core.Vector{x, at.Y}
		fb.drawColumn(int(sourceX), source, outPosition, size.Y, fb.height)
	}
}

func (fb *FrameBuffer) DrawLine(from, to _core.Vector, color color.Color) {
	difference := _core.SubstractVectors(to, from)
	var stepCount int
	sign := -1.0
	var step _core.Vector
	if math.Abs(difference.X) > math.Abs(difference.Y) {
		stepCount = int(math.Ceil(math.Abs(difference.X)))
		if difference.X > 0 {
			sign = 1.0
		}
		step = _core.Vector{1, difference.Y / difference.X}
	} else {
		stepCount = int(math.Ceil(math.Abs(difference.Y)))
		if difference.Y > 0 {
			sign = 1.0
		}
		step = _core.Vector{difference.X / difference.Y, 1}
	}
	step.Multiply(sign)

	point := from

	for i := 0; i < stepCount; i++ {
		fb.SetColorAt(int(point.X), int(point.Y), color)
		point.Add(step)
	}
}

func (fb *FrameBuffer) resetFrameBuffer() {
	fb.img.Clear()
}

// ToRGBA converts the FrameBuffer into an Png-Image
// Buffer and returns an io.Reader
func (fb *FrameBuffer) ToImage() *ebiten.Image {
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
		uint8(alpha * 255)}

	//fmt.Printf("%v %v %v %v %v %v %v \n", r, g, b, a, effectColor, opacity, alpha)

	for y := 0; y < fb.height; y++ {
		for x := 0; x < fb.width; x++ {
			fb.blendPixel(x, y, effectColor)
		}
	}
}
