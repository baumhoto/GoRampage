package main

import (
	"errors"
	"flag"
	"fmt"
	_asset "github.com/baumhoto/GoRampage/engine/asset"
	_ "image/png"
	"log"
	"math"
	"os"
	"runtime/pprof"
	"time"

	_consts "github.com/baumhoto/GoRampage/engine/consts"
	_entity "github.com/baumhoto/GoRampage/engine/entity"
	_input "github.com/baumhoto/GoRampage/engine/input"
	_render "github.com/baumhoto/GoRampage/engine/render"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
)

type Game struct {
}

var world _entity.World
var renderer _render.Renderer
var textureManger *_asset.TextureManager
var fullScreen bool
var pause bool
var lastFrameTime float64
var lastTime time.Time
var count = 0
var profile = false
var lastStatsTime = time.Now()
var showFps = false
var fixedUpdateCyclesCount = 0

func main() {

	var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")
	var fps = flag.Int("showFps", 0, "Show FPS stats `int`")
	var updateCycles = flag.Int("fixedUpdateCycles", 0, "Set the number of update cycles after which program execution ends.")

	flag.Parse()

	if *fps == 1 {
		showFps = true
	}

	if *updateCycles > 0 {
		fixedUpdateCyclesCount = *updateCycles
	}

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		defer f.Close() // error handling omitted for example
		profile = true
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	world = _entity.NewWorld()
	renderer = _render.NewRenderer()
	textureManger = _asset.NewTextureManager()
	game := &Game{}
	ebiten.SetWindowSize(_consts.SCREEN_WIDTH*_consts.SCREEN_SCALE, _consts.SCREEN_HEIGHT*_consts.SCREEN_SCALE)
	ebiten.SetWindowTitle("GoRampage")
	//ebiten.SetVsyncEnabled(true)
	if err := ebiten.RunGame(game); err != nil {
		log.Print(err)
	}
}

// Update proceeds the game state.
// Update is called every tick (1/60 [s] by default).
func (g *Game) Update(screen *ebiten.Image) error {
	if inpututil.IsKeyJustReleased(ebiten.KeyEscape) {
		os.Exit(1)
	}

	if inpututil.IsKeyJustReleased(ebiten.KeyP) {
		pause = !pause
	}
	if inpututil.IsKeyJustReleased(ebiten.KeyR) {
		world.Reset()
	}

	if !pause {
		worldSteps := math.Ceil(_consts.TIMESTEP / _consts.WORLD_TIMESTEP)
		for i := 0; i < int(worldSteps); i++ {
			world.Update(_consts.TIMESTEP/worldSteps, _input.GetInput(world.Player.TurningSpeed, lastFrameTime), textureManger)
		}
	}

	if fixedUpdateCyclesCount > 0 {
		count++
		if count >= fixedUpdateCyclesCount {
			return errors.New("fixed update cycles reached. Exiting", )
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
	renderer.Draw(world, screen, textureManger)

	if !lastTime.IsZero() {
		lastFrameTime = time.Since(lastTime).Seconds()
		if showFps && time.Since(lastStatsTime).Seconds() >= 3 {
			fmt.Printf("%v %v\n", lastFrameTime, ebiten.CurrentFPS())
			lastStatsTime = time.Now()
		}
	}
	lastTime = time.Now()
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return _consts.SCREEN_WIDTH, _consts.SCREEN_HEIGHT
}
