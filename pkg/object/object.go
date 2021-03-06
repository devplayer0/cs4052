package object

import (
	"fmt"
	"io/ioutil"

	"github.com/devplayer0/cs4052/pkg/pb"
	"github.com/devplayer0/cs4052/pkg/util"
	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"google.golang.org/protobuf/proto"
)

type joint struct {
	ID uint32

	InverseBind mgl32.Mat4
}

type vec3Key struct {
	Time  float32
	Value mgl32.Vec3
}
type quatKey struct {
	Time  float32
	Value mgl32.Quat
}
type nodeAnim struct {
	Pos   []vec3Key
	Rot   []quatKey
	Scale []vec3Key
}

func (a nodeAnim) findPos(t float32) (vec3Key, vec3Key) {
	for i := 0; i < len(a.Pos)-1; i++ {
		if t < a.Pos[i+1].Time {
			return a.Pos[i], a.Pos[i+1]
		}
	}

	return a.Pos[0], a.Pos[0]
}
func (a nodeAnim) findRot(t float32) (quatKey, quatKey) {
	for i := 0; i < len(a.Rot)-1; i++ {
		if t < a.Rot[i+1].Time {
			return a.Rot[i], a.Rot[i+1]
		}
	}

	return a.Rot[0], a.Rot[0]
}
func (a nodeAnim) findScale(t float32) (vec3Key, vec3Key) {
	for i := 0; i < len(a.Scale)-1; i++ {
		if t < a.Scale[i+1].Time {
			return a.Scale[i], a.Scale[i+1]
		}
	}

	return a.Scale[0], a.Scale[0]
}

type nodeDebug struct {
	vao uint32

	pathVAO    uint32
	pathBuffer *util.Buffer
}

type node struct {
	Name      string
	Transform mgl32.Mat4
	Joint     *joint

	Parent   *node
	Children []*node

	debug *nodeDebug
}

func buildNodeHierarchy(pb *pb.Object, cnid uint32, anims []*Animation, current *node) {
	cn := pb.Hierarchy[cnid]
	current.Name = cn.Name
	current.Transform = util.PBMat4(cn.Transform)

	if cn.JointID != nil {
		current.Joint = &joint{
			ID: *cn.JointID,

			InverseBind: util.PBMat4(pb.Joints[*cn.JointID].InverseBind),
		}
	}
	for i, a := range pb.Animations {
		for _, c := range a.Channels {
			if c.NodeID == cnid {
				aChan := nodeAnim{}

				aChan.Pos = make([]vec3Key, len(c.PosFrames))
				for j, p := range c.PosFrames {
					aChan.Pos[j] = vec3Key{
						Time:  p.Time,
						Value: util.PBVec3(p.Value),
					}
				}
				aChan.Rot = make([]quatKey, len(c.RotFrames))
				for j, r := range c.RotFrames {
					aChan.Rot[j] = quatKey{
						Time:  r.Time,
						Value: util.PBQuat(r.Value),
					}
				}
				aChan.Scale = make([]vec3Key, len(c.ScaleFrames))
				for j, s := range c.ScaleFrames {
					aChan.Scale[j] = vec3Key{
						Time:  s.Time,
						Value: util.PBVec3(s.Value),
					}
				}

				anims[i].channels[current] = aChan
			}
		}
	}

	for _, ccnid := range cn.Children {
		child := &node{
			Parent: current,
		}
		current.Children = append(current.Children, child)

		buildNodeHierarchy(pb, ccnid, anims, child)
	}
}

func (n *node) setupDebug(p *util.Program) {
	d := &nodeDebug{}

	gl.GenVertexArrays(1, &d.vao)
	gl.BindVertexArray(d.vao)
	cubeBuffer := util.NewBuffer(gl.ARRAY_BUFFER)
	cubeBuffer.Bind()
	cubeBuffer.SetVec3(util.CubeVertices)
	cubeBuffer.LinkVertexPointer(p, "frag_pos", 3, gl.FLOAT, 0, 0)

	gl.GenVertexArrays(1, &d.pathVAO)
	gl.BindVertexArray(d.pathVAO)
	d.pathBuffer = util.NewBuffer(gl.ARRAY_BUFFER)
	d.pathBuffer.Bind()
	d.pathBuffer.LinkVertexPointer(p, "frag_pos", 3, gl.FLOAT, 0, 0)

	n.debug = d

	for _, c := range n.Children {
		c.setupDebug(p)
	}
}

type nodeTraverseCallback func(n *node, parent, local, final mgl32.Mat4)

func (n *node) traverse(parent mgl32.Mat4, anim *Animation, t float32, cb nodeTraverseCallback) {
	local := n.Transform
	if anim != nil {
		if aChan, ok := anim.channels[n]; ok {
			pa, pb := aChan.findPos(t)
			pFactor := (t - pa.Time) / (pb.Time - pa.Time)
			pVec := pa.Value.Add(pb.Value.Sub(pa.Value).Mul(pFactor))
			pos := mgl32.Translate3D(pVec.X(), pVec.Y(), pVec.Z())

			ra, rb := aChan.findRot(t)
			rFactor := (t - ra.Time) / (rb.Time - ra.Time)
			rot := util.QuatSlerp(ra.Value, rb.Value, rFactor).Normalize().Mat4()

			sa, sb := aChan.findScale(t)
			sFactor := (t - sa.Time) / (sb.Time - sa.Time)
			sVec := util.InterpolateVec3(sa.Value, sb.Value, sFactor)
			scale := mgl32.Scale3D(sVec.X(), sVec.Y(), sVec.Z())

			local = pos.Mul4(rot).Mul4(scale)
		}
	}

	final := parent.Mul4(local)
	cb(n, parent, local, final)

	for _, c := range n.Children {
		c.traverse(final, anim, t, cb)
	}
}

