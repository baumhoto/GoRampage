package core

func Linear(t float64) float64 {
	return t
}

func EaseIn(t float64) float64 {
	return t * t
}

func EaseOut(t float64) float64 {
	return 1 - EaseIn(1-t)
}

func EaseInEaseOut(t float64) float64 {
	if t < 0.5 {
		return 2 * EaseIn(t)
	} else {
		return 4*t - 2*EaseIn(t) - 1
	}
}
