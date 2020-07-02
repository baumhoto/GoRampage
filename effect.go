package main

import (
	"image/color"
	"math"
)

type EffectType int

const (
	fadeIn EffectType = iota
)

type Effect struct {
	effectType EffectType
	color      color.Color
	duration   float64
	time       float64
}

func NewEffect(effectType EffectType, color color.Color, duration float64) Effect {
	return Effect{
		effectType,
		color,
		duration,
		0,
	}
}

func (e Effect) isCompleted() bool {
	return e.time >= e.duration
}

func (e Effect) progress() float64 {
	return math.Min(1.0, e.time/e.duration)
}
