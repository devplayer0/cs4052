package util

import (
	"fmt"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

// ShadowResolution is the size of a single depth map's face (square)
const ShadowResolution = 2048

const lightingFragShaderFile = "assets/shaders/lighting.fs"
const lightingDepthFragShaderFile = "assets/shaders/shadows_depth.fs"
const lightingDepthGeoShaderFile = "assets/shaders/shadows_depth.gs"

const nearPlane = float32(1.0)
const farPlane = float32(45.0)

// AttenuationParams represents the lighting attenuation coefficients
type AttenuationParams struct {
	Constant  float32
	Linear    float32
	Quadratic float32
}

// DirectionalLight represents a directional light source
type DirectionalLight struct {
	Direction mgl32.Vec3

	Ambient  mgl32.Vec3
	Diffuse  mgl32.Vec3
	Specular mgl32.Vec3
}

// Lamp represents a point light source
type Lamp struct {
	Attenuation AttenuationParams

	Position mgl32.Vec3

	Ambient  mgl32.Vec3
	Diffuse  mgl32.Vec3
	Specular mgl32.Vec3
}

// Spotlight represents a directional (cone) point light source
type Spotlight struct {
	Attenuation AttenuationParams

	Position  mgl32.Vec3
	Direction mgl32.Vec3

	Cutoff      float32
	OuterCutoff float32

	Ambient  mgl32.Vec3
	Diffuse  mgl32.Vec3
	Specular mgl32.Vec3
}

type shaderTemplateData struct {
	Dirs       []*DirectionalLight
	Lamps      []*Lamp
	Spotlights []*Spotlight
}

// Lighting represents a shader to colour an object with lighting
type Lighting struct {
	viewPos    mgl32.Vec3
	dirs       []*DirectionalLight
	lamps      []*Lamp
	spotlights []*Spotlight

	cubeVAO    uint32
	cubeShader *Program

	ShadowsEnabled       bool
	depthUpdateLamps     []bool
	lampShadowTransforms []mgl32.Mat4
	depthMapsFBO         *Framebuffer
	DepthMaps            *Texture

	fragSource      string
	depthFragSource string
	depthGeoSource  string
}

// NewLighting creates a new lighting shader from a given vertex shader and set
// of lamps
func NewLighting(dirs []*DirectionalLight, lamps []*Lamp, spotlights []*Spotlight) (*Lighting, error) {
	shaderTplParams := shaderTemplateData{dirs, lamps, spotlights}
	fsSource, err := TemplateFile(lightingFragShaderFile, shaderTplParams)
	if err != nil {
		return nil, fmt.Errorf("failed to generate fragment shader source: %w", err)
	}

	dfsSource, err := TemplateFile(lightingDepthFragShaderFile, shaderTplParams)
	if err != nil {
		return nil, fmt.Errorf("failed to generate depth map fragment shader source: %w", err)
	}
	dgsSource, err := TemplateFile(lightingDepthGeoShaderFile, shaderTplParams)
	if err != nil {
		return nil, fmt.Errorf("failed to generate depth map geometry shader source: %w", err)
	}

	l := &Lighting{
		dirs:       dirs,
		lamps:      lamps,
		spotlights: spotlights,

		depthUpdateLamps:     make([]bool, len(lamps)),
		lampShadowTransforms: make([]mgl32.Mat4, len(lamps)*6),

		fragSource:      fsSource,
		depthFragSource: dfsSource,
		depthGeoSource:  dgsSource,
	}
	if l.initDebugCubes(); err != nil {
		return nil, fmt.Errorf("failed to initialize lamp cubes: %w", err)
	}
	if l.initDepthMaps(); err != nil {
		return nil, fmt.Errorf("failed to initialize depth map: %w", err)
	}

	for i := range lamps {
		l.UpdateLampI(i, []*Program{}...)
	}

	return l, nil
}

func (l *Lighting) initDebugCubes() error {
	l.cubeShader = NewProgram()
	if err := l.cubeShader.LinkFiles("assets/shaders/generic_3d.vs", "assets/shaders/uniform_color.fs", ""); err != nil {
		return fmt.Errorf("failed to initialize shader program: %w", err)
	}

	gl.GenVertexArrays(1, &l.cubeVAO)
	gl.BindVertexArray(l.cubeVAO)
	cubeBuffer := NewBuffer(gl.ARRAY_BUFFER)
	cubeBuffer.Bind()
	cubeBuffer.SetVec3(CubeVertices)
	cubeBuffer.LinkVertexPointer(l.cubeShader, "frag_pos", 3, gl.FLOAT, 0, 0)

	return nil
}

func (l *Lighting) initDepthMaps() {
	l.DepthMaps = NewTexture(gl.TEXTURE_CUBE_MAP_ARRAY)
	l.DepthMaps.SetData3D(gl.TEXTURE_CUBE_MAP_ARRAY, 0, gl.DEPTH_COMPONENT, ShadowResolution, ShadowResolution, int32(len(l.lamps)*6), 0, gl.DEPTH_COMPONENT, gl.FLOAT, nil)

	l.DepthMaps.SetIParameter(gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	l.DepthMaps.SetIParameter(gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	l.DepthMaps.SetIParameter(gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	l.DepthMaps.SetIParameter(gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	l.DepthMaps.SetIParameter(gl.TEXTURE_WRAP_R, gl.CLAMP_TO_EDGE)
	l.DepthMaps.SetIParameter(gl.TEXTURE_BASE_LEVEL, 0)
	l.DepthMaps.SetIParameter(gl.TEXTURE_MAX_LEVEL, 0)

	l.depthMapsFBO = NewFramebuffer(gl.FRAMEBUFFER)
	l.depthMapsFBO.Bind()
	l.depthMapsFBO.SetTexture(gl.DEPTH_ATTACHMENT, l.DepthMaps, 0)
	gl.DrawBuffer(gl.NONE)
	gl.ReadBuffer(gl.NONE)
	l.depthMapsFBO.Unbind()
}

// MakeFragShader creates a new fragment shader defined by this Lighting
func (l *Lighting) MakeFragShader() (*Shader, error) {
	fs := NewShader(gl.FRAGMENT_SHADER, l.fragSource)
	if err := fs.Compile(); err != nil {
		return nil, fmt.Errorf("failed to compile: %w", err)
	}

	return fs, nil
}

// ProgramVS is a convenience function which creates a new program with a
// new Lighting fragment shader and the provided vertex shader
func (l *Lighting) ProgramVS(vs *Shader) (*Program, error) {
	fs, err := l.MakeFragShader()
	if err != nil {
		return nil, fmt.Errorf("failed to compile fragment shader: %w", err)
	}

	p := NewProgram()
	if err := p.Link(vs, fs, nil); err != nil {
		return nil, fmt.Errorf("failed to link shaders: %w", err)
	}

	return p, nil
}

// ProgramVSFile is a convenience function which creates a new program with a
// new Lighting fragment shader and the provided vertex shader file
func (l *Lighting) ProgramVSFile(vsFile string) (*Program, error) {
	vs, err := NewShaderFile(gl.VERTEX_SHADER, vsFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load: %w", err)
	}
	if err := vs.Compile(); err != nil {
		return nil, fmt.Errorf("failed to compile: %w", err)
	}

	return l.ProgramVS(vs)
}

// ProgramVSTemplateFile is a convenience function which creates a new program with a
// new Lighting fragment shader and the provided vertex shader file template
func (l *Lighting) ProgramVSTemplateFile(vsTplFile string, tplData interface{}) (*Program, error) {
	vs, err := NewShaderTemplateFile(gl.VERTEX_SHADER, vsTplFile, tplData)
	if err != nil {
		return nil, fmt.Errorf("failed to load: %w", err)
	}
	if err := vs.Compile(); err != nil {
		return nil, fmt.Errorf("failed to compile: %w", err)
	}

	return l.ProgramVS(vs)
}

// DepthProgramVS is a convenience function which creates a new program with
// the depth geometry + fragment shaders and the provided vertex shader
func (l *Lighting) DepthProgramVS(vs *Shader) (*Program, error) {
	fs := NewShader(gl.FRAGMENT_SHADER, l.depthFragSource)
	if err := fs.Compile(); err != nil {
		return nil, fmt.Errorf("failed to compile fragment shader: %w", err)
	}
	gs := NewShader(gl.GEOMETRY_SHADER, l.depthGeoSource)
	if err := gs.Compile(); err != nil {
		return nil, fmt.Errorf("failed to compile geometry shader: %w", err)
	}

	p := NewProgram()
	if err := p.Link(vs, fs, gs); err != nil {
		return nil, fmt.Errorf("failed to link shaders: %w", err)
	}

	return p, nil
}

// DepthProgramVSFile is a convenience function which creates a new program with
// the depth geometry + fragment shaders and the provided vertex shader file
func (l *Lighting) DepthProgramVSFile(vsFile string) (*Program, error) {
	vs, err := NewShaderFile(gl.VERTEX_SHADER, vsFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load: %w", err)
	}
	if err := vs.Compile(); err != nil {
		return nil, fmt.Errorf("failed to compile: %w", err)
	}

	return l.DepthProgramVS(vs)
}

// DepthProgramVSTemplateFile is a convenience function which creates a new program with
// the depth geometry + fragment shaders and the provided vertex shader file template
func (l *Lighting) DepthProgramVSTemplateFile(vsTplFile string, tplData interface{}) (*Program, error) {
	vs, err := NewShaderTemplateFile(gl.VERTEX_SHADER, vsTplFile, tplData)
	if err != nil {
		return nil, fmt.Errorf("failed to load: %w", err)
	}
	if err := vs.Compile(); err != nil {
		return nil, fmt.Errorf("failed to compile: %w", err)
	}

	return l.DepthProgramVS(vs)
}

// SetViewPos sets the view position vector
func (l *Lighting) SetViewPos(pos mgl32.Vec3) {
	l.viewPos = pos
}

// UpdateLampI updates a single lamp by index
func (l *Lighting) UpdateLampI(index int, ps ...*Program) {
	lamp := l.lamps[index]
	base := fmt.Sprintf("lamps[%v]", index)
	for _, p := range ps {
		p.SetUniformFloat32(base+".attenuation.constant", lamp.Attenuation.Constant)
		p.SetUniformFloat32(base+".attenuation.linear", lamp.Attenuation.Linear)
		p.SetUniformFloat32(base+".attenuation.quadratic", lamp.Attenuation.Quadratic)

		p.SetUniformVec3(base+".position", lamp.Position)

		p.SetUniformVec3(base+".ambient", lamp.Ambient)
		p.SetUniformVec3(base+".diffuse", lamp.Diffuse)
		p.SetUniformVec3(base+".specular", lamp.Specular)
	}

	// Calculate transforms for each face of the cubemap
	shadowProj := mgl32.Perspective(mgl32.DegToRad(90), 1, nearPlane, farPlane)
	lp := lamp.Position
	ts := l.lampShadowTransforms

	ts[index*6+0] = shadowProj.Mul4(mgl32.LookAtV(lp, lp.Add(mgl32.Vec3{1, 0, 0}), mgl32.Vec3{0, -1, 0}))
	ts[index*6+1] = shadowProj.Mul4(mgl32.LookAtV(lp, lp.Add(mgl32.Vec3{-1, 0, 0}), mgl32.Vec3{0, -1, 0}))
	ts[index*6+2] = shadowProj.Mul4(mgl32.LookAtV(lp, lp.Add(mgl32.Vec3{0, 1, 0}), mgl32.Vec3{0, 0, 1}))
	ts[index*6+3] = shadowProj.Mul4(mgl32.LookAtV(lp, lp.Add(mgl32.Vec3{0, -1, 0}), mgl32.Vec3{0, 0, -1}))
	ts[index*6+4] = shadowProj.Mul4(mgl32.LookAtV(lp, lp.Add(mgl32.Vec3{0, 0, 1}), mgl32.Vec3{0, -1, 0}))
	ts[index*6+5] = shadowProj.Mul4(mgl32.LookAtV(lp, lp.Add(mgl32.Vec3{0, 0, -1}), mgl32.Vec3{0, -1, 0}))

	l.depthUpdateLamps[index] = true
}

// UpdateLamps updates all lamps (expensive)
func (l *Lighting) UpdateLamps(ps ...*Program) {
	for i := range l.lamps {
		l.UpdateLampI(i, ps...)
	}
}

// UpdateLamp updates a single lamp (for the shaders provided), and also marks
// it as needing an update on the next depth map pass
func (l *Lighting) UpdateLamp(lamp *Lamp, ps ...*Program) {
	index := -1
	for i, l := range l.lamps {
		if l == lamp {
			index = i
			break
		}
	}
	if index == -1 {
		return
	}

	l.UpdateLampI(index, ps...)
}

// Update re-sets all of the light parameter uniforms (except for point lamps,
// which should be updated individually for performance)
func (l *Lighting) Update(ps ...*Program) {
	for _, p := range ps {
		p.SetUniformVec3("view_pos", l.viewPos)
		p.SetUniformFloat32("far_plane", farPlane)
		p.SetUniformBool("shadows_enabled", l.ShadowsEnabled)

		for i, spot := range l.spotlights {
			base := fmt.Sprintf("spotlights[%v]", i)

			p.SetUniformFloat32(base+".attenuation.constant", spot.Attenuation.Constant)
			p.SetUniformFloat32(base+".attenuation.linear", spot.Attenuation.Linear)
			p.SetUniformFloat32(base+".attenuation.quadratic", spot.Attenuation.Quadratic)

			p.SetUniformVec3(base+".position", spot.Position)
			p.SetUniformVec3(base+".direction", spot.Direction)

			p.SetUniformFloat32(base+".cutoff", spot.Cutoff)
			p.SetUniformFloat32(base+".outer_cutoff", spot.OuterCutoff)

			p.SetUniformVec3(base+".ambient", spot.Ambient)
			p.SetUniformVec3(base+".diffuse", spot.Diffuse)
			p.SetUniformVec3(base+".specular", spot.Specular)
		}

		for i, dir := range l.dirs {
			base := fmt.Sprintf("dirs[%v]", i)

			p.SetUniformVec3(base+".direction", dir.Direction)

			p.SetUniformVec3(base+".ambient", dir.Ambient)
			p.SetUniformVec3(base+".diffuse", dir.Diffuse)
			p.SetUniformVec3(base+".specular", dir.Specular)
		}
	}
}

// DrawCubes renders the lights in the scene as cubes (coloured by their diffuse
// colour)
func (l *Lighting) DrawCubes(projection mgl32.Mat4, c *Camera) {
	l.cubeShader.Use()
	for _, lamp := range l.lamps {
		l.cubeShader.Project(projection, c, TransFromPos(lamp.Position).Mul4(mgl32.Scale3D(0.4, 0.4, 0.4)))
		l.cubeShader.SetUniformVec3("color", lamp.Diffuse)

		gl.BindVertexArray(l.cubeVAO)
		gl.DrawArrays(gl.TRIANGLES, 0, 36)
	}
}

// DepthMapParamsApplicator is a function passed to the DepthMapRenderFunc which
// should be called by the user to set uniforms used by the depth mapping shader
type DepthMapParamsApplicator = func(p *Program)

// DepthMapRenderFunc is a user function which will be called to actually draw
// the geometry once the framebuffer is set up
type DepthMapRenderFunc = func(DepthMapParamsApplicator)

// ShadowsDepthPass renders the scene from each lamp's perspective to generate
// a set of depth maps
func (l *Lighting) ShadowsDepthPass(cb DepthMapRenderFunc) {
	if !l.ShadowsEnabled {
		return
	}

	gl.Viewport(0, 0, ShadowResolution, ShadowResolution)
	l.depthMapsFBO.Bind()
	gl.Clear(gl.DEPTH_BUFFER_BIT)

	lampPositions := make([]mgl32.Vec3, len(l.lamps))
	for i, lamp := range l.lamps {
		lampPositions[i] = lamp.Position
	}

	cb(func(p *Program) {
		p.Use()

		// For geometry shader
		p.SetUniformBoolSlice("update_lamps", l.depthUpdateLamps)
		p.SetUniformMat4Slice("shadow_transforms", l.lampShadowTransforms)

		// For fragment shader
		p.SetUniformVec3Slice("lamp_positions", lampPositions)
		p.SetUniformFloat32("far_plane", farPlane)
	})

	for i := range l.depthUpdateLamps {
		l.depthUpdateLamps[i] = false
	}

	l.depthMapsFBO.Unbind()
}
