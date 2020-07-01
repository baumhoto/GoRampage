package main

// Rect is a Rectangle
type Rect struct {
	min Vector
	max Vector
}

func (r Rect) intersection(rect Rect) (bool, Vector) {
	left := Vector{r.max.x - rect.min.x, 0}
	if left.x <= 0 {
		return false, Vector{}
	}
	right := Vector{r.min.x - rect.max.x, 0}
	if right.x >= 0 {
		return false, Vector{}
	}
	up := Vector{0, r.max.y - rect.min.y}
	if up.y <= 0 {
		return false, Vector{}
	}
	down := Vector{0, r.min.y - rect.max.y}
	if down.y >= 0 {
		return false, Vector{}
	}
	vectors := [4]Vector{left, right, up, down}
	var result Vector
	for i, e := range vectors {
		if i == 0 || e.length() < result.length() {
			result = e
		}
	}
	return true, result
}
