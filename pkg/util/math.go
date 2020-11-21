package util

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"
)

// Sin calculates single precision sin()
func Sin(x float32) float32 {
	return float32(math.Sin(float64(x)))
}

// Cos calculates single precision cos()
func Cos(x float32) float32 {
	return float32(math.Cos(float64(x)))
}

// TransFromPos converts a position vector to a translation matrix
func TransFromPos(pos mgl32.Vec3) mgl32.Mat4 {
	return mgl32.Translate3D(pos.X(), pos.Y(), pos.Z())
}
