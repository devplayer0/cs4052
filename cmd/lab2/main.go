package main

import (
	"fmt"
	"log"
	"runtime"
	"unsafe"

	"github.com/devplayer0/cs4052/pkg/util"
	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

var window *glfw.Window

var colorProg, yellowProg *util.Program

func setup() error {
	colorProg = util.NewProgram()
	if err := colorProg.LinkFiles("assets/shaders/color_tri.vs", "assets/shaders/color_tri.fs"); err != nil {
		return fmt.Errorf("failed to set up color triangle program: %w", err)
	}

	vertexBuf := util.NewBuffer()
	vertexBuf.SetVertices([]mgl32.Vec3{
		{-0.9, -0.9, 0},
		{0.9, -0.9, 0},
		{0, -0.1, 0},

		{0.9, 0.9, 0},
		{-0.9, 0.9, 0},
		{0, 0.1, 0},
	})
	colorProg.LinkVertexPointer("vPosition", 3, gl.FLOAT, vertexBuf, 0)

	colorBuf := util.NewBuffer()
	colorBuf.SetVec4([]mgl32.Vec4{
		{0, 1, 0, 1},
		{1, 0, 0, 1},
		{0, 0, 1, 1},

		{1, 0, 0, 1},
		{0, 0, 1, 1},
		{0, 1, 0, 1},
	})
	colorProg.LinkVertexPointer("vColor", 4, gl.FLOAT, colorBuf, 0)

	yellowProg = util.NewProgram()
	if err := yellowProg.LinkFiles("assets/shaders/yellow_tri.vs", "assets/shaders/yellow_tri.fs"); err != nil {
		return fmt.Errorf("failed to set up yellow triangle program: %w", err)
	}

	vertexBuf = util.NewBuffer()
	vertexBuf.SetVertices([]mgl32.Vec3{
		{-1, -1, 0},
		{1, -1, 0},
		{0, 0, 0},
	})
	yellowProg.LinkVertexPointer("vPosition", 3, gl.FLOAT, vertexBuf, 0)

	return nil
}

func loop() {
	frames := 0
	lastFPS := glfw.GetTime()
	previousTime := lastFPS
	var angle float32

	window.SetKeyCallback(func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		switch key {
		case glfw.KeyLeft:
			angle -= 5
		case glfw.KeyRight:
			angle += 5
		}
	})

	for !window.ShouldClose() {
		t := glfw.GetTime()
		if t-lastFPS > 1 {
			log.Printf("FPS: %v", frames)

			frames = 0
			lastFPS = t
		}

		_ = t - previousTime
		if angle >= 360 {
			angle = 0
		}
		yTrans := mgl32.Scale3D(0.3, 0.3, 0).
			Mul4(mgl32.Translate3D(0.7, 0, 0)).
			Mul4(mgl32.HomogRotate3DZ(mgl32.DegToRad(angle)))

		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		colorProg.Use()
		gl.DrawArrays(gl.TRIANGLES, 0, 6)

		yellowProg.Use()
		gl.UniformMatrix4fv(yellowProg.Uniform("transform"), 1, false, &yTrans[0])
		gl.DrawArrays(gl.TRIANGLES, 0, 3)

		window.SwapBuffers()
		glfw.PollEvents()
		frames++

		previousTime = t
	}
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

	loop()
}
