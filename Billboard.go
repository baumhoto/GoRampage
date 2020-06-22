package main

type Billboard struct {
	start     Vector
	end       Vector
	direction Vector
	length    float64
}

func NewBillBoard(start Vector, direction Vector, length float64) Billboard {
	end := MultiplyVector(direction, length)
	end.Add(start)
	return Billboard{start, end, direction, length}
}
