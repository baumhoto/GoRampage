package entity

import (
	"math"

	"github.com/baumhoto/GoRampage/engine/asset"
	_core "github.com/baumhoto/GoRampage/engine/core"
	_input "github.com/baumhoto/GoRampage/engine/input"
	_map "github.com/baumhoto/GoRampage/engine/map"
)

type PlayerState int

const (
	playerStateIdle PlayerState = iota
	playerStateFiring
)

// Player is a container for a Player
type Player struct {
	speed          float64
	TurningSpeed   float64
	radius         float64
	Position       _core.Vector
	velocity       _core.Vector
	Direction      _core.Vector
	health         float64
	state          PlayerState
	Animation      string
	AnimationTime  float64
	attackCooldown float64
}

// NewPlayer creates a new Player
func NewPlayer(position _core.Vector) Player {
	return Player{speed: 2, TurningSpeed: math.Pi, radius: 0.25, Position: position,
		velocity: _core.Vector{0, 0}, Direction: _core.Vector{1, 0},
		health: 100, state: playerStateIdle, Animation: asset.PistolIdleAnimation,
		attackCooldown: 0.25}
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

func (p Player) canFire() bool {
	switch p.state {
	case playerStateIdle:
		return true
	case playerStateFiring:
		return p.AnimationTime >= p.attackCooldown
	default:
		return true
	}
}

func (p *Player) update(world *World, input _input.Input) {
	p.Direction = p.Direction.Rotated(input.Rotation)
	p.velocity = _core.MultiplyVector(p.Direction, input.Speed*p.speed)

	switch p.state {
	case playerStateIdle:
		if input.IsFiring && p.canFire() {
			p.state = playerStateFiring
			p.AnimationTime = 0.0
			p.Animation = asset.PistolFireAnimation
			ray := _core.Ray{p.Position, p.Direction}
			hitIndex := world.hitTest(ray)
			if hitIndex >= 0 {
				world.hurtMonster(hitIndex, 10)
			}
		}
	case playerStateFiring:
		if p.AnimationTime >= p.attackCooldown { // TODO if animation.isCompleted()
			p.state = playerStateIdle
			p.AnimationTime = 0.0
			p.Animation = asset.PistolIdleAnimation
		}
	}
}
