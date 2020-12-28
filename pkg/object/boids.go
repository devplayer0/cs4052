package object

import (
	"math/rand"

	"github.com/devplayer0/cs4052/pkg/util"
	"github.com/go-gl/mathgl/mgl32"
)

type Boid struct {
	Velocity mgl32.Vec3
	Position mgl32.Vec3
}

func (b *Boid) Cohesion(bs *Boids) {
	var perceivedCOM mgl32.Vec3
	for _, ob := range bs.Instances {
		if ob == b {
			continue
		}

		perceivedCOM = perceivedCOM.Add(ob.Position)
	}
	perceivedCOM = perceivedCOM.Mul(1 / float32(len(bs.Instances)-1))
	comVec := perceivedCOM.Sub(b.Position)
	comVec = comVec.Mul(1 / comVec.Len()).Mul(bs.MaxSpeed)

	//cohesion := perceivedCOM.Sub(b.Position).Mul(1 / bs.Cohesion)
	cohesion := comVec.Sub(b.Position)
	if cohesion.Len() > bs.Cohesion {
		cohesion = cohesion.Mul(1 / cohesion.Len()).Mul(bs.Cohesion)
	}
	b.Velocity = b.Velocity.Add(cohesion)
}

//func (b *Boid) Separation(bs *Boids) {
//	var sep mgl32.Vec3
//	for _, ob := range bs.Instances {
//		if ob == b {
//			continue
//		}
//
//		diff := ob.Position.Sub(b.Position)
//		if diff.Len() < bs.Separation {
//			sep = sep.Sub(diff)
//		}
//	}
//
//	b.Velocity = b.Velocity.Add(sep)
//}
func (b *Boid) Separation(bs *Boids) {
	var avgVec mgl32.Vec3
	for _, ob := range bs.Instances {
		if ob == b {
			continue
		}

		dist := ob.Position.Sub(b.Position).Len()
		diff := b.Position.Sub(ob.Position)
		diff = diff.Mul(1 / dist)
		avgVec = avgVec.Add(diff)
	}

	avgVec = avgVec.Mul(1 / (float32(len(bs.Instances) - 1)))
	separation := avgVec.Sub(b.Velocity)
	if separation.Len() > bs.Separation {
		separation = separation.Mul(1 / separation.Len()).Mul(bs.Separation)
	}
	b.Velocity = b.Velocity.Add(separation)
}

func (b *Boid) Alignment(bs *Boids) {
	var perceivedV mgl32.Vec3
	for _, ob := range bs.Instances {
		if ob == b {
			continue
		}

		perceivedV = perceivedV.Add(ob.Velocity)
	}
	perceivedV = perceivedV.Mul(1 / float32(len(bs.Instances)-1))
	perceivedV = perceivedV.Mul(1 / perceivedV.Len()).Mul(bs.Alignment)

	alignment := perceivedV.Sub(b.Velocity)
	b.Velocity = b.Velocity.Add(alignment)
}

func (b *Boid) Edges(bs *Boids) {
	if b.Position.X() > bs.Bounds.Max.X() {
		b.Velocity[0] = -bs.MaxSpeed
	} else if b.Position.X() < bs.Bounds.Min.X() {
		b.Velocity[0] = bs.MaxSpeed
	}
	if b.Position.Y() > bs.Bounds.Max.Y() {
		b.Velocity[1] = -bs.MaxSpeed
	} else if b.Position.Y() < bs.Bounds.Min.Y() {
		b.Velocity[1] = bs.MaxSpeed
	}
	if b.Position.Z() > bs.Bounds.Max.Z() {
		b.Velocity[2] = -bs.MaxSpeed
	} else if b.Position.Z() < bs.Bounds.Min.Z() {
		b.Velocity[2] = bs.MaxSpeed
	}
}

func (b *Boid) LimitSpeed(bs *Boids) {
	if b.Velocity.Len() > bs.MaxSpeed {
		b.Velocity = b.Velocity.Normalize().Mul(bs.MaxSpeed)
	}
}

type Boids struct {
	Bounds     util.Bounds
	MaxSpeed   float32
	Cohesion   float32
	Alignment  float32
	Separation float32

	Instances []*Boid
}

func NewBoids(bounds util.Bounds, maxSpeed float32) *Boids {
	b := &Boids{
		Bounds:     bounds,
		MaxSpeed:   maxSpeed,
		Cohesion:   0.00001,
		Alignment:  0.0000001,
		Separation: 0.0001,
	}

	return b
}

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

func (bs *Boids) Update() {
	for _, b := range bs.Instances {
		b.Cohesion(bs)
		b.Separation(bs)
		//b.Alignment(bs)
		b.Edges(bs)

		b.LimitSpeed(bs)
		b.Position = b.Position.Add(b.Velocity)
	}
}
