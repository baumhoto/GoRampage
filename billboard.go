package main

import "math"

type Billboard struct {
	start     Vector
	end       Vector
	direction Vector
	length    float64
	texture   Texture
}

func NewBillBoard(start Vector, direction Vector, length float64, texture Texture) Billboard {
	end := MultiplyVector(direction, length)
	end.Add(start)
	return Billboard{start, end, direction, length, texture}
}

func (b Billboard) hitTest(ray Ray) Vector {
	lhs := ray
	rhs := Ray{b.start, b.direction}

	// Ensure rays are never exactly vertical
	epsilon := 0.00001
	if math.Abs(lhs.direction.x) < epsilon {
		lhs.direction.x = epsilon
	}
	if math.Abs(rhs.direction.x) < epsilon {
		rhs.direction.x = epsilon
	}

	// Calculate slopes and intercepts
	slope1, intercept1 := lhs.slopeIntercept()
	slope2, intercept2 := rhs.slopeIntercept()

	// Check if slopes are parallel
	if slope1 == slope2 {
		return Vector{}
	}

	// Find intersection point
	x := (intercept1 - intercept2) / (slope2 - slope1)
	y := slope1*x + intercept1

	// Check intersection point is in range
	distanceAlongRay := (x - lhs.origin.x) / lhs.direction.x
	if distanceAlongRay < 0 {
		return Vector{}
	}

	distanceAlongBillboard := (x - rhs.origin.x) / rhs.direction.x
	if distanceAlongBillboard < 0 || distanceAlongBillboard > b.length {
		return Vector{}
	}

	return Vector{x, y}
}
