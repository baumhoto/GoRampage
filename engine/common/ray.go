package common

type Ray struct {
	Origin, Direction Vector
}

func (r Ray) SlopeIntercept() (float64, float64) {
	slope := r.Direction.Y / r.Direction.X
	intercept := r.Origin.Y - slope*r.Origin.X
	return slope, intercept
}
