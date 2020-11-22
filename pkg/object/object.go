package object

import (
	"github.com/devplayer0/cs4052/pkg/util"
	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

var (
	jointDebugShader *util.Program
)

type jointDebug struct {
	vao uint32

	boneVAO    uint32
	boneBuffer *util.Buffer
}

// Joint represents a joint in the skeleton (hierarchical)
type Joint struct {
	Keyframes []mgl32.Mat4

	Path     string
	Parent   *Joint
	Children map[string]*Joint

	debug *jointDebug
}

// SetupHierarchy recursively wires up the path and parent
func (j *Joint) SetupHierarchy() {
	for n, c := range j.Children {
		c.Path = j.Path + "." + n
		c.Parent = j
		c.SetupHierarchy()
	}
}

// SetupDebug recursively generates VBO's for rendering the skeleton (joints and
// bones)
func (j *Joint) SetupDebug(p *util.Program) {
	d := &jointDebug{}

	gl.GenVertexArrays(1, &d.vao)
	gl.BindVertexArray(d.vao)
	cubeBuffer := util.NewBuffer(gl.ARRAY_BUFFER)
	cubeBuffer.Bind()
	cubeBuffer.SetVec3(util.CubeVertices)
	cubeBuffer.LinkVertexPointer(p, "frag_pos", 3, gl.FLOAT, 0, 0)

	gl.GenVertexArrays(1, &d.boneVAO)
	gl.BindVertexArray(d.boneVAO)
	d.boneBuffer = util.NewBuffer(gl.ARRAY_BUFFER)
	d.boneBuffer.Bind()
	d.boneBuffer.LinkVertexPointer(p, "frag_pos", 3, gl.FLOAT, 0, 0)

	j.debug = d

	for _, c := range j.Children {
		c.SetupDebug(p)
	}
}

type jointDrawCallback func(j *Joint, parentTrans, localTrans, finalTrans mgl32.Mat4)

// Draw poses the skeleton recursively (with animation)
func (j *Joint) Draw(proj mgl32.Mat4, cam *util.Camera, trans mgl32.Mat4, t float32, cb jointDrawCallback) {
	local := j.Keyframes[0]

	// normalize time
	t = t - util.Floor(t)

	approxFrame := t * float32(len(j.Keyframes)-1)
	curFrame := util.Floor(approxFrame)

	curFrameI := int(curFrame)
	if curFrameI != len(j.Keyframes)-1 {
		// how far (normalized) between two frames
		interp := approxFrame - curFrame
		local = util.InterpolateMat4(j.Keyframes[curFrameI], j.Keyframes[curFrameI+1], interp)
	}

	finalTrans := trans.Mul4(local)
	cb(j, trans, local, finalTrans)

	for _, c := range j.Children {
		c.Draw(proj, cam, finalTrans, t, cb)
	}
}

// Mesh represents a mesh with bone vertex weights
type Mesh struct {
	Mesh *util.Mesh

	VertexWeights map[string][]float32
}

// Object represents an object with many meshes and a skeleton
type Object struct {
	transform    mgl32.Mat4
	invTransform mgl32.Mat4

	Debug       bool
	debugShader *util.Program
	Skeleton    *Joint
	Meshes      map[string]*Mesh
}

// NewObject creates a new object
func NewObject(s *Joint, t mgl32.Mat4, ds *util.Program) *Object {
	o := &Object{
		invTransform: t.Inv(),

		debugShader: ds,
		Skeleton:    s,
	}
	o.SetTransform(t)

	o.Skeleton.Path = "root"
	o.Skeleton.SetupHierarchy()
	if ds != nil {
		o.Skeleton.SetupDebug(ds)
	}

	return o
}

// GetTransform gets the object's transform
func (o *Object) GetTransform() mgl32.Mat4 {
	return o.transform
}

// SetTransform sets the object's transform
func (o *Object) SetTransform(t mgl32.Mat4) {
	o.transform = t
	o.invTransform = o.transform.Inv()
}

// Draw poses the skeleton and calculates skin matrices (for every vertex of
// every mesh) and finally renders each of the meshes (optionally rendering the
// skeleton if debugging is enabled)
func (o *Object) Draw(p *util.Program, proj mgl32.Mat4, cam *util.Camera, t float32) {
	boneTransforms := map[string]mgl32.Mat4{}

	o.Skeleton.Draw(proj, cam, o.transform, t, func(j *Joint, parentTrans, localTrans, finalTrans mgl32.Mat4) {
		// Use the inverted base transform since vertices will be transformed externally later
		boneTransforms[j.Path] = o.invTransform.Mul4(finalTrans)

		if o.Debug && o.debugShader != nil {
			o.debugShader.Use()
			o.debugShader.Project(proj, cam, finalTrans.Mul4(mgl32.Scale3D(0.05, 0.05, 0.05)))
			gl.BindVertexArray(j.debug.vao)
			gl.DrawArrays(gl.TRIANGLES, 0, int32(len(util.CubeVertices)))

			if j.Parent != nil {
				o.debugShader.Project(proj, cam, parentTrans)

				gl.BindVertexArray(j.debug.boneVAO)
				j.debug.boneBuffer.SetVec3([]mgl32.Vec3{
					{0, 0, 0},
					util.PosFromTrans(localTrans),
				})
				gl.DrawArrays(gl.LINES, 0, 2)
			}
		}
	})

	for _, m := range o.Meshes {
		transformedVertices := make([]util.Vertex, len(m.Mesh.Vertices))
		for i, v := range m.Mesh.Vertices {
			newVertex := util.Vertex{
				UV: v.UV,
			}

			var totalWeight float32
			for jointName, weights := range m.VertexWeights {
				weight := float32(1)
				if len(weights) == 1 {
					weight = weights[0]
				} else if len(weights) != 0 {
					weight = weights[i]
				}

				newVertex.Position = newVertex.Position.Add(boneTransforms[jointName].Mul4x1(v.Position.Vec4(1)).Mul(weight).Vec3())
				newVertex.Normal = newVertex.Normal.Add(boneTransforms[jointName].Mul4x1(v.Normal.Vec4(1)).Mul(weight).Vec3())

				totalWeight += weight
			}

			if totalWeight != 1 {
				normWeight := 1 / totalWeight
				newVertex.Position = newVertex.Position.Mul(normWeight)
				newVertex.Normal = newVertex.Normal.Mul(normWeight)
			}

			transformedVertices[i] = newVertex
		}

		m.Mesh.ReplaceVertices(transformedVertices)
		m.Mesh.Draw(p, proj, cam, o.transform)
	}
}
