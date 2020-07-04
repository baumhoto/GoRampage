package asset

import (
	"image"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	_map "github.com/baumhoto/GoRampage/engine/map"
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

func (tm TextureManager) GetTextureByKey(key string) Texture {
	return tm.textures[key]
}

func (tm *TextureManager) initAnimations() {
	tm.Animations = make(map[string]Animation, 7)
	tm.Animations[MonsterIdleAnimation] = Animation{
		frames:   []Texture{tm.textures[MonsterIdleTexture]},
		duration: 0,
	}

	tm.Animations[MonsterWalkAnimation] = Animation{
		frames: []Texture{
			tm.textures[MonsterWalkTexture1],
			tm.textures[MonsterIdleTexture],
			tm.textures[MonsterWalkTexture2],
			tm.textures[MonsterIdleTexture]},
		duration: 0.5,
	}

	tm.Animations[MonsterScratchAnimation] = Animation{
		frames: []Texture{
			tm.textures[MonsterScratchTexture1],
			tm.textures[MonsterScratchTexture2],
			tm.textures[MonsterScratchTexture3],
			tm.textures[MonsterScratchTexture4],
			tm.textures[MonsterScratchTexture5],
			tm.textures[MonsterScratchTexture6],
			tm.textures[MonsterScratchTexture7],
			tm.textures[MonsterScratchTexture8],
		},
		duration: 0.8,
	}

	tm.Animations[AnimationMonsterHurt] = Animation{
		frames: [] Texture{
			tm.textures[TextureMonsterHurt],
		},
		duration: 0.2,
	}

	tm.Animations[AnimationMonsterDeath] = Animation{
		frames: [] Texture{
			tm.textures[TextureMonsterHurt],
			tm.textures[TextureMonsterDeath1],
			tm.textures[TextureMonsterDeath2],
		},
		duration: 0.5,
	}

	tm.Animations[AnimationMonsterDead] = Animation{
		frames: [] Texture{
			tm.textures[TextureMonsterDead],
		},
		duration: 0.0,
	}

	tm.Animations[PistolIdleAnimation] = Animation{
		frames: []Texture{
			tm.textures[PistolIdleTexture],
		},
		duration: 0.0,
	}

	tm.Animations[PistolFireAnimation] = Animation{
		frames: []Texture{
			tm.textures[PistolFireTexture1],
			tm.textures[PistolFireTexture2],
			tm.textures[PistolFireTexture3],
			tm.textures[PistolFireTexture4],
		},
		duration: 0.5,
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
		_,_,_, a := img.At(0, 0).RGBA()

		if img != nil {
			texture := Texture{
				name:     textureName,
				category: GetTextureCategory(textureName),
				Image:    img,
				IsOpaque: uint8(a) == 255, // TODO not the best way to check
			}
			textures[textureId] = texture
		}

		//fmt.Printf("%v %v\n", textureName, uint8(a) ==255)
	}

	return textures
}
