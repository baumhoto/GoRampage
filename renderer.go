package main

import (
	"github.com/hajimehoshi/ebiten"
	"math"
)

// Renderer renders to the window
type Renderer struct {
	frameBuffer FrameBuffer
}

func NewRenderer(width int, height int) Renderer {
	fb := NewFrameBuffer(width, height, black)
	return Renderer{fb}
}

// draw renders the world into the window
func (r *Renderer) draw(world World, screen *ebiten.Image) {
	width, height := screen.Size()

	focalLength := 1.0
	viewWidth := float64(width) / float64(height)
	viewPlane := world.player.direction.orthogonal()
	viewPlane.Multiply(viewWidth)
	viewCenter := MultiplyVector(world.player.direction, focalLength)
	viewCenter.Add(world.player.position)
	viewStart := DivideVector(viewPlane, 2)
	viewStart = SubstractVectors(viewCenter, viewStart)

	// Cast rays
	columns := width
	step := DivideVector(viewPlane, float64(columns))
	columnPosition := viewStart
	for x := 0; x < columns; x++ {
		rayDirection := SubstractVectors(columnPosition, world.player.position)
		viewPlaneDistance := rayDirection.length()
		ray := Ray{world.player.position, DivideVector(rayDirection, viewPlaneDistance)}
		lineEnd := world.worldmap.hitTest(ray)
		wallDistance := SubstractVectors(lineEnd, ray.origin).length()

		// Draw wall
		wallHeight := 1.0
		distanceRatio := viewPlaneDistance / focalLength
		perpendicular := wallDistance / distanceRatio
		realHeight := wallHeight * focalLength / perpendicular * float64(height)
		wallColor := gray
		if math.Floor(lineEnd.x) == lineEnd.x {
			wallColor = white
		}

		// original code would use drawline but somehow looks weird
		// using FillRect instead with width of 1 fixes the problem
		// window.DrawLine(i, int((float64(height)-realHeight)/2), i,
		//	int((float64(height)+realHeight)/2), wallColor)
		from := Vector{float64(x), (float64(height) - realHeight) / 2}
		to := Vector{float64(x), (float64(height) + realHeight) / 2}
		r.frameBuffer.DrawLine(from, to, wallColor)
		columnPosition.Add(step)
	}
	screen.DrawImage(r.frameBuffer.ToImage(), nil)
}
