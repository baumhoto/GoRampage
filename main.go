package main

import (
	"encoding/json"
	"fmt"
	"github.com/hajimehoshi/ebiten/inpututil"
	"image"
	_ "image/png"
	"io/ioutil"
	"log"
	"math"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten"
)

type Game struct{}

var world World
var lastRenderFinishedTime time.Time
var renderer Renderer
var fullScreen bool

// TODO make frametime available from underlying window
const timeStep = 1.0 / 60.0
const maximumTimeStep = 1.0 / 20.0
const worldTimeStep = 1.0 / 120.0
const screenwidth = 320
const screenheight = 240
const screenscale = 2

func main() {
	world = NewWorld(loadMap())
	renderer = NewRenderer(screenwidth, screenheight)
	game := &Game{}
	ebiten.SetWindowSize(screenwidth*screenscale, screenheight*screenscale)
	ebiten.SetWindowTitle("GoRampage")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

// Update proceeds the game state.
// Update is called every tick (1/60 [s] by default).
func (g *Game) Update(screen *ebiten.Image) error {
	if inpututil.IsKeyJustReleased(ebiten.KeyEscape) {
		os.Exit(1)
	}

	worldSteps := math.Ceil(timeStep / worldTimeStep)
	for i := 0; i < int(worldSteps); i++ {
		world.update(float64(timeStep/worldSteps), GetInput())
	}

	return nil
}

// Draw draws the game screen.
// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (g *Game) Draw(screen *ebiten.Image) {
	//renderer.draw2d(world, screen)
	renderer.draw(world, screen)
	renderer.frameBuffer.resetFrameBuffer()
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return screenwidth, screenheight
}

func loadMap() Tilemap {
	file, _ := ioutil.ReadFile("map.json")

	data := Tilemap{}

	err := json.Unmarshal([]byte(file), &data)
	if err != nil {
		fmt.Printf("%v\n", err)
		return Tilemap{}
	}

	return data
}

func loadTextures() TextureManager {
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

	return TextureManager{textures: textures}
}

func GetInput() Input {
	inputVector := Vector{}
	velocity := float64(1)

	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		inputVector.y = velocity
	} else if ebiten.IsKeyPressed(ebiten.KeyUp) {
		inputVector.y = velocity * -1
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		inputVector.x = velocity * -1
	} else if ebiten.IsKeyPressed(ebiten.KeyRight) {
		inputVector.x = velocity
	}

	rotation := inputVector.x * world.player.turningSpeed * worldTimeStep

	return Input{
		speed:    -inputVector.y,
		rotation: NewRotation(math.Sin(rotation), math.Cos(rotation)),
	}
}
