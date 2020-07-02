package core

import "math"

// Vector is a 2D-Vector
type Vector struct {
	X, Y float64
}

// Add Vector
func (v *Vector) Add(add Vector) {
	v.X = v.X + add.X
	v.Y = v.Y + add.Y
}

// Substract Vectors
func (v *Vector) Substract(substract Vector) *Vector {
	v.X -= substract.X
	v.Y -= substract.Y
	return v
}

// Multiply Vector
func (v *Vector) Multiply(multiplier float64) *Vector {
	v.X *= multiplier
	v.Y *= multiplier
	return v
}

// Divide Vector
func (v *Vector) Divide(divisor float64) *Vector {
	v.X /= divisor
	v.Y /= divisor
	return v
}

// rotated returns a new Vector with the rotation applied
func (v Vector) Rotated(rotation Rotation) Vector {
	return Vector{v.X*rotation.m1 + v.Y*rotation.m2,
		v.X*rotation.m3 + v.Y*rotation.m4}
}

// length returns the length of the vector
func (v Vector) Length() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

// orthogonal returns the orthogonal vector
func (v Vector) Orthogonal() Vector {
	return Vector{-v.Y, v.X}
}

// AddVectors adds 2 Vectors returning a new Vector
func AddVectors(a Vector, b Vector) Vector {
	return Vector{a.X + b.X, a.Y + b.Y}
}

// SubstractVectors substracts 2 Vectors returning a new Vector
func SubstractVectors(a Vector, b Vector) Vector {
	return Vector{a.X - b.X, a.Y - b.Y}
}

// MultiplyVector multiplies a Vector with a multiplier
func MultiplyVector(a Vector, multiplier float64) Vector {
	return Vector{a.X * multiplier, a.Y * multiplier}
}

// DivideVector divides a Vector with a divisor
func DivideVector(a Vector, divisor float64) Vector {
	return Vector{a.X / divisor, a.Y / divisor}
}
