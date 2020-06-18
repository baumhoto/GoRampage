package main

// Player is a container for a Player
type Player struct {
	speed     float64
	radius    float64
	position  Vector
	velocity  Vector
	direction Vector
}

// NewPlayer creates a new Player
func NewPlayer(position Vector) Player {
	return Player{2, 0.25, position, Vector{0, 0}, Vector{1, 0}}
}

// rect return the player position as Rect
func (p Player) rect() Rect {
	var halfSize = Vector{p.radius, p.radius}
	return Rect{SubstractVectors(p.position, halfSize),
		AddVectors(p.position, halfSize)}
}

func (p Player) intersection(tileMap Tilemap) (bool, Vector) {
	rect := p.rect()
	minX := int(rect.min.x)
	maxX := int(rect.max.x)
	minY := int(rect.min.y)
	maxY := int(rect.max.y)

	largestIntersection := Vector{}
	result := false
	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			if tileMap.GetTile(x, y).isWall() {
				min := Vector{float64(x), float64(y)}
				max := Vector{float64(x + 1), float64(y + 1)}
				wallRect := Rect{min, max}
				if ok, intersection := rect.intersection(wallRect); ok && intersection.length() > largestIntersection.length() {
					largestIntersection = intersection
					result = true
				}
			}
		}
	}
	return result, largestIntersection
}
