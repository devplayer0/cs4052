package util

import (
	"github.com/go-gl/gl/v4.6-core/gl"
)

// Framebuffer represents an OpenGL framebuffer
type Framebuffer struct {
	target uint32

	id uint32
}

// NewFramebuffer creates a new OpenGL framebuffer
func NewFramebuffer(target uint32) *Framebuffer {
	f := &Framebuffer{target: target}
	gl.GenFramebuffers(1, &f.id)

	return f
}

// Bind binds the framebuffer
func (f *Framebuffer) Bind() {
	gl.BindFramebuffer(f.target, f.id)
}

// Unbind sets the currently bound framebuffer to 0
func (f *Framebuffer) Unbind() {
	gl.BindFramebuffer(f.target, 0)
}

// SetTexture sets the texture associated with the framebuffer
func (f *Framebuffer) SetTexture(attachment uint32, t *Texture, level int32) {
	f.Bind()
	gl.FramebufferTexture(f.target, attachment, t.id, level)
}

// SetTextureLayer sets the texture associated with a specific layer in the framebuffer
func (f *Framebuffer) SetTextureLayer(attachment uint32, t *Texture, level, layer int32) {
	f.Bind()
	gl.FramebufferTextureLayer(f.target, attachment, t.id, level, layer)
}
