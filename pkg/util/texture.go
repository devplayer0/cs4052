package util

import (
	"bytes"
	"fmt"
	"image/jpeg"
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
func (t *Texture) SetData2D(target uint32, level, internalformat, width, height, border int32, format, xtype uint32, pixels []byte) {
	t.Bind()
	gl.TexImage2D(target, level, internalformat, width, height, border, format, xtype, gl.Ptr(pixels))
}

// LoadPNG decodes a PNG and uploads it to the texture on the GPU
func (t *Texture) LoadPNG(target uint32, data []byte) error {
	img, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("failed to decode: %w", err)
	}

	b := img.Bounds()
	size := b.Max.Sub(b.Min)
	pixels := make([]byte, size.X*size.Y*4)
	for x := 0; x < size.X; x++ {
		for y := 0; y < size.Y; y++ {
			r, g, b, a := img.At(x, y).RGBA()
			pixels[0+(x*4)+(y*4*size.X)] = byte(r >> 8)
			pixels[1+(x*4)+(y*4*size.X)] = byte(g >> 8)
			pixels[2+(x*4)+(y*4*size.X)] = byte(b >> 8)
			pixels[3+(x*4)+(y*4*size.X)] = byte(a >> 8)
		}
	}

	t.SetData2D(target, 0, gl.RGBA, int32(size.X), int32(size.Y), 0, gl.RGBA, gl.UNSIGNED_BYTE, pixels)
	return nil
}

// LoadJPEG decodes a JPEG and uploads it to the texture on the GPU
func (t *Texture) LoadJPEG(target uint32, data []byte) error {
	img, err := jpeg.Decode(bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("failed to decode: %w", err)
	}

	b := img.Bounds()
	size := b.Max.Sub(b.Min)
	pixels := make([]byte, size.X*size.Y*3)
	for x := 0; x < size.X; x++ {
		for y := 0; y < size.Y; y++ {
			r, g, b, _ := img.At(x, y).RGBA()
			pixels[0+(x*3)+(y*3*size.X)] = byte(r >> 8)
			pixels[1+(x*3)+(y*3*size.X)] = byte(g >> 8)
			pixels[2+(x*3)+(y*3*size.X)] = byte(b >> 8)
		}
	}

	t.SetData2D(target, 0, gl.RGB, int32(size.X), int32(size.Y), 0, gl.RGB, gl.UNSIGNED_BYTE, pixels)
	return nil
}
