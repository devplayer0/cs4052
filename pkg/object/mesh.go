package object

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/sheenobu/go-obj/obj"

	"github.com/devplayer0/cs4052/pkg/pb"
	"github.com/devplayer0/cs4052/pkg/util"
)

// VertexSize is the native size of the Vertex struct (each Vec element is a
// 32-bit float)
const VertexSize = (3 + 3 + 2 + 3 + 3) * 4

// MaxJoints is the maximum number of global joints
const MaxJoints = 256

// MaxWeights is the maximum number of
const MaxWeights = 8

// WeightsSize is the size of a single weights vertex attribute
const WeightsSize = (MaxWeights + MaxWeights) * 4

// MeshWireFrame when enabled, renders meshes in wireframe mode
var MeshWireFrame = false

var zeroVec3 = mgl32.Vec3{}

// Vertex represents a vertex in a mesh (position, normal and UV coordinates)
type Vertex struct {
	Position  mgl32.Vec3
	Normal    mgl32.Vec3
	UV        mgl32.Vec2
	Tangent   mgl32.Vec3
	Bitangent mgl32.Vec3
}

// JointWeights represents the weights specific global joints have on a
// particular vertex
type JointWeights struct {
	JointIDs [MaxWeights]uint32
	Values   [MaxWeights]float32
}

// Material represents a material with textures and specular shininess
type Material struct {
	Diffuse   mgl32.Vec3
	Specular  mgl32.Vec3
	Emmissive mgl32.Vec3

	DiffuseTexture   *util.Texture
	SpecularTexture  *util.Texture
	NormalTexture    *util.Texture
	EmmissiveTexture *util.Texture

	Shininess float32
}

func loadSOBJTexture(pbTex *pb.Texture) (*util.Texture, error) {
	t := util.NewTexture(gl.TEXTURE_2D)
	if err := t.LoadPNG(pbTex.Data); err != nil {
		return nil, fmt.Errorf("failed to decode PNG: %w", err)
	}

	t.GenerateMipmap()
	t.SetIParameter(gl.TEXTURE_WRAP_S, gl.REPEAT)
	t.SetIParameter(gl.TEXTURE_WRAP_T, gl.REPEAT)
	t.SetIParameter(gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_LINEAR)
	t.SetIParameter(gl.TEXTURE_MAG_FILTER, gl.LINEAR)

	return t, nil
}

// LoadSOBJMaterial loads a material from a parsed SOBJ
func LoadSOBJMaterial(m *pb.Material) (*Material, error) {
	mat := &Material{
		Shininess: m.Shininess,
	}
	if mat.Shininess < 16 {
		mat.Shininess = 16
	}

	var err error
	if m.Diffuse != nil {
		mat.DiffuseTexture, err = loadSOBJTexture(m.Diffuse)
		if err != nil {
			return nil, fmt.Errorf("failed to load diffuse texture: %w", err)
		}
	}
	if m.Specular != nil {
		mat.SpecularTexture, err = loadSOBJTexture(m.Specular)
		if err != nil {
			return nil, fmt.Errorf("failed to load specular texture: %w", err)
		}
	}
	if m.Normal != nil {
		mat.NormalTexture, err = loadSOBJTexture(m.Normal)
		if err != nil {
			return nil, fmt.Errorf("failed to load normal map texture: %w", err)
		}
	}
	if m.Emissive != nil {
		mat.EmmissiveTexture, err = loadSOBJTexture(m.Emissive)
		if err != nil {
			return nil, fmt.Errorf("failed to load emissive texture: %w", err)
		}
	}

	return mat, nil
}

// Mesh represents a mesh (indices, vertices and optional UV's and skinning)
type Mesh struct {
	Indices  []uint32
	Vertices []Vertex
	Weights  []JointWeights

	Material *Material

	VAO          uint32
	indexBuffer  *util.Buffer
	vertexBuffer *util.Buffer
	skinBuffer   *util.Buffer
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
func NewOBJMesh(obj *obj.Object, mat *Material) *Mesh {
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
		Vertices: vertices,
		Indices:  indices,

		Material: mat,
	}
	m.init()

	return m
}

// NewOBJMeshFile loads a mesh from a .obj file
func NewOBJMeshFile(objFile string, mat *Material) (*Mesh, error) {
	obj, err := ReadOBJFile(objFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load obj: %w", err)
	}

	return NewOBJMesh(obj, mat), nil
}

func pbVertex(i *pb.Vertex) Vertex {
	return Vertex{
		Position:  util.PBVec3(i.Position),
		Normal:    util.PBVec3(i.Normal),
		UV:        util.PBVec2(i.Uv),
		Tangent:   util.PBVec3(i.Tangent),
		Bitangent: util.PBVec3(i.Bitangent),
	}
}

type vertexWeight struct {
	jointID uint32
	weight  float32
}

