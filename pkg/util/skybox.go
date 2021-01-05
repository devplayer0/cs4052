package util

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

var cubeSides = map[string]uint32{
	"right":  gl.TEXTURE_CUBE_MAP_POSITIVE_X,
	"left":   gl.TEXTURE_CUBE_MAP_NEGATIVE_X,
	"top":    gl.TEXTURE_CUBE_MAP_POSITIVE_Y,
	"bottom": gl.TEXTURE_CUBE_MAP_NEGATIVE_Y,
	"front":  gl.TEXTURE_CUBE_MAP_POSITIVE_Z,
	"back":   gl.TEXTURE_CUBE_MAP_NEGATIVE_Z,
}

// Skybox represents a cube-mapped skybox
type Skybox struct {
	Texture *Texture
	shader  *Program
	vao     uint32
}

// NewSkybox creates a new skybox
func NewSkybox(pathBase string) (*Skybox, error) {
	t := NewTexture(gl.TEXTURE_CUBE_MAP)
	t.Bind()

	for side, glTarget := range cubeSides {
		path := pathBase + side + ".jpg"
		if err := t.LoadJPEGFile(glTarget, path); err != nil {
			return nil, fmt.Errorf("failed to upload %v texture to GPU: %w", side, err)
		}
	}

	t.SetIParameter(gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	t.SetIParameter(gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	t.SetIParameter(gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	t.SetIParameter(gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	t.SetIParameter(gl.TEXTURE_WRAP_R, gl.CLAMP_TO_EDGE)

	shader := NewProgram()
	if err := shader.LinkFiles("assets/shaders/skybox.vs", "assets/shaders/skybox.fs", ""); err != nil {
		return nil, fmt.Errorf("failed to initialize shader: %w", err)
	}

	s := &Skybox{
		Texture: t,
		shader:  shader,
	}
	gl.GenVertexArrays(1, &s.vao)

	gl.BindVertexArray(s.vao)
	buf := &bytes.Buffer{}
	binary.Write(buf, NativeOrder, CubeVertices)

	vertexBuf := NewBuffer(gl.ARRAY_BUFFER)
	vertexBuf.SetData(buf.Bytes())
	vertexBuf.LinkVertexPointer(s.shader, "frag_pos", 3, gl.FLOAT, 12, 0)

	return s, nil
}

// Draw renders the skybox (should be done last)
func (s *Skybox) Draw(proj mgl32.Mat4, c *Camera) {
	// Use LEQUAL depth test function so depth test passes when values are equal
	// to the buffer's content (see also vertex shader)
	gl.DepthFunc(gl.LEQUAL)

	s.Texture.Activate(s.shader, "skybox", 0)

	s.shader.Use()
	s.shader.SetUniformMat4("projection", proj)
	// Convert the camera's transform to a Mat3 and back to discard translation
	s.shader.SetUniformMat4("camera", c.Transform().Mat3().Mat4())

	gl.BindVertexArray(s.vao)
	gl.DrawArrays(gl.TRIANGLES, 0, 36)

	gl.DepthFunc(gl.LESS)
}
