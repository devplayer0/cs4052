package util

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

// Program represents an OpenGL program
type Program struct {
	ID uint32

	Vertex   *Shader
	Fragment *Shader

	uniforms map[string]int32
}

// NewProgram creates a new OpenGL program
func NewProgram() *Program {
	p := &Program{
		ID: gl.CreateProgram(),

		uniforms: make(map[string]int32),
	}

	return p
}

// Link attaches a vertex and fragment shader and links the OpenGL program
func (p *Program) Link(vertex, fragment *Shader) error {
	gl.AttachShader(p.ID, vertex.ID)
	gl.AttachShader(p.ID, fragment.ID)
	gl.LinkProgram(p.ID)

	var status int32
	gl.GetProgramiv(p.ID, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var len int32
		gl.GetProgramiv(p.ID, gl.INFO_LOG_LENGTH, &len)

		log := strings.Repeat("\x00", int(len+1))
		gl.GetProgramInfoLog(p.ID, len, nil, gl.Str(log))

		return errors.New(log)
	}

	gl.ValidateProgram(p.ID)
	if status == gl.FALSE {
		var len int32
		gl.GetProgramiv(p.ID, gl.INFO_LOG_LENGTH, &len)

		log := strings.Repeat("\x00", int(len+1))
		gl.GetProgramInfoLog(p.ID, len, nil, gl.Str(log))

		return errors.New(log)
	}

	return nil
}

// LinkFiles is a shortcut to load and compile a vertex and fragment shader from file
func (p *Program) LinkFiles(vertex, fragment string) error {
	v, err := NewShaderFile(gl.VERTEX_SHADER, vertex)
	if err != nil {
		return fmt.Errorf("failed to load vertex shader: %w", err)
	}
	if err := v.Compile(); err != nil {
		return fmt.Errorf("failed to compile vertex shader: %w", err)
	}

	f, err := NewShaderFile(gl.FRAGMENT_SHADER, fragment)
	if err != nil {
		return fmt.Errorf("failed to load fragment shader: %w", err)
	}
	if err := f.Compile(); err != nil {
		return fmt.Errorf("failed to compile fragment shader: %w", err)
	}

	return p.Link(v, f)
}

// Use sets up OpenGL to use this program
func (p *Program) Use() {
	gl.UseProgram(p.ID)
}

// Uniform gets a uniform's location
func (p *Program) Uniform(n string) int32 {
	p.Use()

	u, ok := p.uniforms[n]
	if ok {
		return u
	}

	u = gl.GetUniformLocation(p.ID, gl.Str(n+"\x00"))
	p.uniforms[n] = u
	return u
}

// SetUniformFloat32 sets a float32 uniform value
func (p *Program) SetUniformFloat32(n string, val float32) {
	gl.Uniform1fv(p.Uniform(n), 1, &val)
}

// SetUniformVec3 sets a vec3 uniform value
func (p *Program) SetUniformVec3(n string, val mgl32.Vec3) {
	gl.Uniform3fv(p.Uniform(n), 1, &val[0])
}

// SetUniformMat3 sets a mat3 uniform value
func (p *Program) SetUniformMat3(n string, val mgl32.Mat3) {
	gl.UniformMatrix3fv(p.Uniform(n), 1, false, &val[0])
}

// SetUniformMat4 sets a mat4 uniform value
func (p *Program) SetUniformMat4(n string, val mgl32.Mat4) {
	gl.UniformMatrix4fv(p.Uniform(n), 1, false, &val[0])
}

// Project sets the uniforms for the required 3D transformation matrices
func (p *Program) Project(projection mgl32.Mat4, c *Camera, model mgl32.Mat4) {
	p.SetUniformMat4("projection", projection)
	p.SetUniformMat4("camera", c.Transform())
	p.SetUniformMat4("model", model)
}
