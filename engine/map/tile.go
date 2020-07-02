package _map

type Tile int

func (t Tile) IsWall() bool {
	switch t {
	case 0, 4:
		return false
	default:
		return true
	}
}
