package asset

import (
	"github.com/baumhoto/go-rampage/engine/common"
	"image/color"
	"math"
)

type EffectType int

const (
	FadeIn EffectType = iota
	FadeOut
	FizzleOut
)

type Effect struct {
	EffectType EffectType
	Color      color.Color
	duration   float64
	Time       float64
}

func NewEffect(effectType EffectType, color color.Color, duration float64) Effect {
	return Effect{
		effectType,
		color,
		duration,
		0,
	}
}

func (e Effect) IsCompleted() bool {
	return e.Time >= e.duration
}

func (e Effect) Progress() float64 {
	t := math.Min(1.0, e.Time/e.duration)
	switch e.EffectType {
	case FadeIn:
		return common.EaseIn(t)
	case FadeOut:
		return common.EaseOut(t)
	case FizzleOut:
		return common.EaseInEaseOut(t)
	default:
		return t
	}
}
