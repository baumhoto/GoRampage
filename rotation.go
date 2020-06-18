package main

type Rotation struct {
	m1, m2, m3, m4 float64
}

func NewRotation(sine, cosine float64) Rotation {
	return Rotation{cosine, -sine, sine, cosine}
}
