package main

type Tile int

func (t Tile) isWall() bool {
	switch t {
	case 0, 4:
		return false
	default:
		return true
	}
}
