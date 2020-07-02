package asset

import (
	"image"
	"image/color"
	"strings"
)

func GetTextureCategory(textureName string) int {
	if strings.Contains(textureName, "wall") {
		return 0
	}

	return 1
}

type Texture struct {
	name     string
	category int
	Image    image.Image
}

func (t Texture) GetColorAt(x, y int) color.Color {
	return t.Image.At(x, y)
}

func (t Texture) GetColorAtNormalized(x, y float64) color.Color {
	return t.Image.At(int(x*float64(t.Width())), int(y*float64(t.Height())))
}

func (t Texture) Width() int {
	return t.Image.Bounds().Size().X
}

func (t Texture) Height() int {
	return t.Image.Bounds().Size().Y
}
