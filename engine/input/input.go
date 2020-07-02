package input

import (
	_common "github.com/baumhoto/go-rampage/engine/common"
	"github.com/hajimehoshi/ebiten"
	"math"
)

type Input struct {
	Speed    float64
	Rotation _common.Rotation
}

func GetInput(playerTurningSpeed float64) Input {
	inputVector := _common.Vector{}
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

	rotation := inputVector.X * playerTurningSpeed * _common.WORLD_TIMESTEP

	return Input{
		Speed:    -inputVector.Y,
		Rotation: _common.NewRotation(math.Sin(rotation), math.Cos(rotation)),
	}
}
