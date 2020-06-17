package main

import "math"

// Vector is a 2D-Vector
type Vector struct {
	x, y float64
}

// Add Vector
func (v *Vector) Add(add Vector) {
	v.x = v.x + add.x
	v.y = v.y + add.y
}

// Substract Vectors
func (v *Vector) Substract(substract Vector) {
	v.x -= substract.x
	v.y -= substract.y
}

// Multiply Vector
func (v *Vector) Multiply(multiplier float64) {
	v.x *= multiplier
	v.y *= multiplier
}

// Divide Vector
func (v *Vector) Divide(divisor float64) {
	v.x /= divisor
	v.y /= divisor
}

// length returns the length of the vector
func (v Vector) length() float64 {
	return math.Sqrt(v.x*v.x + v.y*v.y)
}

// AddVectors adds 2 Vectors returning a new Vector
func AddVectors(a Vector, b Vector) Vector {
	return Vector{a.x + b.x, a.y + b.y}
}

// SubstractVectors substracts 2 Vectors returning a new Vector
func SubstractVectors(a Vector, b Vector) Vector {
	return Vector{a.x - b.x, a.y - b.y}
}

// MultiplyVector multiplies a Vector with a multiplier
func MultiplyVector(a Vector, multiplier float64) Vector {
	return Vector{a.x * multiplier, a.y * multiplier}
}

// DivideVector divides a Vector with a divisor
func DivideVector(a Vector, divisor float64) Vector {
	return Vector{a.x / divisor, a.y / divisor}
}
