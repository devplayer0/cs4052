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

// MeshWireFrame when enabled, renders meshes in wireframe mode
var MeshWireFrame = false

// Vertex represents a vertex in a mesh (position, normal and UV coordinates)
type Vertex struct {
	Position mgl32.Vec3
	Normal   mgl32.Vec3

	UV mgl32.Vec2
}

// Mesh represents a mesh (indices and vertices)
type Mesh struct {
	Indices   []uint32
	Vertices  []Vertex
	Transform mgl32.Mat4

	VAO          uint32
	indexBuffer  *Buffer
	vertexBuffer *Buffer
}

// ReadOBJFile reads and parses a Wavefront .obj file
func ReadOBJFile(objFile string) (*obj.Object, error) {
	f, err := os.Open(objFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %v: %w", objFile, err)
	}
	defer f.Close()

	obj, err := obj.NewReader(f).Read()
	if err != nil {
		return nil, fmt.Errorf("failed to parse: %w", err)
	}

	return obj, nil
}

// NewOBJMesh creates a new mesh from a parsed OBJ
func NewOBJMesh(obj *obj.Object, trans mgl32.Mat4) *Mesh {
	var vertices []Vertex
	var indices []uint32
	var i uint32
	for _, f := range obj.Faces {
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
		Vertices:  vertices,
		Indices:   indices,
		Transform: trans,
	}

	gl.GenVertexArrays(1, &m.VAO)
	gl.BindVertexArray(m.VAO)

	m.indexBuffer = NewBuffer(gl.ELEMENT_ARRAY_BUFFER)
	m.vertexBuffer = NewBuffer(gl.ARRAY_BUFFER)

	return m
}

// NewOBJMeshFile loads a mesh from a .obj file
func NewOBJMeshFile(objFile string, trans mgl32.Mat4) (*Mesh, error) {
	obj, err := ReadOBJFile(objFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load obj: %w", err)
	}

	return NewOBJMesh(obj, trans), nil
}

// ReplaceVertices re-uploads mesh data into the vertex buffers
func (m *Mesh) ReplaceVertices(p *Program, vertices []Vertex) {
	gl.BindVertexArray(m.VAO)

	buf := &bytes.Buffer{}
	binary.Write(buf, nativeOrder, vertices)

	m.vertexBuffer.SetData(buf.Bytes())
}

// Upload writes the index buffer and vertex buffer to the GPU
func (m *Mesh) Upload(p *Program) *Mesh {
	gl.BindVertexArray(m.VAO)

	buf := &bytes.Buffer{}
	binary.Write(buf, nativeOrder, m.Indices)
	m.indexBuffer.SetData(buf.Bytes())

	m.vertexBuffer.LinkVertexPointer(p, "frag_pos", 3, gl.FLOAT, 32, 0)
	m.vertexBuffer.LinkVertexPointer(p, "normal", 3, gl.FLOAT, 32, 12)
	m.vertexBuffer.LinkVertexPointer(p, "uv", 2, gl.FLOAT, 32, 24)

	m.ReplaceVertices(p, m.Vertices)
	return m
}

// Draw renders the mesh with the given shader and projection
func (m *Mesh) Draw(p *Program, proj mgl32.Mat4, c *Camera, trans mgl32.Mat4) {
	p.Use()
	p.Project(proj, c, trans.Mul4(m.Transform))

	// Hardcode to white for now
	p.SetUniformVec3("in_color", mgl32.Vec3{1, 1, 1})
	// Hardcode shininess for now
	p.SetUniformFloat32("mat.spec_exponent", 64)

	gl.BindVertexArray(m.VAO)
	if MeshWireFrame {
		gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
	}
	gl.DrawElements(gl.TRIANGLES, int32(len(m.Indices)), gl.UNSIGNED_INT, gl.PtrOffset(0))
	gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)
}
