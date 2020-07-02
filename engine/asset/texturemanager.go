package asset

import (
	_map "github.com/baumhoto/go-rampage/engine/map"
	"image"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type TextureManager struct {
	textures   map[string]Texture
	Animations map[string]Animation
}

func NewTextureManager() *TextureManager {
	instance := TextureManager{loadTextures(), nil} // <-- thread safe
	instance.initAnimations()
	return &instance
}

func (tm TextureManager) GetWallTextureByTile(tile _map.Tile, isVertical bool) Texture {
	searchId := strconv.Itoa(int(tile))
	if isVertical {
		searchId = searchId + "v"
	}
	return tm.textures[searchId]
}

func (tm TextureManager) GetFloorCeilingTextureByTile(tile _map.Tile, isCeiling bool) Texture {
	searchId := strconv.Itoa(int(tile))

	if tile == 4 && isCeiling {
		searchId = "0c"
	} else if isCeiling {
		searchId = searchId + "c"
	} else {
		searchId = searchId + "f"
	}

	result := tm.textures[searchId]

	if result.Image == nil {
		result = tm.GetWallTextureByTile(tile, isCeiling)
	}

	return result
}

func (tm *TextureManager) initAnimations() {
	tm.Animations = make(map[string]Animation, 2)
	tm.Animations[MonsterIdleAnimation] = Animation{
		frames:   []Texture{tm.textures["5"]},
		duration: 0,
	}

	tm.Animations[MonsterWalkAnimation] = Animation{
		frames: []Texture{
			tm.textures["5aw1"],
			tm.textures["5"],
			tm.textures["5aw2"],
			tm.textures["5"]},
		duration: 0.5,
	}

	tm.Animations[MonsterScratchAnimation] = Animation{
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
				Image:    img,
			}
			textures[textureId] = texture
		}
	}

	return textures
}
