package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"time"

	"github.com/gonutz/prototype/draw"
)

var world World
var lastRenderFinishedTime time.Time
var renderer Renderer
// TODO make frametime available from underlying window
const timeStep = 1.0 / 60.0
const maximumTimeStep = 1.0 /20.0
const worldTimeStep = 1.0 / 120.0

func main() {
	world = NewWorld(loadMap())
	renderer = Renderer{}
	draw.RunWindow("Title", 640, 640, update)
}

func update(window draw.Window) {
	if window.WasKeyPressed(draw.KeyEscape) {
		window.Close()
	}

    worldSteps := math.Round(timeStep / worldTimeStep)
	for i:=0; i < int(worldSteps); i++ {
		world.update(float64(timeStep / worldSteps), GetInput(window))
	}
	renderer.draw(world, window)
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
	velocity := float64(1)

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
