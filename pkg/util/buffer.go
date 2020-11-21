package util

import (
	"bytes"
	"encoding/binary"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

// TODO: Determine actual system byte order
var nativeOrder = binary.LittleEndian

// Buffer represents an OpenGL buffer
type Buffer struct {
	t uint32

	id uint32
}

// NewBuffer creates a new OpenGL VBO
func NewBuffer(t uint32) *Buffer {
	b := &Buffer{t: t}
	gl.GenBuffers(1, &b.id)

	return b
}

// Bind binds the buffer
func (b *Buffer) Bind() {
	gl.BindBuffer(b.t, b.id)
}

// SetData sets the buffers's data
func (b *Buffer) SetData(data []byte) {
	b.Bind()
	gl.BufferData(b.t, len(data), gl.Ptr(data), gl.STATIC_DRAW)
}

// LinkVertexPointer sets up a named vertex attribute to point to a buffer
func (b *Buffer) LinkVertexPointer(p *Program, va string, size int32, vType uint32, stride int32, offset int) {
	b.Bind()

	attrib := uint32(gl.GetAttribLocation(p.ID, gl.Str(va+"\x00")))
	gl.EnableVertexAttribArray(attrib)
	gl.VertexAttribPointer(attrib, size, vType, false, stride, gl.PtrOffset(offset))
}

// SetVec2 sets the VBO's data to the list of Vec2 (e.g. vertices)
func (b *Buffer) SetVec2(vertices []mgl32.Vec2) {
	buf := &bytes.Buffer{}
	for _, v := range vertices {
		binary.Write(buf, nativeOrder, v.X())
		binary.Write(buf, nativeOrder, v.Y())
	}

	b.SetData(buf.Bytes())
}

// SetVec3 sets the VBO's data to the list of Vec3 (e.g. vertices)
func (b *Buffer) SetVec3(vertices []mgl32.Vec3) {
	buf := &bytes.Buffer{}
	for _, v := range vertices {
		binary.Write(buf, nativeOrder, v.X())
		binary.Write(buf, nativeOrder, v.Y())
		binary.Write(buf, nativeOrder, v.Z())
	}

	b.SetData(buf.Bytes())
}

// SetVec4 sets the VBO's data to the list of Vec4 (e.g. colors)
func (b *Buffer) SetVec4(vecs []mgl32.Vec4) {
	buf := &bytes.Buffer{}
	for _, v := range vecs {
		binary.Write(buf, nativeOrder, v.X())
		binary.Write(buf, nativeOrder, v.Y())
		binary.Write(buf, nativeOrder, v.Z())
		binary.Write(buf, nativeOrder, v.W())
	}

	b.SetData(buf.Bytes())
}
