package main

// Player is a container for a Player
type Player struct {
	speed    float64
	radius   float64
	position Vector
	velocity Vector
}

// NewPlayer creates a new Player
func NewPlayer(position Vector) Player {
	return Player{1.5, 0.25, position, Vector{0, 0}}
}

// rect return the player position as Rect
func (p Player) rect() Rect {
	var halfSize = Vector{p.radius, p.radius}
	return Rect{SubstractVectors(p.position, halfSize),
		AddVectors(p.position, halfSize)}
}

func (p Player) isIntersecting(tileMap Tilemap) bool {
	rect := p.rect()
	minX := int(rect.min.x)
	maxX := int(rect.max.x)
	minY := int(rect.min.y)
	maxY := int(rect.max.y)

	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			if tileMap.GetTile(x, y).isWall() {
				return true
			}
		}
	}
	return false
}
