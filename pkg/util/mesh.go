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

// VertexSize is the native size of the Vertex struct (each Vec element is a
// 32-bit float, (3 + 3 + 2)*4)
const VertexSize = 32

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

	VAO uint32
}

// NewOBJMesh loads a mesh from a .obj file
func NewOBJMesh(objFile string) (*Mesh, error) {
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

// Upload uploads mesh data into buffers (vertex and element) attached to a new VAO
func (m *Mesh) Upload(p *Program) {
	gl.GenVertexArrays(1, &m.VAO)
	gl.BindVertexArray(m.VAO)

	buf := &bytes.Buffer{}
	binary.Write(buf, nativeOrder, m.Indices)

	indexBuffer := NewBuffer(gl.ELEMENT_ARRAY_BUFFER)
	indexBuffer.SetData(buf.Bytes())

	buf = &bytes.Buffer{}
	binary.Write(buf, nativeOrder, m.Vertices)

	vertexBuffer := NewBuffer(gl.ARRAY_BUFFER)
	vertexBuffer.SetData(buf.Bytes())

	vertexBuffer.LinkVertexPointer(p, "frag_pos", 3, gl.FLOAT, 32, 0)
	vertexBuffer.LinkVertexPointer(p, "normal", 3, gl.FLOAT, 32, 12)
	vertexBuffer.LinkVertexPointer(p, "uv", 2, gl.FLOAT, 32, 24)
}

// Draw renders the mesh with the given shader and projection
func (m *Mesh) Draw(p *Program, proj mgl32.Mat4, c *Camera, trans mgl32.Mat4) {
	p.Use()
	p.Project(proj, c, trans)

	// Hardcode to white for now
	p.SetUniformVec3("in_color", mgl32.Vec3{1, 1, 1})
	// Hardcode shininess for now
	p.SetUniformFloat32("mat.spec_exponent", 64)

	gl.BindVertexArray(m.VAO)
	gl.DrawElements(gl.TRIANGLES, int32(len(m.Indices)), gl.UNSIGNED_INT, gl.PtrOffset(0))
}
