package main

const (
	Floor = iota
	Wall
)

type Tile struct {
	Tiletype int
}

func (t Tile) isWall() bool {
	switch t.Tiletype {
	case 1:
		return true
	default:
		return false
	}
}
