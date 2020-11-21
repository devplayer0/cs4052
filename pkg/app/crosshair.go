package app

import (
	"fmt"

	"github.com/devplayer0/cs4052/pkg/util"
	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

// Crosshair represents a centred 2D crosshair
type Crosshair struct {
	vao    uint32
	shader *util.Program
}

// NewCrosshair creates a new crosshair object
func NewCrosshair(win *glfw.Window) (*Crosshair, error) {
	c := &Crosshair{}

	c.shader = util.NewProgram()
	if err := c.shader.LinkFiles("assets/shaders/crosshair.vs", "assets/shaders/white.fs"); err != nil {
		return nil, fmt.Errorf("failed to set up program: %w", err)
	}

	wi, hi := win.GetSize()
	w := float32(wi)
	h := float32(hi)
	c.shader.SetUniformMat4("projection", mgl32.Ortho2D(0, w, 0, h))
	c.shader.SetUniformMat4("model", mgl32.Translate3D(w/2, h/2, 0).Mul4(mgl32.Scale3D(8, 8, 0)))

	gl.GenVertexArrays(1, &c.vao)
	gl.BindVertexArray(c.vao)
	vertexBuf := util.NewBuffer(gl.ARRAY_BUFFER)
	vertexBuf.SetVec2([]mgl32.Vec2{
		{-1, 0},
		{1, 0},

		{0, -1},
		{0, 1},
	})
	vertexBuf.LinkVertexPointer(c.shader, "frag_pos", 2, gl.FLOAT, 0, 0)

	return c, nil
}

// Draw renders the crosshair on screen
func (c *Crosshair) Draw() {
	c.shader.Use()
	gl.BindVertexArray(c.vao)
	gl.DrawArrays(gl.LINES, 0, 4)
}
