package util

import (
	"math"
	"math/rand"

	"github.com/go-gl/mathgl/mgl32"
)

// Floor calculates single precision floor()
func Floor(x float32) float32 {
	return float32(math.Floor(float64(x)))
}

// Ceil calculates single precision ceil()
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

// Atan2 calculates single precision atan2()
func Atan2(y, x float32) float32 {
	return float32(math.Atan2(float64(y), float64(x)))
}

// Mod returns floating point remainder of x / y
func Mod(x, y float32) float32 {
	return float32(math.Mod(float64(x), float64(y)))
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

// InterpolateVec3 calculates the linear interpolation between two Vec3's
func InterpolateVec3(a, b mgl32.Vec3, t float32) mgl32.Vec3 {
	return mgl32.Vec3{
		Interpolate(a.X(), b.X(), t),
		Interpolate(a.Y(), b.Y(), t),
		Interpolate(a.Z(), b.Z(), t),
	}
}

// InterpolateMat4 calculates the interpolation between two 4x4 matrices
func InterpolateMat4(a, b mgl32.Mat4, t float32) mgl32.Mat4 {
	interp := mgl32.Mat4{}
	for i := 0; i < len(a); i++ {
		interp[i] = Interpolate(a[i], b[i], t)
	}

	return interp
}

// QuatSlerp performs Spherical Linear Interpolation between q1 and q2
func QuatSlerp(q1, q2 mgl32.Quat, amount float32) mgl32.Quat {
	q1, q2 = q1.Normalize(), q2.Normalize()
	dot := q1.Dot(q2)

	// Make sure we're going the right direction, this is missing from mgl32!
	if dot < 0 {
		dot = -dot
		q2 = q2.Scale(-1)
	}

	// If the inputs are too close for comfort, linearly interpolate and normalize the result.
	if dot > 0.9995 {
		return mgl32.QuatNlerp(q1, q2, amount)
	}

	// This is here for precision errors, I'm perfectly aware that *technically* the dot is bound [-1,1], but since Acos will freak out if it's not (even if it's just a liiiiitle bit over due to normal error) we need to clamp it
	dot = mgl32.Clamp(dot, -1, 1)

	theta := float32(math.Acos(float64(dot))) * amount
	c, s := float32(math.Cos(float64(theta))), float32(math.Sin(float64(theta)))
	rel := q2.Sub(q1.Scale(dot)).Normalize()

	return q1.Scale(c).Add(rel.Scale(s))
}

// Bounds represents a 3D box
type Bounds struct {
	Min mgl32.Vec3
	Max mgl32.Vec3
}

// RandVec3 returns a Vec3 with random components
func RandVec3() mgl32.Vec3 {
	return mgl32.Vec3{rand.Float32(), rand.Float32(), rand.Float32()}
}
