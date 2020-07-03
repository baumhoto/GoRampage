package entity

import (
	"image/color"

	_asset "github.com/baumhoto/GoRampage/engine/asset"
	_consts "github.com/baumhoto/GoRampage/engine/consts"
	_core "github.com/baumhoto/GoRampage/engine/core"
	_input "github.com/baumhoto/GoRampage/engine/input"
	_map "github.com/baumhoto/GoRampage/engine/map"
)

// World is a container for the world
type World struct {
	Worldmap _map.Tilemap
	Player   Player
	Monsters []Monster
	Effects  []Effect
}

// NewWorld creates a new World.
func NewWorld() World {
	var player Player
	world := World{_map.LoadMap(), player, nil, nil}
	world.Reset()
	return world
}

// update updates the World
func (w *World) Update(timeStep float64, input _input.Input) {
	// update effects
	var effectsInProgress []Effect
	for _, effect := range w.Effects {
		effect.Time += timeStep
		if !effect.IsCompleted() {
			effectsInProgress = append(effectsInProgress, effect)
		}
	}

	w.Effects = effectsInProgress

	// update player
	if !w.Player.isDead() {
		w.Player.AnimationTime += timeStep
		w.Player.update(w, input)
		w.Player.velocity.Multiply(timeStep)
		w.Player.Position.Add(w.Player.velocity)
	} else if len(w.Effects) == 0 {
		w.Reset()
		w.Effects = append(w.Effects, NewEffect(FadeIn, _consts.RED, 0.5))
		return
	}

	// update monsters
	for i, _ := range w.Monsters {
		monster := w.Monsters[i]
		monster.animationTime += timeStep
		monster.update(w)
		monster.position.Add(_core.MultiplyVector(monster.velocity, timeStep))
		w.Monsters[i] = monster
	}

	// handle collisions
	for i, _ := range w.Monsters {
		if w.Monsters[i].isDead() {
			continue
		}
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
	ray := _core.Ray{
		Origin:    w.Player.Position,
		Direction: w.Player.Direction,
	}
	var result []_asset.Billboard
	for _, monster := range w.Monsters {
		billboard := monster.billboard(ray)
		billboard.Texture = tm.Animations[monster.animation].Texture(monster.animationTime)
		result = append(result, billboard)
	}
	return result
}

func (w *World) hurtPlayer(damage float64) {
	if w.Player.isDead() {
		return
	}
	w.Player.health -= damage
	w.Player.velocity = _core.Vector{0.0, 0.0}
	w.Effects = append(w.Effects, NewEffect(FadeIn, color.RGBA{
		R: 255,
		G: 0,
		B: 0,
		A: 191,
	}, 0.2))
	if w.Player.isDead() {
		w.Effects = append(w.Effects,
			NewEffect(FizzleOut, _consts.RED, 2))
	}
}

func (w *World) hurtMonster(index int, damage float64) {
	monster := w.Monsters[index]
	if monster.isDead() {
		return
	}
	monster.health -= damage
	monster.velocity = _core.Vector{0.0, 0.0}
	if monster.isDead() {
		monster.state = MonsterStateDead
		monster.animation = _asset.AnimationMonsterDeath
	} else {
		monster.state = MonsterStateHurt
		monster.animation = _asset.AnimationMonsterHurt
	}
	w.Monsters[index] = monster
	//fmt.Printf("%v %v\n", index, w.Monsters[index].health)
}

func (w *World) Reset() {
	w.Monsters = w.Monsters[:0]
	w.Effects = w.Effects[:0]
	for y := 0; y < w.Worldmap.Height; y++ {
		for x := 0; x < w.Worldmap.Width; x++ {
			position := _core.Vector{X: float64(x) + 0.5, Y: float64(y) + 0.5}
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

func (w World) hitTest(ray _core.Ray) int {
	wallHit := w.Worldmap.HitTest(ray)
	distance := _core.SubstractVectors(wallHit, ray.Origin).Length()
	result := -1
	for i, monster := range w.Monsters {
		monsterHit := monster.hitTest(ray)
		if monsterHit.IsNil() {
			continue
		}
		hitDistance := _core.SubstractVectors(monsterHit, ray.Origin).Length()
		if hitDistance >= distance {
			continue
		}
		result = i
		distance = hitDistance
	}

	return result
}
