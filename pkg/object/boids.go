package object

import (
	"math/rand"

	"github.com/devplayer0/cs4052/pkg/util"
	"github.com/go-gl/mathgl/mgl32"
)

// Boid represents a single boid (position, velocity and acceleration)
type Boid struct {
	Position mgl32.Vec3

	Velocity     mgl32.Vec3
	Acceleration mgl32.Vec3
}

// Distance finds the distance between two boids
func (b *Boid) Distance(ob *Boid) float32 {
	return ob.Position.Sub(b.Position).Len()
}

// Cohesion applies the cohesion rule, where a boid tries to steer towards the
// centre of mass for the local flock (other boids within a certain distance)
// This is basically the average of position
func (b *Boid) Cohesion(bs *Boids) {
	n := 0
	var perceivedCOM mgl32.Vec3
	for _, ob := range bs.Instances {
		if ob == b || b.Distance(ob) > bs.Perception {
			continue
		}

		perceivedCOM = perceivedCOM.Add(ob.Position)
		n++
	}
	if n == 0 {
		return
	}

	perceivedCOM = perceivedCOM.Mul(1 / float32(n))

	cohesion := perceivedCOM.Sub(b.Position).Mul(bs.CohesionFactor)
	b.Acceleration = b.Acceleration.Add(cohesion)
}

// Separation applies the separation rule, where a boid avoids other boids
// If a boid is too close to another, it will accelerate itself in the opposite
// direction with the repulsion force
func (b *Boid) Separation(bs *Boids) {
	n := 0
	var avg mgl32.Vec3
	for _, ob := range bs.Instances {
		diff := ob.Position.Sub(b.Position)
		distance := diff.Len()
		if ob == b || distance > bs.SeparationDistance {
			continue
		}

		avg = avg.Sub(diff.Mul(1 / distance))
		n++
	}
	if n == 0 {
		return
	}

	sep := avg.Mul(1 / float32(n)).Normalize().Mul(bs.RepelForce)

	b.Acceleration = b.Acceleration.Add(sep)
}

// Alignment applies the alignment rule, where a boid tries to steer its
// velocity to match others within a certain distance (the "local flock")
// This is basically the average of velocity
func (b *Boid) Alignment(bs *Boids) {
	n := 0
	var perceivedV mgl32.Vec3
	for _, ob := range bs.Instances {
		if ob == b || b.Distance(ob) > bs.Perception {
			continue
		}

		perceivedV = perceivedV.Add(ob.Velocity)
		n++
	}
	if n == 0 {
		return
	}

	perceivedV = perceivedV.Mul(1 / float32(n))

	alignment := perceivedV.Sub(b.Velocity)
	alignment = alignment.Mul(bs.AlignmentFactor)
	b.Acceleration = b.Acceleration.Add(alignment)
}

// Edges makes sure the boids can't leave their bounds by applying an opposing
// force if they are at the edge
func (b *Boid) Edges(bs *Boids) {
	if b.Position.X() > bs.Bounds.Max.X() {
		b.Velocity[0] = -bs.RepelForce
	} else if b.Position.X() < bs.Bounds.Min.X() {
		b.Velocity[0] = bs.RepelForce
	}
	if b.Position.Y() > bs.Bounds.Max.Y() {
		b.Velocity[1] = -bs.MaxSpeed
	} else if b.Position.Y() < bs.Bounds.Min.Y() {
		b.Velocity[1] = bs.RepelForce
	}
	if b.Position.Z() > bs.Bounds.Max.Z() {
		b.Velocity[2] = -bs.RepelForce
	} else if b.Position.Z() < bs.Bounds.Min.Z() {
		b.Velocity[2] = bs.RepelForce
	}
}

// LimitSpeed ensures a boid's velocity never exceeds a maximum value
func (b *Boid) LimitSpeed(bs *Boids) {
	if b.Velocity.Len() > bs.MaxSpeed {
		b.Velocity = b.Velocity.Normalize().Mul(bs.MaxSpeed)
	}
}

// Boids manages a set of boids
type Boids struct {
	Bounds             util.Bounds
	MaxSpeed           float32
	RepelForce         float32
	Perception         float32
	CohesionFactor     float32
	AlignmentFactor    float32
	SeparationDistance float32

	Instances []*Boid
}

// NewBoids creates a new boids manager
func NewBoids(bounds util.Bounds, maxSpeed float32) *Boids {
	b := &Boids{
		Bounds:             bounds,
		MaxSpeed:           maxSpeed,
		Perception:         6,
		RepelForce:         0.0003,
		CohesionFactor:     0.000004,
		AlignmentFactor:    0.0001,
		SeparationDistance: 2,
	}

	return b
}

// MakeBoid creates a new boid with a random position (within the bounds) and
// a random velocity
func (bs *Boids) MakeBoid() *Boid {
	hi := bs.Bounds.Max.Sub(bs.Bounds.Min)
	hi[0] *= rand.Float32()
	hi[1] *= rand.Float32()
	hi[2] *= rand.Float32()

	b := &Boid{
		Velocity: util.RandVec3().Mul(bs.MaxSpeed),
		Position: hi.Add(bs.Bounds.Min),
	}

	return b
}

// Update applies each of the rules to each boid and updates the current
// velocity / position
func (bs *Boids) Update() {
	for _, b := range bs.Instances {
		b.Cohesion(bs)
		b.Separation(bs)
		b.Alignment(bs)
		b.Edges(bs)

		b.Velocity = b.Velocity.Add(b.Acceleration)
		b.LimitSpeed(bs)
		b.Acceleration = mgl32.Vec3{}

		b.Position = b.Position.Add(b.Velocity)
	}
}
