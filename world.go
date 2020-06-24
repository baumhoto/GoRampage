package main

import (
	"math"
)

// World is a container for the world
type World struct {
	worldmap Tilemap
	player   Player
	monsters []Monster
}

// NewWorld creates a new World.
func NewWorld(worldmap Tilemap) World {
	var player Player
	var monsters []Monster
	for y := 0; y < worldmap.Height; y++ {
		for x := 0; x < worldmap.Width; x++ {
			position := Vector{float64(x) + 0.5, float64(y) + 0.5}
			thing := worldmap.Things[y*worldmap.Width+x]
			switch thing {
			case 0:
				break
			case 1:
				player = NewPlayer(position)
				break
			case 2:
				monsters = append(monsters, Monster{position: position})
				break
			}
		}
	}
	return World{worldmap, player, monsters}
}

// update updates the World
func (w *World) update(timeStep float64, input Input) {
	w.player.direction = w.player.direction.rotated(input.rotation)
	w.player.velocity = MultiplyVector(w.player.direction, input.speed*w.player.speed)
	w.player.velocity.Multiply(timeStep)
	w.player.position.Add(w.player.velocity)

	for {
		if ok, intersection := w.player.intersection(w.worldmap); ok {
			intersection.Multiply(-0.01)
			w.player.position.Substract(intersection)
		} else {
			break
		}
	}

	w.player.position.x = math.Mod(w.player.position.x, float64(w.worldmap.Width))
	w.player.position.y = math.Mod(w.player.position.y, float64(w.worldmap.Height))
}

func (w World) sprites() []Billboard {
	spritePlane := w.player.direction.orthogonal()
	var result []Billboard
	for _, monster := range w.monsters {
		start := DivideVector(spritePlane, 2)
		start = SubstractVectors(monster.position, start)
		result = append(result, NewBillBoard(start, spritePlane, 1))
	}
	return result
}
