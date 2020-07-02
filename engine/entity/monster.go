package entity

import (
	"fmt"
	_asset "github.com/baumhoto/go-rampage/engine/asset"
	_core "github.com/baumhoto/go-rampage/engine/core"
	_map "github.com/baumhoto/go-rampage/engine/map"
)

type MonsterState int

const (
	MonsterStateIdle MonsterState = iota
	MonsterStateChasing
	MonsterStateScratching
)

type Monster struct {
	speed          float64
	position       _core.Vector
	velocity       _core.Vector
	radius         float64
	state          MonsterState
	animation      string
	animationTime  float64
	attackCoolDown float64
	lastAttackTime float64
}

func NewMonster(position _core.Vector) Monster {
	return Monster{speed: 0.5, position: position, velocity: _core.Vector{0, 0},
		radius: 0.4375, state: MonsterStateIdle, animation: _asset.MonsterIdleAnimation,
		attackCoolDown: 0.4}
}

func (m Monster) rect() _core.Rect {
	return rect(m.radius, m.position)
}

func (m Monster) intersection(tileMap _map.Tilemap) (bool, _core.Vector) {
	return intersection(m.rect(), tileMap)
}

func (m Monster) canSeePlayer(world World) bool {
	direction := _core.SubstractVectors(world.Player.Position, m.position)
	playerDistance := direction.Length()
	ray := _core.Ray{
		Origin:    m.position,
		Direction: *direction.Divide(playerDistance),
	}
	wallHit := world.Worldmap.HitTest(ray)
	wallDistance := wallHit.Substract(m.position).Length()
	return wallDistance > playerDistance
}

func (m Monster) canReachPlayer(world World) bool {
	reach := 0.25
	playerDistance := _core.SubstractVectors(world.Player.Position, m.position).Length()
	return playerDistance-m.radius-world.Player.radius < reach
}

func (m *Monster) update(world *World) {
	switch m.state {
	case MonsterStateIdle:
		if m.canSeePlayer(*world) {
			m.state = MonsterStateChasing
			m.velocity = _core.Vector{0, 0}
			m.animation = _asset.MonsterWalkAnimation
			m.animationTime = 0.0
		}
	case MonsterStateChasing:
		if !m.canSeePlayer(*world) {
			m.state = MonsterStateIdle
			m.animation = _asset.MonsterIdleAnimation
			m.animationTime = 0.0
			break
		}
		if m.canReachPlayer(*world) {
			m.state = MonsterStateScratching
			m.animation = _asset.MonsterScratchAnimation
			m.lastAttackTime = -m.attackCoolDown
		}
		direction := _core.SubstractVectors(world.Player.Position, m.position)
		m.velocity = *direction.Multiply(m.speed / direction.Length())
	case MonsterStateScratching:
		if !m.canReachPlayer(*world) {
			m.state = MonsterStateChasing
			m.animation = _asset.MonsterWalkAnimation
			m.animationTime = 0.0
		}
		if m.animationTime-m.lastAttackTime >= m.attackCoolDown {
			m.lastAttackTime = m.animationTime
			world.hurtPlayer(10)
		}
	default:
		fmt.Printf("default\n")
	}
}
