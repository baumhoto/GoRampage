package entity

import (
	_asset "github.com/baumhoto/go-rampage/engine/asset"
	_common "github.com/baumhoto/go-rampage/engine/common"
	_input "github.com/baumhoto/go-rampage/engine/input"
	_map "github.com/baumhoto/go-rampage/engine/map"
	"image/color"
)

// World is a container for the world
type World struct {
	Worldmap _map.Tilemap
	Player   Player
	Monsters []Monster
	Effects  []_asset.Effect
}

// NewWorld creates a new World.
func NewWorld(worldmap _map.Tilemap) World {
	var player Player
	world := World{worldmap, player, nil, nil}
	world.reset()
	return world
}

// update updates the World
func (w *World) Update(timeStep float64, input _input.Input) {
	// update effects
	var effectsInProgress []_asset.Effect
	for _, effect := range w.Effects {
		effect.Time += timeStep
		if !effect.IsCompleted() {
			effectsInProgress = append(effectsInProgress, effect)
		}
	}

	w.Effects = effectsInProgress

	// update player
	if !w.Player.isDead() {
		w.Player.Direction = w.Player.Direction.Rotated(input.Rotation)
		w.Player.velocity = _common.MultiplyVector(w.Player.Direction, input.Speed*w.Player.speed)
		w.Player.velocity.Multiply(timeStep)
		w.Player.Position.Add(w.Player.velocity)
	} else if len(w.Effects) == 0 {
		w.reset()
		w.Effects = append(w.Effects, _asset.NewEffect(_asset.FadeIn, _common.Red, 0.5))
		return
	}

	// update monsters
	for i, _ := range w.Monsters {
		monster := w.Monsters[i]
		monster.update(w)
		monster.position.Add(_common.MultiplyVector(monster.velocity, timeStep))
		monster.animationTime += timeStep
		w.Monsters[i] = monster
	}

	// handle collisions
	for i, _ := range w.Monsters {
		// monster player
		if success, intersection := w.Player.Rect().Intersection(w.Monsters[i].rect()); success {
			intersection.Divide(2)
			w.Player.Position.Substract(intersection)
			w.Monsters[i].position.Add(intersection)
		}

		// monster monster
		for j := i + 1; j < len(w.Monsters); j++ {
			if success, intersection := w.Monsters[i].rect().Intersection(w.Monsters[j].rect()); success {
				intersection.Divide(2)
				w.Monsters[i].position.Substract(intersection)
				w.Monsters[j].position.Add(intersection)
			}
		}

		// monster world
		for {
			if success, intersection := w.Monsters[i].intersection(w.Worldmap); success {
				w.Monsters[i].position.Substract(intersection)
			} else {
				break
			}
		}
	}

	// player world
	for {
		if ok, intersection := w.Player.intersection(w.Worldmap); ok {
			w.Player.Position.Substract(intersection)
		} else {
			break
		}
	}
}

func (w World) Sprites(tm _asset.TextureManager) []_asset.Billboard {
	spritePlane := w.Player.Direction.Orthogonal()
	var result []_asset.Billboard
	for _, monster := range w.Monsters {
		start := _common.DivideVector(spritePlane, 2)
		start = _common.SubstractVectors(monster.position, start)
		result = append(result, _asset.NewBillBoard(start, spritePlane, 1,
			tm.Animations[monster.animation].Texture(monster.animationTime)))
	}
	return result
}

func (w *World) hurtPlayer(damage float64) {
	if w.Player.isDead() {
		return
	}
	w.Player.health -= damage
	w.Effects = append(w.Effects, _asset.NewEffect(_asset.FadeIn, color.RGBA{
		R: 255,
		G: 0,
		B: 0,
		A: 191,
	}, 0.2))
	if w.Player.isDead() {
		w.Effects = append(w.Effects,
			_asset.NewEffect(_asset.FizzleOut, _common.Red, 2))
	}
}

func (w *World) reset() {
	w.Monsters = w.Monsters[:0]
	w.Effects = w.Effects[:0]
	for y := 0; y < w.Worldmap.Height; y++ {
		for x := 0; x < w.Worldmap.Width; x++ {
			position := _common.Vector{float64(x) + 0.5, float64(y) + 0.5}
			thing := w.Worldmap.Things[y*w.Worldmap.Width+x]
			switch thing {
			case 0:
				break
			case 1:
				w.Player = NewPlayer(position)
				break
			case 2:
				w.Monsters = append(w.Monsters, NewMonster(position))
				break
			}
		}
	}
}
