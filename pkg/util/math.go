package util

import "math"

// Sin calculates single precision sin()
func Sin(x float32) float32 {
	return float32(math.Sin(float64(x)))
}

// Cos calculates single precision cos()
func Cos(x float32) float32 {
	return float32(math.Cos(float64(x)))
}
