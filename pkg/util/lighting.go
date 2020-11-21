package util

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"text/template"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

const lightingFragShaderFile = "assets/shaders/lighting.fs"

// AttenuationParams represents the lighting attenuation coefficients
type AttenuationParams struct {
	Constant  float32
	Linear    float32
	Quadratic float32
}

// Lamp represents a point light source
type Lamp struct {
	Attenuation AttenuationParams

	Position mgl32.Vec3

	Ambient  mgl32.Vec3
	Diffuse  mgl32.Vec3
	Specular mgl32.Vec3
}

// Lighting represents a shader to colour an object with lighting
type Lighting struct {
	lamps      []*Lamp
	cubeVAO    uint32
	cubeShader *Program

	Shader *Program
}

func makeLightingFS(lamps []*Lamp) (*Shader, error) {
	fsTemplateData, err := ioutil.ReadFile(lightingFragShaderFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read source file %v: %w", lightingFragShaderFile, err)
	}
	fsTemplate, err := template.New(lightingFragShaderFile).Parse(string(fsTemplateData))
	if err != nil {
		return nil, fmt.Errorf("failed to parse template: %w", err)
	}

	fsSourceBuf := &bytes.Buffer{}
	if err := fsTemplate.Execute(fsSourceBuf, struct {
		Lamps []*Lamp
	}{
		Lamps: lamps,
	}); err != nil {
		return nil, fmt.Errorf("failed to execute template: %w", err)
	}

	fs := NewShader(gl.FRAGMENT_SHADER, fsSourceBuf.String())
	if err := fs.Compile(); err != nil {
		return nil, fmt.Errorf("failed to compile shader: %w", err)
	}

	return fs, nil
}

// NewLighting creates a new lighting shader from a given vertex shader and set
// of lamps
func NewLighting(vs *Shader, lamps []*Lamp) (*Lighting, error) {
	fs, err := makeLightingFS(lamps)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize fragment shader: %w", err)
	}
	prog := NewProgram()
	if err := prog.Link(vs, fs); err != nil {
		return nil, fmt.Errorf("failed to link shaders to program: %w", err)
	}

	cubeProg := NewProgram()
	if err := cubeProg.LinkFiles("assets/shaders/generic_3d.vs", "assets/shaders/uniform_color.fs"); err != nil {
		return nil, fmt.Errorf("failed to initialize lamp cube program: %w", err)
	}

	l := &Lighting{
		lamps:      lamps,
		cubeShader: cubeProg,

		Shader: prog,
	}

	gl.GenVertexArrays(1, &l.cubeVAO)
	gl.BindVertexArray(l.cubeVAO)
	cubeBuffer := NewBuffer(gl.ARRAY_BUFFER)
	cubeBuffer.Bind()
	cubeBuffer.SetVec3(cubeVertices)
	cubeBuffer.LinkVertexPointer(cubeProg, "frag_pos", 3, gl.FLOAT, 0, 0)

	l.Update()

	return l, nil
}

// NewLightingVSFile is a convenience function that creates a new lighting
// shader by compiling the vertex shader source file
func NewLightingVSFile(vsFile string, lamps []*Lamp) (*Lighting, error) {
	vs, err := NewShaderFile(gl.VERTEX_SHADER, vsFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read vertex shader source: %w", err)
	}
	if err := vs.Compile(); err != nil {
		return nil, fmt.Errorf("failed to compile vertex shader: %w", err)
	}

	return NewLighting(vs, lamps)
}

// Update re-sets all of the lamp parameter uniforms
func (l *Lighting) Update() {
	for i, lamp := range l.lamps {
		base := fmt.Sprintf("lamps[%v]", i)

		l.Shader.SetUniformFloat32(base+".attenuation.constant", lamp.Attenuation.Constant)
		l.Shader.SetUniformFloat32(base+".attenuation.linear", lamp.Attenuation.Linear)
		l.Shader.SetUniformFloat32(base+".attenuation.quadratic", lamp.Attenuation.Quadratic)

		l.Shader.SetUniformVec3(base+".position", lamp.Position)
		l.Shader.SetUniformVec3(base+".ambient", lamp.Ambient)
		l.Shader.SetUniformVec3(base+".diffuse", lamp.Diffuse)
		l.Shader.SetUniformVec3(base+".specular", lamp.Specular)
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