// NewSOBJMesh loads a mesh from a protobuf mesh
func NewSOBJMesh(im *pb.Mesh, mat *Material) *Mesh {
	m := &Mesh{
		Vertices: make([]Vertex, len(im.Vertices)),
		Indices:  make([]uint32, len(im.Faces)*3),
		Weights:  make([]JointWeights, len(im.Vertices)),

		Material: mat,
	}

	for i, v := range im.Vertices {
		m.Vertices[i] = pbVertex(v)
	}

	weightMap := make([][]vertexWeight, len(m.Vertices))
	for j, ws := range im.Weights {
		for _, w := range ws.Weights {
			weightMap[w.Vertex] = append(weightMap[w.Vertex], vertexWeight{j, w.Weight})
		}
	}
	for i, ws := range weightMap {
		if len(ws) <= MaxWeights {
			for j, w := range ws {
				vw := &m.Weights[i]
				vw.JointIDs[j] = w.jointID
				vw.Values[j] = w.weight
			}
		} else {
			log.Printf("Warning: %v weights is too many! (max %v)", len(ws), MaxWeights)
			sort.Slice(ws, func(i, j int) bool {
				return ws[i].weight < ws[j].weight
			})

			for j := 0; j < MaxWeights; j++ {
				vw := &m.Weights[i]
				vw.JointIDs[j] = ws[j].jointID
				vw.Values[j] = ws[j].weight
			}
		}
	}

	for i, f := range im.Faces {
		m.Indices[i*3] = f.A
		m.Indices[i*3+1] = f.B
		m.Indices[i*3+2] = f.C
	}

	m.init()

	return m
}

func (m *Mesh) init() {
	gl.GenVertexArrays(1, &m.VAO)
	gl.BindVertexArray(m.VAO)

	m.indexBuffer = util.NewBuffer(gl.ELEMENT_ARRAY_BUFFER)
	m.vertexBuffer = util.NewBuffer(gl.ARRAY_BUFFER)

	if len(m.Weights) > 0 {
		m.skinBuffer = util.NewBuffer(gl.ARRAY_BUFFER)
	}
}

// ReplaceVertices re-uploads mesh data into the vertex buffers
func (m *Mesh) ReplaceVertices(vertices []Vertex) {
	gl.BindVertexArray(m.VAO)

	buf := &bytes.Buffer{}
	binary.Write(buf, util.NativeOrder, vertices)

	m.vertexBuffer.SetData(buf.Bytes())
}

// Upload writes the index buffer and vertex buffer to the GPU
func (m *Mesh) Upload(p *util.Program) *Mesh {
	gl.BindVertexArray(m.VAO)

	buf := &bytes.Buffer{}
	binary.Write(buf, util.NativeOrder, m.Indices)
	m.indexBuffer.SetData(buf.Bytes())

	m.vertexBuffer.LinkVertexPointer(p, "frag_pos", 3, gl.FLOAT, VertexSize, 0)
	m.vertexBuffer.LinkVertexPointer(p, "normal", 3, gl.FLOAT, VertexSize, 12)
	m.vertexBuffer.LinkVertexPointer(p, "uv_in", 2, gl.FLOAT, VertexSize, 24)
	m.vertexBuffer.LinkVertexPointer(p, "tangent", 3, gl.FLOAT, VertexSize, 32)
	m.vertexBuffer.LinkVertexPointer(p, "bitangent", 3, gl.FLOAT, VertexSize, 44)
	m.ReplaceVertices(m.Vertices)

	if m.skinBuffer != nil {
		m.skinBuffer.LinkVertexIPointer(p, "joint_ids_a", 4, gl.UNSIGNED_INT, WeightsSize, 0)
		m.skinBuffer.LinkVertexIPointer(p, "joint_ids_b", 4, gl.UNSIGNED_INT, WeightsSize, 16)
		m.skinBuffer.LinkVertexPointer(p, "weights_a", 4, gl.FLOAT, WeightsSize, 32)
		m.skinBuffer.LinkVertexPointer(p, "weights_b", 4, gl.FLOAT, WeightsSize, 48)

		buf = &bytes.Buffer{}
		binary.Write(buf, util.NativeOrder, m.Weights)
		m.skinBuffer.SetData(buf.Bytes())
	}

	return m
}

func v3NonZero(v mgl32.Vec3) mgl32.Vec3 {
	if v == zeroVec3 {
		return mgl32.Vec3{0.0001, 0.0001, 0.0001}
	}

	return v
}

// Draw renders the mesh with the given shader and projection
func (m *Mesh) Draw(p *util.Program, proj mgl32.Mat4, c *util.Camera, trans mgl32.Mat4) {
	p.Use()
	p.Project(proj, c, trans)

	// Hardcode to white for now
	p.SetUniformVec3("in_color", mgl32.Vec3{1, 1, 1})

	if m.Material != nil {
		if m.Material.DiffuseTexture != nil {
			m.Material.DiffuseTexture.Activate(p, "tex_diffuse", 0)
		} else {
			p.SetUniformVec3("m_diffuse_color", v3NonZero(m.Material.Diffuse))
		}

		if m.Material.SpecularTexture != nil {
			m.Material.SpecularTexture.Activate(p, "tex_specular", 1)
		} else {
			p.SetUniformVec3("m_specular_color", v3NonZero(m.Material.Specular))
		}

		if m.Material.NormalTexture != nil {
			m.Material.NormalTexture.Activate(p, "tex_normal", 2)
			p.SetUniformInt("normal_map", 1)
		} else {
			p.SetUniformInt("normal_map", 0)
		}

		if m.Material.EmmissiveTexture != nil {
			m.Material.EmmissiveTexture.Activate(p, "tex_emissive", 3)
		} else {
			p.SetUniformVec3("m_emmissive_color", v3NonZero(m.Material.Emmissive))
		}

		p.SetUniformFloat32("shininess", m.Material.Shininess)
	} else {
		p.SetUniformFloat32("shininess", 16)
	}

	gl.BindVertexArray(m.VAO)
	if MeshWireFrame {
		gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
	}
	gl.DrawElements(gl.TRIANGLES, int32(len(m.Indices)), gl.UNSIGNED_INT, gl.PtrOffset(0))
	gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)
}
