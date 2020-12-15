package util

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"text/template"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

const lightingFragShaderFile = "assets/shaders/lighting_nmap.fs"

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

type fsTemplateData struct {
	Dirs       []*DirectionalLight
	Lamps      []*Lamp
	Spotlights []*Spotlight
}

// Lighting represents a shader to colour an object with lighting
type Lighting struct {
	dirs       []*DirectionalLight
	lamps      []*Lamp
	spotlights []*Spotlight

	cubeVAO    uint32
	cubeShader *Program

	fragSource string
}

// NewLighting creates a new lighting shader from a given vertex shader and set
// of lamps
func NewLighting(dirs []*DirectionalLight, lamps []*Lamp, spotlights []*Spotlight) (*Lighting, error) {
	fsSourceData, err := ioutil.ReadFile(lightingFragShaderFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read fragment shader source file %v: %w", lightingFragShaderFile, err)
	}
	fsSourceTpl, err := template.New(lightingFragShaderFile).Parse(string(fsSourceData))
	if err != nil {
		return nil, fmt.Errorf("failed to parse fragment shader template: %w", err)
	}
	fsSourceBuf := &bytes.Buffer{}
	if err := fsSourceTpl.Execute(fsSourceBuf, fsTemplateData{dirs, lamps, spotlights}); err != nil {
		return nil, fmt.Errorf("failed to execute template: %w", err)
	}

	cubeProg := NewProgram()
	if err := cubeProg.LinkFiles("assets/shaders/generic_3d.vs", "assets/shaders/uniform_color.fs"); err != nil {
		return nil, fmt.Errorf("failed to initialize lamp cube program: %w", err)
	}

	l := &Lighting{
		dirs:       dirs,
		lamps:      lamps,
		spotlights: spotlights,

		cubeShader: cubeProg,

		fragSource: fsSourceBuf.String(),
	}

	gl.GenVertexArrays(1, &l.cubeVAO)
	gl.BindVertexArray(l.cubeVAO)
	cubeBuffer := NewBuffer(gl.ARRAY_BUFFER)
	cubeBuffer.Bind()
	cubeBuffer.SetVec3(CubeVertices)
	cubeBuffer.LinkVertexPointer(cubeProg, "frag_pos", 3, gl.FLOAT, 0, 0)

	return l, nil
}

// MakeFragShader creates a new fragment shader defined by this Lighting
func (l *Lighting) MakeFragShader() (*Shader, error) {
	fs := NewShader(gl.FRAGMENT_SHADER, l.fragSource)
	if err := fs.Compile(); err != nil {
		return nil, fmt.Errorf("failed to compile: %w", err)
	}

	return fs, nil
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

	fs, err := l.MakeFragShader()
	if err != nil {
		return nil, fmt.Errorf("failed to compile fragment shader: %w", err)
	}

	p := NewProgram()
	if err := p.Link(vs, fs); err != nil {
		return nil, fmt.Errorf("failed to link shaders: %w", err)
	}

	return p, nil
}

// Update re-sets all of the lamp parameter uniforms
func (l *Lighting) Update(ps ...*Program) {
	for _, p := range ps {
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
		for i, lamp := range l.lamps {
			base := fmt.Sprintf("lamps[%v]", i)

			p.SetUniformFloat32(base+".attenuation.constant", lamp.Attenuation.Constant)
			p.SetUniformFloat32(base+".attenuation.linear", lamp.Attenuation.Linear)
			p.SetUniformFloat32(base+".attenuation.quadratic", lamp.Attenuation.Quadratic)

			p.SetUniformVec3(base+".position", lamp.Position)

			p.SetUniformVec3(base+".ambient", lamp.Ambient)
			p.SetUniformVec3(base+".diffuse", lamp.Diffuse)
			p.SetUniformVec3(base+".specular", lamp.Specular)
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
