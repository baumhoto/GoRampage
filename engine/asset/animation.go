package asset

import (
	"math"
)

type Animation struct {
	frames   []Texture
	duration float64
}

const (
	MonsterIdleAnimation    = "monsterIdle"
	MonsterWalkAnimation    = "monsterWalk"
	MonsterScratchAnimation = "monsterScratch"
	AnimationMonsterHurt = "monsterHurt"
	AnimationMonsterDeath = "monsterDeath"
	AnimationMonsterDead = "monsterDead"
	PistolIdleAnimation     = "pistolIdle"
	PistolFireAnimation     = "pistolFire"
)

func (a Animation) Texture(time float64) Texture {
	if a.duration == 0 {
		return a.frames[0]
	}
	t := math.Mod(time, a.duration) / a.duration
	return a.frames[int(float64(len(a.frames))*t)]
}

func (a Animation) IsCompleted(animationTime float64) bool {
	return animationTime >= a.duration
}

