package main

import "image/color"

// World is a container for the world
type World struct {
	worldmap Tilemap
	player   Player
	monsters []Monster
	effects  []Effect
}

// NewWorld creates a new World.
func NewWorld(worldmap Tilemap) World {
	var player Player
	world := World{worldmap, player, nil, nil}
	world.reset()
	return world
}

// update updates the World
func (w *World) update(timeStep float64, input Input) {
	// update effects
	var effectsInProgress []Effect
	for _, effect := range w.effects {
		effect.time += timeStep
		if !effect.isCompleted() {
			effectsInProgress = append(effectsInProgress, effect)
		}
	}

	w.effects = effectsInProgress

	// update player
	if !w.player.isDead() {
		w.player.direction = w.player.direction.rotated(input.rotation)
		w.player.velocity = MultiplyVector(w.player.direction, input.speed*w.player.speed)
		w.player.velocity.Multiply(timeStep)
		w.player.position.Add(w.player.velocity)
	} else if len(w.effects) == 0 {
		w.reset()
		world.effects = append(world.effects, NewEffect(fadeIn, red, 0.5))
		return
	}

	// update monsters
	for i, _ := range w.monsters {
		monster := w.monsters[i]
		monster.update(w)
		monster.position.Add(MultiplyVector(monster.velocity, timeStep))
		monster.animationTime += timeStep
		w.monsters[i] = monster
	}

	// handle collisions
	for i, _ := range w.monsters {
		// monster player
		if success, intersection := w.player.rect().intersection(w.monsters[i].rect()); success {
			intersection.Divide(2)
			w.player.position.Substract(intersection)
			w.monsters[i].position.Add(intersection)
		}

		// monster monster
		for j := i + 1; j < len(w.monsters); j++ {
			if success, intersection := w.monsters[i].rect().intersection(w.monsters[j].rect()); success {
				intersection.Divide(2)
				w.monsters[i].position.Substract(intersection)
				w.monsters[j].position.Add(intersection)
			}
		}

		// monster world
		for {
			if success, intersection := w.monsters[i].intersection(w.worldmap); success {
				w.monsters[i].position.Substract(intersection)
			} else {
				break
			}
		}
	}

	// player world
	for {
		if ok, intersection := w.player.intersection(w.worldmap); ok {
			w.player.position.Substract(intersection)
		} else {
			break
		}
	}
}

func (w World) sprites(tm TextureManager) []Billboard {
	spritePlane := w.player.direction.orthogonal()
	var result []Billboard
	for _, monster := range w.monsters {
		start := DivideVector(spritePlane, 2)
		start = SubstractVectors(monster.position, start)
		result = append(result, NewBillBoard(start, spritePlane, 1, tm.animations[monster.animation].Texture(monster.animationTime)))
	}
	return result
}

func (w *World) hurtPlayer(damage float64) {
	if w.player.isDead() {
		return
	}
	w.player.health -= damage
	w.effects = append(w.effects, NewEffect(fadeIn, color.RGBA{
		R: 255,
		G: 0,
		B: 0,
		A: 191,
	}, 0.2))
	if w.player.isDead() {
		w.effects = append(w.effects, NewEffect(fizzleOut, red, 2))
	}
}

func (w *World) reset() {
	w.monsters = w.monsters[:0]
	w.effects = w.effects[:0]
	for y := 0; y < w.worldmap.Height; y++ {
		for x := 0; x < w.worldmap.Width; x++ {
			position := Vector{float64(x) + 0.5, float64(y) + 0.5}
			thing := w.worldmap.Things[y*w.worldmap.Width+x]
			switch thing {
			case 0:
				break
			case 1:
				w.player = NewPlayer(position)
				break
			case 2:
				w.monsters = append(w.monsters, NewMonster(position))
				break
			}
		}
	}
}
