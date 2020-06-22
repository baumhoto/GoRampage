package main

import (
	"image"
	"image/color"
	"strconv"
	"strings"
)

type TextureManager struct {
	textures map[string]Texture
}

func (tm TextureManager) GetWallTextureByTile(tile Tile, isVertical bool) Texture {
	searchId := strconv.Itoa(int(tile))
	if isVertical {
		searchId = searchId + "v"
	}
	return tm.textures[searchId]
}

func (tm TextureManager) GetFloorCeilingTextureByTile(tile Tile, isCeiling bool) Texture {
	searchId := strconv.Itoa(int(tile))

	if tile == 4 && isCeiling {
		searchId = "0c"
	} else if isCeiling {
		searchId = searchId + "c"
	} else {
		searchId = searchId + "f"
	}

	result := tm.textures[searchId]

	if result.image == nil {
		result = tm.GetWallTextureByTile(tile, isCeiling)
	}

	return result
}

func GetTextureCategory(textureName string) int {
	if strings.Contains(textureName, "wall") {
		return 0
	}

	return 1
}

type Texture struct {
	name     string
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
