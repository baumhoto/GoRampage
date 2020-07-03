package asset

import (
	"math"

	_core "github.com/baumhoto/GoRampage/engine/core"
)

type Billboard struct {
	Start     _core.Vector
	End       _core.Vector
	direction _core.Vector
	Length    float64
	Texture   Texture
}

func NewBillBoard(start _core.Vector, direction _core.Vector, length float64) Billboard {
	end := _core.MultiplyVector(direction, length)
	end.Add(start)
	return Billboard{
		Start:     start,
		End:       end,
		direction: direction,
		Length:    length,
	}
}

func (b Billboard) HitTest(ray _core.Ray) _core.Vector {
	lhs := ray
	rhs := _core.Ray{b.Start, b.direction}

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
		return _core.NilVector()
	}

	// Find intersection point
	x := (intercept1 - intercept2) / (slope2 - slope1)
	y := slope1*x + intercept1

	// Check intersection point is in range
	distanceAlongRay := (x - lhs.Origin.X) / lhs.Direction.X
	if distanceAlongRay < 0 {
		return _core.NilVector()
	}

	distanceAlongBillboard := (x - rhs.Origin.X) / rhs.Direction.X
	if distanceAlongBillboard < 0 || distanceAlongBillboard > b.Length {
		return _core.NilVector()
	}

	return _core.Vector{x, y}
}
