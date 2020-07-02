package main

import "fmt"

type MonsterState int

const (
	MonsterStateIdle MonsterState = iota
	MonsterStateChasing
	MonsterStateScratching
)

type Monster struct {
	speed          float64
	position       Vector
	velocity       Vector
	radius         float64
	state          MonsterState
	animation      string
	animationTime  float64
	attackCoolDown float64
	lastAttackTime float64
}

func NewMonster(position Vector) Monster {
	return Monster{speed: 0.5, position: position, velocity: Vector{0, 0},
		radius: 0.4375, state: MonsterStateIdle, animation: MonsterIdleAnimation,
		attackCoolDown: 0.4}
}

func (m Monster) rect() Rect {
	return rect(m.radius, m.position)
}

func (m Monster) intersection(tileMap Tilemap) (bool, Vector) {
	return intersection(m.rect(), tileMap)
}

func (m Monster) canSeePlayer(world World) bool {
	direction := SubstractVectors(world.player.position, m.position)
	playerDistance := direction.length()
	ray := Ray{
		origin:    m.position,
		direction: *direction.Divide(playerDistance),
	}
	wallHit := world.worldmap.hitTest(ray)
	wallDistance := wallHit.Substract(m.position).length()
	return wallDistance > playerDistance
}

func (m Monster) canReachPlayer(world World) bool {
	reach := 0.25
	playerDistance := SubstractVectors(world.player.position, m.position).length()
	return playerDistance-m.radius-world.player.radius < reach
}

func (m *Monster) update(world *World) {
	switch m.state {
	case MonsterStateIdle:
		if m.canSeePlayer(*world) {
			m.state = MonsterStateChasing
			m.velocity = Vector{0, 0}
			m.animation = MonsterWalkAnimation
			m.animationTime = 0.0
		}
	case MonsterStateChasing:
		if !m.canSeePlayer(*world) {
			m.state = MonsterStateIdle
			m.animation = MonsterIdleAnimation
			m.animationTime = 0.0
			break
		}
		if m.canReachPlayer(*world) {
			m.state = MonsterStateScratching
			m.animation = MonsterScratchAnimation
			m.lastAttackTime = -m.attackCoolDown
		}
		direction := SubstractVectors(world.player.position, m.position)
		m.velocity = *direction.Multiply(m.speed / direction.length())
	case MonsterStateScratching:
		if !m.canReachPlayer(*world) {
			m.state = MonsterStateChasing
			m.animation = MonsterWalkAnimation
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
