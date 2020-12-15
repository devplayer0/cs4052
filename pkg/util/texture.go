package util

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/png"

	"github.com/go-gl/gl/v4.6-core/gl"
)

// Texture represents an OpenGL texture
type Texture struct {
	ty uint32

	id uint32
}

// NewTexture creates a new OpenGL texture
func NewTexture(ty uint32) *Texture {
	t := &Texture{ty: ty}
	gl.GenTextures(1, &t.id)

	return t
}

// Bind binds the texture
func (t *Texture) Bind() {
	gl.BindTexture(t.ty, t.id)
}

// Activate sets up the provided texture unit, binds the texture and sets the
// uniform corresponding to the sampler in the shader
func (t *Texture) Activate(p *Program, uniform string, unit uint32) {
	gl.ActiveTexture(gl.TEXTURE0 + unit)
	t.Bind()
	p.SetUniformInt(uniform, int32(unit))
}

// GenerateMipmap generates a mipmap for the texture
func (t *Texture) GenerateMipmap() {
	t.Bind()
	gl.GenerateMipmap(t.ty)
}

// SetIParameter allows for parameters such as texture wrapping to be set
func (t *Texture) SetIParameter(param uint32, value int32) {
	t.Bind()
	gl.TexParameteri(t.ty, param, value)
}

// SetData2D uploads 2D pixel data to the texture on the GPU
func (t *Texture) SetData2D(level, internalformat, width, height, border int32, format, xtype uint32, pixels []byte) {
	t.Bind()
	gl.TexImage2D(t.ty, level, internalformat, width, height, border, format, xtype, gl.Ptr(pixels))
}

// LoadPNG decodes a PNG and uploads it to the texture on the GPU
func (t *Texture) LoadPNG(data []byte) error {
	img, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("failed to decode: %w", err)
	}

	var pixels []byte
	switch img.ColorModel() {
	case color.RGBAModel:
		pixels = img.(*image.RGBA).Pix
	case color.NRGBAModel:
		pixels = img.(*image.NRGBA).Pix
	default:
		return errors.New("unknown color model")
	}

	b := img.Bounds()
	size := b.Max.Sub(b.Min)
	t.SetData2D(0, gl.RGBA, int32(size.X), int32(size.Y), 0, gl.RGBA, gl.UNSIGNED_BYTE, pixels)
	return nil
}
