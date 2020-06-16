package main

import (
	"fmt"
	"image/color"

	"github.com/gonutz/prototype/draw"
)

// Renderer renders to the window
type Renderer struct {
	frameBuffer FrameBuffer
}

func NewRenderer(width int, height int) Renderer {
	fb := NewFrameBuffer(width, height, color.White)
	return Renderer{fb}
}

// draw renders the world into the window
func (r *Renderer) draw(world World, window draw.Window) {
	_, height := window.Size()
	scale := float64(height) / world.size.y

	rect := world.player.rect()
	rect.min.Multiply(scale)
	rect.max.Multiply(scale)
	r.frameBuffer.Fill(rect, color.Black)

	err := window.DrawFrameBuffer(r.frameBuffer.ToTexture())
	// err := window.DrawImageFile("image.png", 0, 0)
	if err != nil {
		fmt.Printf("LoadTexture error")
	}
}
