package entity

import (
	_core "github.com/baumhoto/go-rampage/engine/core"
	_map "github.com/baumhoto/go-rampage/engine/map"
	"math"
)

// Player is a container for a Player
type Player struct {
	speed        float64
	TurningSpeed float64
	radius       float64
	Position     _core.Vector
	velocity     _core.Vector
	Direction    _core.Vector
	health       float64
}

// NewPlayer creates a new Player
func NewPlayer(position _core.Vector) Player {
	return Player{2, math.Pi, 0.25, position,
		_core.Vector{0, 0}, _core.Vector{1, 0}, 100}
}

func (p Player) Rect() _core.Rect {
	return rect(p.radius, p.Position)
}

func (p Player) intersection(tileMap _map.Tilemap) (bool, _core.Vector) {
	return intersection(p.Rect(), tileMap)
}

func (p Player) isDead() bool {
	return p.health <= 0
}
