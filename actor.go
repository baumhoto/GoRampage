package main

type Actor interface {
	intersection(tileMap Tilemap) (bool, Vector)
	rect() Rect
}

func rect(radius float64, position Vector) Rect {
	var halfSize = Vector{radius, radius}
	return Rect{SubstractVectors(position, halfSize),
		AddVectors(position, halfSize)}
}

func intersection(rect Rect, tileMap Tilemap) (bool, Vector) {
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
