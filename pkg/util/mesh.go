package util

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/sheenobu/go-obj/obj"
)

// Vertex represents a vertex in a mesh (position, normal and UV coordinates)
type Vertex struct {
	Position mgl32.Vec3
	Normal   mgl32.Vec3

	UV mgl32.Vec2
}

// Mesh represents a mesh (indices and vertices)
type Mesh struct {
	Indices  []uint32
	Vertices []Vertex
}

// NewMesh loads a mesh from a .obj file
func NewMesh(objFile string) (*Mesh, error) {
	f, err := os.Open(objFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open obj file %v: %w", objFile, err)
	}
	defer f.Close()

	o, err := obj.NewReader(f).Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read obj: %w", err)
	}

	var vertices []Vertex
	var indices []uint32
	var i uint32
	for _, f := range o.Faces {
		for _, p := range f.Points {
			var uv mgl32.Vec2
			if p.Texture != nil {
				uv = mgl32.Vec2{float32(p.Texture.U), float32(p.Texture.V)}
			}

			indices = append(indices, i)
			vertices = append(vertices, Vertex{
				Position: mgl32.Vec3{float32(p.Vertex.X), float32(p.Vertex.Y), float32(p.Vertex.Z)},
				Normal:   mgl32.Vec3{float32(p.Normal.X), float32(p.Normal.Y), float32(p.Normal.Z)},

				UV: uv,
			})

			i++
		}
	}

	m := &Mesh{
		Vertices: vertices,
		Indices:  indices,
	}

	return m, nil
}

// UploadToProgram uploads mesh data into a new buffer attached to a shader program
func (m *Mesh) UploadToProgram(p *Program) {
	p.Use()

	buf := &bytes.Buffer{}
	binary.Write(buf, nativeOrder, m.Indices)

	indexBuffer := NewBuffer(gl.ELEMENT_ARRAY_BUFFER)
	indexBuffer.SetData(buf.Bytes())

	buf = &bytes.Buffer{}
	binary.Write(buf, nativeOrder, m.Vertices)

	vertexBuffer := NewBuffer(gl.ARRAY_BUFFER)
	vertexBuffer.SetData(buf.Bytes())

	p.LinkVertexPointer("vPosition", 3, gl.FLOAT, 32, vertexBuffer, 0)
}
