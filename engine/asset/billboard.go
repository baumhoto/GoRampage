package asset

import (
	_common "github.com/baumhoto/go-rampage/engine/common"
	"math"
)

type Billboard struct {
	Start     _common.Vector
	End       _common.Vector
	direction _common.Vector
	Length    float64
	Texture   Texture
}

func NewBillBoard(start _common.Vector, direction _common.Vector, length float64, texture Texture) Billboard {
	end := _common.MultiplyVector(direction, length)
	end.Add(start)
	return Billboard{start, end, direction, length, texture}
}

func (b Billboard) HitTest(ray _common.Ray) _common.Vector {
	lhs := ray
	rhs := _common.Ray{b.Start, b.direction}

	// Ensure rays are never exactly vertical
	epsilon := 0.00001
	if math.Abs(lhs.Direction.X) < epsilon {
		lhs.Direction.X = epsilon
	}
	if math.Abs(rhs.Direction.X) < epsilon {
		rhs.Direction.X = epsilon
	}

	// Calculate slopes and intercepts
	slope1, intercept1 := lhs.SlopeIntercept()
	slope2, intercept2 := rhs.SlopeIntercept()

	// Check if slopes are parallel
	if slope1 == slope2 {
		return _common.Vector{}
	}

	// Find intersection point
	x := (intercept1 - intercept2) / (slope2 - slope1)
	y := slope1*x + intercept1

	// Check intersection point is in range
	distanceAlongRay := (x - lhs.Origin.X) / lhs.Direction.X
	if distanceAlongRay < 0 {
		return _common.Vector{}
	}

	distanceAlongBillboard := (x - rhs.Origin.X) / rhs.Direction.X
	if distanceAlongBillboard < 0 || distanceAlongBillboard > b.Length {
		return _common.Vector{}
	}

	return _common.Vector{x, y}
}
