package main

type Ray struct {
	origin, direction Vector
}

func (r Ray) slopeIntercept() (float64, float64) {
	slope := r.direction.y / r.direction.x
	intercept := r.origin.y - slope*r.origin.x
	return slope, intercept
}
