package main

import "math"

type Tilemap struct {
	Tiles  []int `json:"tiles"`
	Things []int `json:"things`
	Width  int   `json:"width"`
	Height int   `json:"height"`
}

func (tm Tilemap) GetTile(x, y int) Tile {
	return Tile(tm.Tiles[y*tm.Width+x])
}

func (tm Tilemap) tile(position, direction Vector) Tile {
	var offsetX, offsetY int
	if math.Floor(position.x) == position.x {
		offsetX = -1
		if direction.x > 0 {
			offsetX = 0
		}
	}
	if math.Floor(position.y) == position.y {
		offsetY = -1
		if direction.y > 0 {
			offsetY = 0
		}
	}

	return tm.GetTile(int(position.x)+offsetX, int(position.y)+offsetY)
}

func (tm Tilemap) hitTest(ray Ray) Vector {
	position := ray.origin
	slope := ray.direction.x / ray.direction.y
	for {
		var edgeDistanceX, edgeDistanceY float64
		if ray.direction.x > 0 {
			edgeDistanceX = math.Floor(position.x) + 1 - position.x
		} else {
			edgeDistanceX = math.Ceil(position.x) - 1 - position.x
		}
		if ray.direction.y > 0 {
			edgeDistanceY = math.Floor(position.y) + 1 - position.y
		} else {
			edgeDistanceY = math.Ceil(position.y) - 1 - position.y
		}

		step1 := Vector{edgeDistanceX, edgeDistanceX / slope}
		step2 := Vector{edgeDistanceY * slope, edgeDistanceY}

		if step1.length() < step2.length() {
			position.Add(step1)
		} else {
			position.Add(step2)
		}

		if tm.tile(position, ray.direction).isWall() == true {
			break
		}
	}

	return position
}
