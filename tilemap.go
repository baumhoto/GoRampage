package main

type Tilemap struct {
	Tiles  []int `json:"tiles"`
	Things []int `json:"things`
	Width  int   `json:"width"`
	Height int   `json:"height"`
}

func (tm Tilemap) GetTile(x, y int) Tile {
	return Tile{tm.Tiles[y*tm.Width+x]}
}
