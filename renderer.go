package main

import (
	"github.com/gonutz/prototype/draw"
)

// Renderer renders to the window
type Renderer struct {
}

// draw renders the world into the window
func (r *Renderer) draw(world World, window draw.Window) {
	_, height := window.Size()
	scale := float64(height) / float64(world.worldmap.Height)

	// Draw map
	for y := 0; y < world.worldmap.Height; y++ {
		for x := 0; x < world.worldmap.Width; x++ {
			if world.worldmap.GetTile(x, y).isWall() {
				window.FillRect(x*int(scale), y*int(scale), 1*int(scale), 1*int(scale), draw.White)
			}
		}
	}

	// Draw player
	rect := world.player.rect()
	rect.min.Multiply(scale)
	rect.max.Multiply(scale)
	window.FillRect(int(rect.min.x), int(rect.min.y), int(rect.max.x-rect.min.x), int(rect.max.y-rect.min.y), draw.Blue)

	// Draw view plane
	focalLength := 1.0
	viewWidth := 1.0
	viewPlane := world.player.direction.orthogonal()
	viewPlane.Multiply(viewWidth)
	viewCenter := MultiplyVector(world.player.direction, focalLength)
	viewCenter.Add(world.player.position)
	viewStart := DivideVector(viewPlane, 2)
	viewStart = SubstractVectors(viewCenter, viewStart)
	viewEnd := AddVectors(viewStart, viewPlane)
	window.DrawLine(int(viewStart.x*scale), int(viewStart.y*scale), int(viewEnd.x*scale), int(viewEnd.y*scale), draw.Red)

	// Cast rays
	columns := 10.0
	step := DivideVector(viewPlane, columns)
	columnPosition := viewStart
	for i := 0; i < int(columns); i++ {
		rayDirection := SubstractVectors(columnPosition, world.player.position)
		viewPlaneDistance := rayDirection.length()
		ray := Ray{world.player.position, DivideVector(rayDirection, viewPlaneDistance)}
		lineEnd := world.worldmap.hitTest(ray)
		window.DrawLine(int(world.player.position.x*scale), int(world.player.position.y*scale), int(lineEnd.x*scale), int(lineEnd.y*scale), draw.Green)
		columnPosition.Add(step)
	}
}
