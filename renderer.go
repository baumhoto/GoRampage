package main

import (
	"fmt"

	"github.com/gonutz/prototype/draw"
)

// Renderer renders to the window
type Renderer struct {
	frameBuffer FrameBuffer
}

// NewRenderer creates a new instance of a Renderer
func NewRenderer(width int, height int) Renderer {
	fb := NewFrameBuffer(width, height, white)
	return Renderer{fb}
}

// draw renders the world into the window
func (r *Renderer) draw(world World, window draw.Window) {
	_, height := window.Size()
	scale := float64(height) / world.size.y

	rect := world.player.rect()
	rect.min.Multiply(scale)
	rect.max.Multiply(scale)
	r.frameBuffer.Fill(rect, blue)

	err := window.DrawImageReader(r.frameBuffer.ToImageReader())
	if err != nil {
		fmt.Printf("LoadTexture error")
	}
}
