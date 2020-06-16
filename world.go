package main

import (
	"math"
)

// World is a container for the world
type World struct {
	size   Vector
	player Player
}

// NewWorld creates a new World.
func NewWorld(size Vector) World {
	return World{size, NewPlayer(Vector{4, 4})}
}

// update updates the World
func (w *World) update(timeStep int64) {
	velocity := MultiplyVector(w.player.velocity, float64(timeStep))
	w.player.position = AddVectors(w.player.position, velocity)

	w.player.position.x = math.Mod(w.player.position.x, w.size.x)
	w.player.position.y = math.Mod(w.player.position.y, w.size.y)
}
