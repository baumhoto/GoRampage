package input

import (
	_consts "github.com/baumhoto/go-rampage/engine/consts"
	_core "github.com/baumhoto/go-rampage/engine/core"
	"github.com/hajimehoshi/ebiten"
	"math"
	"time"
)

var lastFiredTime time.Time

type Input struct {
	Speed    float64
	Rotation _core.Rotation
	IsFiring bool
}

func GetInput(playerTurningSpeed float64, lastFrameTime float64) Input {
	inputVector := _core.Vector{}
	velocity := float64(1)

	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		inputVector.Y = velocity
	} else if ebiten.IsKeyPressed(ebiten.KeyUp) {
		inputVector.Y = velocity * -1
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		inputVector.X = velocity * -1
	} else if ebiten.IsKeyPressed(ebiten.KeyRight) {
		inputVector.X = velocity
	}

	fire := false

	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		if !lastFiredTime.IsZero() {
			test := time.Since(lastFiredTime)
			if test.Seconds() > lastFrameTime {
				//ßfmt.Printf("%v %v\n", test.Seconds(), lastFrameTime)
				fire = true
			}
		} else {
			fire = true
		}
		lastFiredTime = time.Now()
	}

	rotation := inputVector.X * playerTurningSpeed * _consts.WORLD_TIMESTEP

	return Input{
		Speed:    -inputVector.Y,
		Rotation: _core.NewRotation(math.Sin(rotation), math.Cos(rotation)),
		IsFiring: fire,
	}
}
