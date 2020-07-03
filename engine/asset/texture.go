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

const (
	CeilingTexture         = "0c"
	FloorTexture           = "0f"
	WallTexture            = "1"
	Wall2Texture           = "1v"
	CrackWallTexture       = "2"
	CrackWall2Texture      = "2v"
	SlimeWallTexture       = "3"
	SlimeWall2Texture      = "3v"
	CrackFloorTexture      = "4f"
	MonsterIdleTexture     = "5"
	MonsterWalkTexture1    = "5aw1"
	MonsterWalkTexture2    = "5aw2"
	MonsterScratchTexture1 = "5as1"
	MonsterScratchTexture2 = "5as2"
	MonsterScratchTexture3 = "5as3"
	MonsterScratchTexture4 = "5as4"
	MonsterScratchTexture5 = "5as5"
	MonsterScratchTexture6 = "5as6"
	MonsterScratchTexture7 = "5as7"
	MonsterScratchTexture8 = "5as8"
	TextureMonsterHurt = "5h"
	TextureMonsterDeath1 = "5ad1"
	TextureMonsterDeath2 = "5ad2"
	TextureMonsterDead = "5d"
	PistolIdleTexture      = "6"
	PistolFireTexture1     = "6af1"
	PistolFireTexture2     = "6af2"
	PistolFireTexture3     = "6af3"
	PistolFireTexture4     = "6af4"
)
