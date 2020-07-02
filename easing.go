package main

func linear(t float64) float64 {
	return t
}

func easeIn(t float64) float64 {
	return t * t
}

func easeOut(t float64) float64 {
	return 1 - easeIn(1-t)
}

func easeInEaseOut(t float64) float64 {
	if t < 0.5 {
		return 2 * easeIn(t)
	} else {
		return 4*t - 2*easeIn(t) - 1
	}
}
