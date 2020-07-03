package entity

import (
	_core "github.com/baumhoto/GoRampage/engine/core"
	_map "github.com/baumhoto/GoRampage/engine/map"
)

type Actor interface {
	intersection(tileMap _map.Tilemap) (bool, _core.Vector)
	rect() _core.Rect
}

func rect(radius float64, position _core.Vector) _core.Rect {
	var halfSize = _core.Vector{radius, radius}
	return _core.Rect{_core.SubstractVectors(position, halfSize),
		_core.AddVectors(position, halfSize)}
}

func intersection(rect _core.Rect, tileMap _map.Tilemap) (bool, _core.Vector) {
	minX := int(rect.Min.X)
	maxX := int(rect.Max.X)
	minY := int(rect.Min.Y)
	maxY := int(rect.Max.Y)

	largestIntersection := _core.Vector{}
	result := false
	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			if tileMap.GetTile(x, y).IsWall() {
				min := _core.Vector{float64(x), float64(y)}
				max := _core.Vector{float64(x + 1), float64(y + 1)}
				wallRect := _core.Rect{min, max}
				if ok, intersection := rect.Intersection(wallRect); ok && intersection.Length() > largestIntersection.Length() {
					largestIntersection = intersection
					result = true
				}
			}
		}
	}
	return result, largestIntersection
}
