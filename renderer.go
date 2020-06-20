package main

import (
	"github.com/hajimehoshi/ebiten"
	"math"
)

// Renderer renders to the window
type Renderer struct {
	frameBuffer FrameBuffer
	textures    TextureManager
}

func NewRenderer(width int, height int) Renderer {
	fb := NewFrameBuffer(width, height, black)
	return Renderer{fb, loadTextures()}
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

		wallTexture := r.textures.textures["textures/wall.png"]

		wallX := lineEnd.x - math.Floor(lineEnd.x)
		if math.Floor(lineEnd.x) == lineEnd.x {
			wallX = lineEnd.y - math.Floor(lineEnd.y)
		}

		textureX := int(wallX * float64(wallTexture.image.Bounds().Size().X))
		wallStart := Vector{float64(x), (float64(height)-realHeight)/2 - 0.001} // hack (add tiny offset)
		r.frameBuffer.drawColumn(textureX, wallTexture, wallStart, realHeight, height, x)

		columnPosition.Add(step)
	}
	screen.DrawImage(r.frameBuffer.ToImage(), nil)
}
