package entity

import (
	_common "github.com/baumhoto/go-rampage/engine/common"
	_map "github.com/baumhoto/go-rampage/engine/map"
	"math"
)

// Player is a container for a Player
type Player struct {
	speed        float64
	TurningSpeed float64
	radius       float64
	Position     _common.Vector
	velocity     _common.Vector
	Direction    _common.Vector
	health       float64
}

// NewPlayer creates a new Player
func NewPlayer(position _common.Vector) Player {
	return Player{2, math.Pi, 0.25, position,
		_common.Vector{0, 0}, _common.Vector{1, 0}, 100}
}

func (p Player) Rect() _common.Rect {
	return rect(p.radius, p.Position)
}

func (p Player) intersection(tileMap _map.Tilemap) (bool, _common.Vector) {
	return intersection(p.Rect(), tileMap)
}

func (p Player) isDead() bool {
	return p.health <= 0
}
