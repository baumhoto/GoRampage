package main

import (
	"fmt"
	"os"

	"github.com/baumhoto/prototype/draw"
)

// Renderer renders to the window
type Renderer struct {
	frameBuffer FrameBuffer
}

// NewRenderer creates a new instance of a Renderer
func NewRenderer(width int, height int) Renderer {
	fb := NewFrameBuffer(width, height, black)
	return Renderer{fb}
}

// draw renders the world into the window
func (r *Renderer) draw(world World, window draw.Window) {
	_, height := window.Size()
	scale := float64(height) / float64(world.worldmap.Height)

	// Draw map
	for y := 0; y < world.worldmap.Height; y++ {
		for x := 0; x < world.worldmap.Width; x++ {
			if world.worldmap.GetTile(x, y).isWall() {
				min := Vector{float64(x), float64(y)}
				min.Multiply(scale)
				max := Vector{float64(x + 1), float64(y + 1)}
				max.Multiply(scale)
				rect := Rect{min, max}
				r.frameBuffer.Fill(rect, gray)
			}
		}
	}

	// Draw player
	rect := world.player.rect()
	rect.min.Multiply(scale)
	rect.max.Multiply(scale)
	r.frameBuffer.Fill(rect, blue)
	err := window.DrawRGBA(r.frameBuffer.ToRGBA())
	if err != nil {
		fmt.Printf("LoadTexture error: %v\n", err)
		os.Exit(-1)
	}
}
