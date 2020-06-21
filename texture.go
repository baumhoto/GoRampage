package main

import (
	"image"
	"image/color"
	"strings"
)

const (
	Category_Wall = 0
)

type TextureManager struct {
	textures map[string]Texture
}

func (tm TextureManager) GetTextureByName(name string) Texture {
	return tm.textures[name]
}

func GetTextureCategory(textureName string) int {
	if strings.Contains(textureName, "wall") {
		return 0
	} else if strings.Contains(textureName, "floor") {
		return 1
	} else if strings.Contains(textureName, "ceiling") {
		return 1
	}

	return 0
}

type Texture struct {
	category int
	image    image.Image
}

func (t Texture) GetColorAt(x, y int) color.Color {
	return t.image.At(x, y)
}

func (t Texture) GetColorAtNormalized(x, y float64) color.Color {
	return t.image.At(int(x*float64(t.Width())), int(y*float64(t.Height())))
}

func (t Texture) Width() int {
	return t.image.Bounds().Size().X
}

func (t Texture) Height() int {
	return t.image.Bounds().Size().Y
}
