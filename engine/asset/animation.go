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
)

func (a Animation) Texture(time float64) Texture {
	if a.duration == 0 {
		return a.frames[0]
	}
	t := math.Mod(time, a.duration) / a.duration
	return a.frames[int(float64(len(a.frames))*t)]
}
