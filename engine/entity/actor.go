package entity

import (
	_common "github.com/baumhoto/go-rampage/engine/common"
	_map "github.com/baumhoto/go-rampage/engine/map"
)

type Actor interface {
	intersection(tileMap _map.Tilemap) (bool, _common.Vector)
	rect() _common.Rect
}

func rect(radius float64, position _common.Vector) _common.Rect {
	var halfSize = _common.Vector{radius, radius}
	return _common.Rect{_common.SubstractVectors(position, halfSize),
		_common.AddVectors(position, halfSize)}
}

func intersection(rect _common.Rect, tileMap _map.Tilemap) (bool, _common.Vector) {
	minX := int(rect.Min.X)
	maxX := int(rect.Max.X)
	minY := int(rect.Min.Y)
	maxY := int(rect.Max.Y)

	largestIntersection := _common.Vector{}
	result := false
	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			if tileMap.GetTile(x, y).IsWall() {
				min := _common.Vector{float64(x), float64(y)}
				max := _common.Vector{float64(x + 1), float64(y + 1)}
				wallRect := _common.Rect{min, max}
				if ok, intersection := rect.Intersection(wallRect); ok && intersection.Length() > largestIntersection.Length() {
					largestIntersection = intersection
					result = true
				}
			}
		}
	}
	return result, largestIntersection
}
