package util

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"
)

// Floor calculates single precision floor()
func Floor(x float32) float32 {
	return float32(math.Floor(float64(x)))
}

// Floor calculates single precision ceil()
func Ceil(x float32) float32 {
	return float32(math.Ceil(float64(x)))
}

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

// PosFromTrans retrieves a position vector from a translation matrix
func PosFromTrans(trans mgl32.Mat4) mgl32.Vec3 {
	return trans.Col(3).Vec3()
}

// Interpolate calculates the linear interpolation between two values
func Interpolate(a, b, t float32) float32 {
	return (b * t) + ((1 - t) * a)
}

// InterpolateMat4 calculates the interpolation between two 4x4 matrices
func InterpolateMat4(a, b mgl32.Mat4, t float32) mgl32.Mat4 {
	interp := mgl32.Mat4{}
	for i := 0; i < len(a); i++ {
		interp[i] = Interpolate(a[i], b[i], t)
	}

	return interp
}
