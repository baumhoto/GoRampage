package main

import (
	"image"
	"image/color"
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

type Texture struct {
	category int
	image    image.Image
}

func (t Texture) GetColorAt(x, y int) color.Color {
	return t.image.At(x, y)
}
