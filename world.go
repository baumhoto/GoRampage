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
func (w *World) update() {
	w.player.position.Add(w.player.velocity)

	w.player.position.x = math.Mod(w.player.position.x, w.size.x)
	w.player.position.y = math.Mod(w.player.position.y, w.size.y)
}
