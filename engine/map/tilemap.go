package _map

import (
	"encoding/json"
	"fmt"
	_common "github.com/baumhoto/go-rampage/engine/common"
	"io/ioutil"
	"math"
)

type Tilemap struct {
	Tiles  []int `json:"tiles"`
	Things []int `json:"things`
	Width  int   `json:"width"`
	Height int   `json:"height"`
}

func (tm Tilemap) GetTile(x, y int) Tile {
	return Tile(tm.Tiles[y*tm.Width+x])
}

func (tm Tilemap) Tile(position, direction _common.Vector) Tile {
	var offsetX, offsetY int
	if math.Floor(position.X) == position.X {
		offsetX = -1
		if direction.X > 0 {
			offsetX = 0
		}
	}
	if math.Floor(position.Y) == position.Y {
		offsetY = -1
		if direction.Y > 0 {
			offsetY = 0
		}
	}

	return tm.GetTile(int(position.X)+offsetX, int(position.Y)+offsetY)
}

func (tm Tilemap) HitTest(ray _common.Ray) _common.Vector {
	position := ray.Origin
	slope := ray.Direction.X / ray.Direction.Y
	for {
		var edgeDistanceX, edgeDistanceY float64
		if ray.Direction.X > 0 {
			edgeDistanceX = math.Floor(position.X) + 1 - position.X
		} else {
			edgeDistanceX = math.Ceil(position.X) - 1 - position.X
		}
		if ray.Direction.Y > 0 {
			edgeDistanceY = math.Floor(position.Y) + 1 - position.Y
		} else {
			edgeDistanceY = math.Ceil(position.Y) - 1 - position.Y
		}

		step1 := _common.Vector{X: edgeDistanceX, Y: edgeDistanceX / slope}
		step2 := _common.Vector{X: edgeDistanceY * slope, Y: edgeDistanceY}

		if step1.Length() < step2.Length() {
			position.Add(step1)
		} else {
			position.Add(step2)
		}

		if tm.Tile(position, ray.Direction).IsWall() == true {
			break
		}
	}

	return position
}

func LoadMap() Tilemap {
	file, _ := ioutil.ReadFile("map.json")

	data := Tilemap{}

	err := json.Unmarshal([]byte(file), &data)
	if err != nil {
		fmt.Printf("%v\n", err)
		return Tilemap{}
	}

	return data
}
