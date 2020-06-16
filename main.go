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
	// if renderer.window == nil {
	// 	renderer.setWindow(window)
	// }
	var lastFrameTime int64 = 1
	// if !lastRenderFinishedTime.IsZero() {
	// 	now := time.Now()
	// 	elapsed := now.Sub(lastRenderFinishedTime)
	// 	lastFrameTime = elapsed.Milliseconds()
	// }
	world.update(lastFrameTime)
	lastRenderFinishedTime = time.Now()

	renderer := NewRenderer(640, 640)
	renderer.draw(world, window)
}
