package main

import (
	"github.com/gonutz/prototype/draw"
	"math"
)

// Renderer renders to the window
type Renderer struct {
}

// draw renders the world into the window
func (r *Renderer) draw(world World, window draw.Window) {
	width, height := window.Size()
	//scale := float64(height) / float64(world.worldmap.Height)

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
	for i := 0; i < columns; i++ {
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
		wallColor := draw.Gray
		if math.Floor(lineEnd.x) == lineEnd.x {
			wallColor = draw.White
		}

		// original code would use drawline but somehow looks weird
		// using FillRect instead with width of 1 fixes the problem
		// window.DrawLine(i, int((float64(height)-realHeight)/2), i,
		//	int((float64(height)+realHeight)/2), wallColor)
		window.FillRect(i, int((float64(height)-realHeight)/2), 1,
			int(realHeight), wallColor)
		columnPosition.Add(step)
	}
}
