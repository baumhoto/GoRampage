package entity

import (
	"fmt"

	_asset "github.com/baumhoto/GoRampage/engine/asset"
	_core "github.com/baumhoto/GoRampage/engine/core"
	_map "github.com/baumhoto/GoRampage/engine/map"
)

type MonsterState int

const (
	MonsterStateIdle MonsterState = iota
	MonsterStateChasing
	MonsterStateScratching
	MonsterStateHurt
	MonsterStateDead
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
	health         float64
}

func NewMonster(position _core.Vector) Monster {
	return Monster{speed: 0.5, position: position, velocity: _core.Vector{0, 0},
		radius: 0.4375, state: MonsterStateIdle, animation: _asset.MonsterIdleAnimation,
		attackCoolDown: 0.4, health: 50}
}

func (m Monster) isDead() bool {
	return m.health <= 0
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

func (m *Monster) update(world *World, tm *_asset.TextureManager) {
	switch m.state {
	case MonsterStateIdle:
		if m.canSeePlayer(*world) {
			m.state = MonsterStateChasing
			m.animation = _asset.MonsterWalkAnimation
			m.animationTime = 0.0
		}
	case MonsterStateChasing:
		if !m.canSeePlayer(*world) {
			m.state = MonsterStateIdle
			m.animation = _asset.MonsterIdleAnimation
			m.animationTime = 0.0
			m.velocity = _core.Vector{}
			break
		}
		if m.canReachPlayer(*world) {
			m.state = MonsterStateScratching
			m.animation = _asset.MonsterScratchAnimation
			m.lastAttackTime = -m.attackCoolDown
			m.velocity = _core.Vector{}
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
	case MonsterStateHurt:
		if tm.Animations[m.animation].IsCompleted(m.animationTime) {
			m.state = MonsterStateIdle
			m.animation = _asset.MonsterIdleAnimation
			m.animationTime = 0.0
		}
	case MonsterStateDead:
		if tm.Animations[m.animation].IsCompleted(m.animationTime) {
			m.animationTime = 0.0
			m.animation = _asset.AnimationMonsterDead
		}
	default:
		fmt.Printf("default\n")
	}
}

func (m Monster) billboard(ray _core.Ray) _asset.Billboard {
	plane := ray.Direction.Orthogonal()
	return _asset.NewBillBoard(
		_core.SubstractVectors(m.position, _core.DivideVector(plane, 2)),
		plane,
		1.0)
}

func (m Monster) hitTest(ray _core.Ray) _core.Vector {
	hit := m.billboard(ray).HitTest(ray)
	if !m.isDead() && !hit.IsNil() {
		if _core.SubstractVectors(hit, m.position).Length() >= m.radius {
			hit = _core.NilVector()
		}
	}
	return hit
}
