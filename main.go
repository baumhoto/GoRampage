package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/baumhoto/prototype/draw"
)

var world World
var lastRenderFinishedTime time.Time
var renderer Renderer

func main() {
	world = NewWorld(loadMap())
	renderer = NewRenderer(640, 640)
	draw.RunWindow("Title", 640, 640, update)
}

func update(window draw.Window) {
	if window.WasKeyPressed(draw.KeyEscape) {
		window.Close()
	}

	world.update(GetInput(window))
	renderer.draw(world, window)
	renderer.frameBuffer.resetFrameBuffer()
}

func loadMap() Tilemap {
	file, _ := ioutil.ReadFile("map.json")

	data := Tilemap{}

	err := json.Unmarshal([]byte(file), &data)
	if err != nil {
		fmt.Printf("%v\n", err)
		return Tilemap{}
	}

	//fmt.Printf("%v\n", data.Tiles[12])

	return data
}

func GetInput(window draw.Window) Input {
	input := Input{}
	velocity := float64(0.01)

	if window.IsKeyDown(draw.KeyDown) {
		input.velocity.y = velocity
	} else if window.IsKeyDown(draw.KeyUp) {
		input.velocity.y = velocity * -1
	}

	if window.IsKeyDown(draw.KeyLeft) {
		input.velocity.x = velocity * -1
	} else if window.IsKeyDown(draw.KeyRight) {
		input.velocity.x = velocity
	}

	return input
}
