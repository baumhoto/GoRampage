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
	case 0, 4:
		return false
	default:
		return true
	}
}
