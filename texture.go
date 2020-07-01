package main

import (
	"image"
	"image/color"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

var once sync.Once

var instance TextureManager

type TextureManager struct {
	textures   map[string]Texture
	animations map[string]Animation
}

func GetInstance() *TextureManager {
	once.Do(func() { // <-- atomic, does not allow repeating
		instance = TextureManager{loadTextures(), nil} // <-- thread safe
		instance.initAnimations()
	})
	return &instance
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

func (tm *TextureManager) initAnimations() {
	tm.animations = make(map[string]Animation, 2)
	tm.animations[MonsterIdleAnimation] = Animation{
		frames:   []Texture{tm.textures["5"]},
		duration: 0,
	}

	tm.animations[MonsterWalkAnimation] = Animation{
		frames: []Texture{
			tm.textures["5aw1"],
			tm.textures["5"],
			tm.textures["5aw2"],
			tm.textures["5"]},
		duration: 0.5,
	}

	tm.animations[MonsterScratchAnimation] = Animation{
		frames: []Texture{
			tm.textures["5as1"],
			tm.textures["5as2"],
			tm.textures["5as3"],
			tm.textures["5as4"],
			tm.textures["5as5"],
			tm.textures["5as6"],
			tm.textures["5as7"],
			tm.textures["5as8"],
		},
		duration: 0.8,
	}
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

func loadTextures() map[string]Texture {
	var textureFiles []string
	root := "textures" + string(os.PathSeparator)
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if strings.Contains(path, ".png") {
			textureFiles = append(textureFiles, path)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	textures := make(map[string]Texture)

	for _, fileName := range textureFiles {
		//fmt.Printf("%v\n", fileName)
		file, err := os.Open(fileName)
		if err != nil {
			panic(err)
		}

		img, _, err := image.Decode(file)
		if err != nil {
			panic(err)
		}
		textureNameString := strings.Split(strings.ToLower(fileName), string(os.PathSeparator))[1]
		textureNameParts := strings.Split(textureNameString, "_")
		textureName := textureNameParts[1]
		textureId := textureNameParts[0]

		if img != nil {
			texture := Texture{
				name:     textureName,
				category: GetTextureCategory(textureName),
				image:    img,
			}
			textures[textureId] = texture
		}
	}

	return textures
}