// Animation represents a skeletal animation
type Animation struct {
	Name     string
	Duration float32
	TPS      float32

	channels map[*node]nodeAnim
}

type meshInstance struct {
	Mesh         *Mesh
	Transform    mgl32.Mat4
	InvTransform mgl32.Mat4
}

// Object represents a multi-mesh hierarchical animation with skeletal animation
// support
type Object struct {
	shader      *util.Program
	depthShader *util.Program

	materials  []*Material
	meshes     []*Mesh
	hierarchy  *node
	Animations []*Animation
	instances  []meshInstance

	Debug       bool
	debugShader *util.Program

	currentTransforms []mgl32.Mat4
	currentAnim       *Animation
	currentATime      float32
}

// NewObject creates a new object
func NewObject(obj *pb.Object, shader, depthShader, ds *util.Program) (*Object, error) {
	o := &Object{
		shader:      shader,
		depthShader: depthShader,

		hierarchy: &node{},

		debugShader: ds,

		currentTransforms: make([]mgl32.Mat4, MaxJoints),
	}

	for _, m := range obj.Materials {
		mat, err := LoadSOBJMaterial(m)
		if err != nil {
			return nil, fmt.Errorf("failed to load material %v: %w", m.Name, err)
		}

		o.materials = append(o.materials, mat)
	}

	for _, m := range obj.Meshes {
		cm := NewSOBJMesh(m, o.materials[m.MaterialID]).
			Upload(shader).
			LinkDepthMap(depthShader)
		o.meshes = append(o.meshes, cm)
	}

	for _, a := range obj.Animations {
		ca := &Animation{
			Name:     a.Name,
			Duration: a.Duration,
			TPS:      a.Tps,

			channels: make(map[*node]nodeAnim),
		}

		o.Animations = append(o.Animations, ca)
	}

	buildNodeHierarchy(obj, 0, o.Animations, o.hierarchy)
	if ds != nil {
		o.hierarchy.setupDebug(ds)
	}

	for _, i := range obj.Instances {
		t := util.PBMat4(i.Transform)
		o.instances = append(o.instances, meshInstance{
			Mesh:         o.meshes[i.MeshID],
			Transform:    t,
			InvTransform: t.Inv(),
		})
	}

	return o, nil
}

// NewObjectFile creates a new object from a file
func NewObjectFile(objFile string, shader, depthShader, ds *util.Program) (*Object, error) {
	data, err := ioutil.ReadFile(objFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var obj pb.Object
	if err := proto.Unmarshal(data, &obj); err != nil {
		return nil, fmt.Errorf("failed to unmarshal: %w", err)
	}

	return NewObject(&obj, shader, depthShader, ds)
}

// Update updates the state of each of the object's joint transforms
func (o *Object) Update(proj mgl32.Mat4, cam *util.Camera, trans mgl32.Mat4, anim *Animation, t float32) {
	o.currentAnim = anim
	if anim != nil {
		o.currentATime = util.Mod(t*anim.TPS, anim.Duration)
	} else {
		o.currentATime = 0
	}

	invTrans := trans.Inv()
	o.hierarchy.traverse(trans, o.currentAnim, o.currentATime, func(n *node, parent, local, final mgl32.Mat4) {
		if n.Joint != nil {
			o.currentTransforms[n.Joint.ID] = invTrans.Mul4(final).Mul4(n.Joint.InverseBind)
		}
	})
}

// DepthMapPass renders the object only for depth information
func (o *Object) DepthMapPass(trans mgl32.Mat4, depthParamsApplicator util.DepthMapParamsApplicator) {
	for _, in := range o.instances {
		ts := make([]mgl32.Mat4, len(o.currentTransforms))
		for i, t := range o.currentTransforms {
			ts[i] = in.InvTransform.Mul4(t)
		}

		o.depthShader.SetUniformMat4Slice("joints", ts)
		in.Mesh.DepthMapPass(o.depthShader, trans.Mul4(in.Transform), depthParamsApplicator)
	}
}

// Draw the object
func (o *Object) Draw(proj mgl32.Mat4, cam *util.Camera, trans mgl32.Mat4, envMap, depthMaps *util.Texture) {
	for _, in := range o.instances {
		ts := make([]mgl32.Mat4, len(o.currentTransforms))
		for i, t := range o.currentTransforms {
			ts[i] = in.InvTransform.Mul4(t)
		}

		o.shader.SetUniformMat4Slice("joints", ts)
		in.Mesh.Draw(o.shader, proj, cam, trans.Mul4(in.Transform), envMap, depthMaps)
	}

	if o.Debug && o.debugShader != nil {
		o.hierarchy.traverse(trans, o.currentAnim, o.currentATime, func(n *node, parent, local, final mgl32.Mat4) {
			o.debugShader.Use()
			o.debugShader.Project(proj, cam, final.Mul4(mgl32.Scale3D(0.05, 0.05, 0.05)))
			gl.BindVertexArray(n.debug.vao)
			gl.DrawArrays(gl.TRIANGLES, 0, int32(len(util.CubeVertices)))

			if n.Parent != nil {
				o.debugShader.Project(proj, cam, parent)

				gl.BindVertexArray(n.debug.pathVAO)
				n.debug.pathBuffer.SetVec3([]mgl32.Vec3{
					{0, 0, 0},
					util.PosFromTrans(local),
				})
				gl.DrawArrays(gl.LINES, 0, 2)
			}
		})
	}
}
