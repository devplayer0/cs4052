package util

import (
	"bytes"
	"encoding/binary"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

// TODO: Determine actual system byte order
var nativeOrder = binary.LittleEndian

// Buffer represents an OpenGL VBO
type Buffer struct {
	ID uint32
}

// NewBuffer creates a new OpenGL VBO
func NewBuffer() *Buffer {
	b := &Buffer{}
	gl.GenBuffers(1, &b.ID)

	return b
}

// Bind binds the buffer
func (b *Buffer) Bind() {
	gl.BindBuffer(gl.ARRAY_BUFFER, b.ID)
}

// SetData sets the VBO's data
func (b *Buffer) SetData(data []byte) {
	b.Bind()
	gl.BufferData(gl.ARRAY_BUFFER, len(data), gl.Ptr(data), gl.STATIC_DRAW)
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