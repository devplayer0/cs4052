package util

import "github.com/go-gl/mathgl/mgl32"

var (
	cameraUp = mgl32.Vec3{0, 1, 0}
)

// Camera represents a 3D camera
type Camera struct {
	Position mgl32.Vec3

	rotation mgl32.Vec2
	lockY    bool

	direction mgl32.Vec3
	transform mgl32.Mat4
}

// NewCamera creates a new camera
func NewCamera(p mgl32.Vec3, r mgl32.Vec2, lockY bool) *Camera {
	c := &Camera{
		Position: p,

		lockY: lockY,
	}
	c.SetRotation(r)

	return c
}

func (c *Camera) update() {
	c.transform = mgl32.LookAtV(c.Position, c.Position.Add(c.direction), cameraUp)
}

// Rotation gets the camera's current rotation
func (c *Camera) Rotation() mgl32.Vec2 {
	return c.rotation
}

// SetRotation sets the camera's rotation
func (c *Camera) SetRotation(r mgl32.Vec2) {
	if r.Y() > 90 {
		r = mgl32.Vec2{r.X(), 90}
	} else if r.Y() < -90 {
		r = mgl32.Vec2{r.X(), -90}
	}
	c.rotation = r

	c.direction = mgl32.Vec3{
		Cos(mgl32.DegToRad(r.X())) * Cos(mgl32.DegToRad(r.Y())),
		Sin(mgl32.DegToRad(r.Y())),
		Sin(mgl32.DegToRad(r.X())) * Cos(mgl32.DegToRad(r.Y())),
	}.Normalize()
	c.update()
}

// MoveY moves the camera up and down
func (c *Camera) MoveY(d float32) {
	c.Position = c.Position.Add(mgl32.Vec3{0, d, 0})
	c.update()
}

// MoveX moves the camera in the left-right direction (based on rotation)
func (c *Camera) MoveX(d float32) {
	c.Position = c.Position.Add(c.direction.Cross(cameraUp).Normalize().Mul(d))
	c.update()
}

// MoveZ moves the camera in the forward-back direction (based on rotation)
func (c *Camera) MoveZ(d float32) {
	if c.lockY {
		c.Position = c.Position.Add(mgl32.Vec3{
			d * c.direction.X(),
			0,
			d * c.direction.Z(),
		})
	} else {
		c.Position = c.Position.Add(mgl32.Vec3{
			d * c.direction.X(),
			d * c.direction.Y(),
			d * c.direction.Z(),
		})
	}

	c.update()
}

// Transform returns the matrix for the camera
func (c *Camera) Transform() mgl32.Mat4 {
	return c.transform
}
