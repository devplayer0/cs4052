package main

import (
	"errors"
	"fmt"
	"log"
	"runtime"
	"strings"
	"unsafe"

	"github.com/MakeNowJust/heredoc"
	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

var window *glfw.Window
var program, vao uint32
var previousTime float64
var frames uint32

var vertexShader = heredoc.Doc(`
	#version 400

	in vec3 vPosition;
	in vec4 vColor;

	out vec4 color;

	void main() {
		gl_Position = vec4(vPosition.x, vPosition.y, vPosition.z, 1.0);
		color = vColor;
	}
`) + "\x00"
var fragmentShader = heredoc.Doc(`
	#version 400

	in vec4 color;
	out vec4 outColor;

	void main() {
		outColor = color;
	}
`) + "\x00"

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()

	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var len int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &len)

		log := strings.Repeat("\x00", int(len+1))
		gl.GetShaderInfoLog(shader, len, nil, gl.Str(log))

		return 0, errors.New(log)
	}

	return shader, nil
}

func createProgram(vertexSource, fragmentSource string) (uint32, error) {
	vertex, err := compileShader(vertexSource, gl.VERTEX_SHADER)
	if err != nil {
		return 0, fmt.Errorf("failed to compile vertex shader: %v", err)
	}

	fragment, err := compileShader(fragmentSource, gl.FRAGMENT_SHADER)
	if err != nil {
		return 0, fmt.Errorf("failed to compile fragment shader: %v", err)
	}

	program := gl.CreateProgram()
	gl.AttachShader(program, vertex)
	gl.AttachShader(program, fragment)
	gl.LinkProgram(program)

	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var len int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &len)

		log := strings.Repeat("\x00", int(len+1))
		gl.GetProgramInfoLog(program, len, nil, gl.Str(log))

		return 0, errors.New(log)
	}

	gl.ValidateProgram(program)
	if status == gl.FALSE {
		var len int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &len)

		log := strings.Repeat("\x00", int(len+1))
		gl.GetProgramInfoLog(program, len, nil, gl.Str(log))

		return 0, errors.New(log)
	}

	return program, nil
}

func generateObjectBuffer(program uint32, vertices []float32, colors []float32) uint32 {
	gl.UseProgram(program)

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)

	// (vertices + colours) * 4 (sizeof(float32))
	gl.BufferData(gl.ARRAY_BUFFER, (len(vertices)+len(colors))*4, nil, gl.STATIC_DRAW)
	gl.BufferSubData(gl.ARRAY_BUFFER, 0, len(vertices)*4, gl.Ptr(vertices))
	gl.BufferSubData(gl.ARRAY_BUFFER, len(vertices)*4, len(colors)*4, gl.Ptr(colors))
	return vbo
}

func linkBufferToProgram(program, buffer uint32, colorOffset int) {
	gl.UseProgram(program)
	gl.BindBuffer(gl.ARRAY_BUFFER, buffer)

	positionAttrib := uint32(gl.GetAttribLocation(program, gl.Str("vPosition\x00")))
	gl.EnableVertexAttribArray(positionAttrib)
	gl.VertexAttribPointer(positionAttrib, 3, gl.FLOAT, false, 0, gl.PtrOffset(0))

	colorAttrib := uint32(gl.GetAttribLocation(program, gl.Str("vColor\x00")))
	gl.EnableVertexAttribArray(colorAttrib)
	gl.VertexAttribPointer(colorAttrib, 4, gl.FLOAT, false, 0, gl.PtrOffset(colorOffset))
}

func setup() error {
	var err error
	program, err = createProgram(vertexShader, fragmentShader)
	if err != nil {
		return fmt.Errorf("failed to create program: %v", err)
	}

	vertices := []float32{
		-1, -1, 0,
		1, -1, 0,
		0, 1, 0,
	}
	colors := []float32{
		0, 1, 0, 1,
		1, 0, 0, 1,
		0, 0, 1, 1,
	}

	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	vbo := generateObjectBuffer(program, vertices, colors)
	linkBufferToProgram(program, vbo, len(vertices)*4)

	return nil
}

func loop() {
	t := glfw.GetTime()
	if t-previousTime > 1 {
		log.Printf("FPS: %v", frames)

		frames = 0
		previousTime = t
	}

	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	gl.UseProgram(program)
	gl.BindVertexArray(vao)
	gl.DrawArrays(gl.TRIANGLES, 0, 3)

	window.SwapBuffers()
	glfw.PollEvents()
	frames++
}

func init() {
	// Keep overything on main thread for OpenGL stuff
	runtime.LockOSThread()
}

func main() {
	if err := glfw.Init(); err != nil {
		log.Fatalf("Error initializing glfw: %v", err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 6)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	var err error
	window, err = glfw.CreateWindow(800, 600, "Hello Triangle", nil, nil)
	if err != nil {
		log.Fatalf("Failed to create window: %v", err)
	}

	window.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		log.Fatalf("Failed to initialize OpenGL: %v", err)
	}

	gl.Enable(gl.DEBUG_OUTPUT)
	gl.DebugMessageCallback(func(source, gltype, id, severity uint32, length int32, message string, userParam unsafe.Pointer) {
		log.Printf("[SEV%v] %v", severity, message)
	}, nil)

	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Printf("OpenGL version: %v", version)

	if err := setup(); err != nil {
		log.Fatalf("Setup failed: %v", err)
	}

	previousTime = glfw.GetTime()
	for !window.ShouldClose() {
		loop()
	}
}
