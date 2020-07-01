package main

import "math"

// Player is a container for a Player
type Player struct {
	speed        float64
	turningSpeed float64
	radius       float64
	position     Vector
	velocity     Vector
	direction    Vector
}

// NewPlayer creates a new Player
func NewPlayer(position Vector) Player {
	return Player{2, math.Pi, 0.25, position, Vector{0, 0}, Vector{1, 0}}
}

func (p Player) rect() Rect {
	return rect(p.radius, p.position)
}

func (p Player) intersection(tileMap Tilemap) (bool, Vector) {
	return intersection(p.rect(), tileMap)
}
