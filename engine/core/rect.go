package core

// Rect is a Rectangle
type Rect struct {
	Min Vector
	Max Vector
}

func (r Rect) Intersection(rect Rect) (bool, Vector) {
	left := Vector{r.Max.X - rect.Min.X, 0}
	if left.X <= 0 {
		return false, Vector{}
	}
	right := Vector{r.Min.X - rect.Max.X, 0}
	if right.X >= 0 {
		return false, Vector{}
	}
	up := Vector{0, r.Max.Y - rect.Min.Y}
	if up.Y <= 0 {
		return false, Vector{}
	}
	down := Vector{0, r.Min.Y - rect.Max.Y}
	if down.Y >= 0 {
		return false, Vector{}
	}
	vectors := [4]Vector{left, right, up, down}
	var result Vector
	for i, e := range vectors {
		if i == 0 || e.Length() < result.Length() {
			result = e
		}
	}
	return true, result
}
