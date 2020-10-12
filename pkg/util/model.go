package util

import (
	"fmt"
	"os"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/sheenobu/go-obj/obj"
)

type Model struct {
	Vertices []mgl32.Vec3
	Normals  []mgl32.Vec3
}

func NewModel(objFile string) (*Model, error) {
	f, err := os.Open(objFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open obj file %v: %w", objFile, err)
	}
	defer f.Close()

	o, err := obj.NewReader(f).Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read obj: %w", err)
	}

	vertices := make([]mgl32.Vec3, len(o.Vertices))
	for i, v := range o.Vertices {
		vertices[i] = mgl32.Vec3{float32(v.X), float32(v.Y), float32(v.Z)}
	}

	normals := make([]mgl32.Vec3, len(o.Normals))
	for i, n := range o.Normals {
		normals[i] = mgl32.Vec3{float32(n.X), float32(n.Y), float32(n.Z)}
	}

	m := &Model{
		Vertices: vertices,
		Normals:  normals,
	}

	return m, nil
}
