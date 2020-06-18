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
var fullScreen bool

// TODO make frametime available from underlying window
const timeStep = 1.0 / 60.0
const maximumTimeStep = 1.0 / 20.0
const worldTimeStep = 1.0 / 120.0

func main() {
	world = NewWorld(loadMap())
	renderer = Renderer{}
	draw.RunWindow("Title", 1280, 720, update)
}

func update(window draw.Window) {
	if window.WasKeyPressed(draw.KeyEscape) {
		window.Close()
	}
	if window.WasKeyPressed(draw.KeyF) {
		if fullScreen {
			window.SetFullscreen(false)
		} else {
			window.SetFullscreen(true)
		}
		fullScreen = !fullScreen
	}

	worldSteps := math.Ceil(timeStep / worldTimeStep)
	for i := 0; i < int(worldSteps); i++ {
		world.update(float64(timeStep/worldSteps), GetInput(window))
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
	inputVector := Vector{}
	velocity := float64(1)

	if window.IsKeyDown(draw.KeyDown) {
		inputVector.y = velocity
	} else if window.IsKeyDown(draw.KeyUp) {
		inputVector.y = velocity * -1
	}
	if window.IsKeyDown(draw.KeyLeft) {
		inputVector.x = velocity * -1
	} else if window.IsKeyDown(draw.KeyRight) {
		inputVector.x = velocity
	}

	rotation := inputVector.x * world.player.turningSpeed * worldTimeStep

	return Input{
		speed:    -inputVector.y,
		rotation: NewRotation(math.Sin(rotation), math.Cos(rotation)),
	}
}
