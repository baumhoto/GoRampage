package main

import (
	"time"

	"github.com/gonutz/prototype/draw"
)

var world World
var lastRenderFinishedTime time.Time

func main() {
	world = NewWorld(Vector{8, 8})
	draw.RunWindow("Title", 640, 640, update)
}

func update(window draw.Window) {
	world.update()
	renderer := NewRenderer(640, 640)
	renderer.draw(world, window)
}
