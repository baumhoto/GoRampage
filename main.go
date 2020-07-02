package main

import (
	_asset "github.com/baumhoto/go-rampage/engine/asset"
	_common "github.com/baumhoto/go-rampage/engine/common"
	_entity "github.com/baumhoto/go-rampage/engine/entity"
	_input "github.com/baumhoto/go-rampage/engine/input"
	_map "github.com/baumhoto/go-rampage/engine/map"
	_render "github.com/baumhoto/go-rampage/engine/render"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
	_ "image/png"
	"log"
	"math"
	"os"
	"time"
)

type Game struct{}

var world _entity.World
var lastRenderFinishedTime time.Time
var renderer _render.Renderer
var fullScreen bool
var textureManager *_asset.TextureManager
var pause bool

// TODO make frametime available from underlying window

func main() {
	textureManager = _asset.GetInstance()
	world = _entity.NewWorld(_map.LoadMap())
	renderer = _render.NewRenderer(_common.SCREEN_WIDTH, _common.SCREEN_HEIGHT, *textureManager)
	game := &Game{}
	ebiten.SetWindowSize(_common.SCREEN_WIDTH*_common.SCREEN_SCALE, _common.SCREEN_HEIGHT*_common.SCREEN_SCALE)
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

	if inpututil.IsKeyJustPressed(ebiten.KeyP) {
		pause = !pause
	}

	if !pause {
		worldSteps := math.Ceil(_common.TIMESTEP / _common.WORLD_TIMESTEP)
		for i := 0; i < int(worldSteps); i++ {
			world.Update(_common.TIMESTEP/worldSteps, _input.GetInput(world.Player.TurningSpeed))
		}
	}

	return nil
}

// Draw draws the game screen.
// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (g *Game) Draw(screen *ebiten.Image) {
	if !pause {
		renderer.ResetFrameBuffer()
	}
	//renderer.Draw2d(world, screen)
	renderer.Draw(world, screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return _common.SCREEN_WIDTH, _common.SCREEN_HEIGHT
}
