package main

import (
	"math"
)

// World is a container for the world
type World struct {
	worldmap Tilemap
	player   Player
}

// NewWorld creates a new World.
func NewWorld(worldmap Tilemap) World {
	var player Player
	for y := 0; y < worldmap.Height; y++ {
		for x := 0; x < worldmap.Width; x++ {
			position := Vector{float64(x) + 0.5, float64(y) + 0.5}
			thing := worldmap.Things[y*worldmap.Width+x]
			switch thing {
			case 0:
				break
			case 1:
				player = NewPlayer(position)
			}
		}
	}
	return World{worldmap, player}
}

// update updates the World
func (w *World) update(input Input) {
	oldPosition := world.player.position
	input.velocity.Multiply(w.player.speed)
	w.player.velocity = input.velocity
	w.player.position.Add(w.player.velocity)

	if w.player.isIntersecting(w.worldmap) {
		w.player.position = oldPosition
	}

	w.player.position.x = math.Mod(w.player.position.x, float64(w.worldmap.Width))
	w.player.position.y = math.Mod(w.player.position.y, float64(w.worldmap.Height))
}
