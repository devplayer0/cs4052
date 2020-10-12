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

const (
	mouseSensitivity = 5
	movementSpeed    = 5
)

var window *glfw.Window

var crosshairProg, colorProg, yellowProg *util.Program
var monkeyModel *util.Model

func setup() error {
	crosshairProg = util.NewProgram()
	if err := crosshairProg.LinkFiles("assets/shaders/crosshair.vs", "assets/shaders/crosshair.fs"); err != nil {
		return fmt.Errorf("failed to set up crosshair program: %w", err)
	}

	vertexBuf := util.NewBuffer()
	vertexBuf.SetVec2([]mgl32.Vec2{
		{-1, 0},
		{1, 0},

		{0, -1},
		{0, 1},
	})
	crosshairProg.LinkVertexPointer("vPosition", 2, gl.FLOAT, vertexBuf, 0)
	crosshairProg.SetUniformMat4("projection", mgl32.Ortho2D(0, 800, 0, 600))
	crosshairProg.SetUniformMat4("model", mgl32.Translate3D(400, 300, 0).Mul4(mgl32.Scale3D(10, 10, 1)))

	colorProg = util.NewProgram()
	if err := colorProg.LinkFiles("assets/shaders/color_3d.vs", "assets/shaders/color_tri.fs"); err != nil {
		return fmt.Errorf("failed to set up color triangle program: %w", err)
	}

	vertexBuf = util.NewBuffer()
	vertexBuf.SetVec3([]mgl32.Vec3{
		{-1, -1, 0},
		{1, -1, 0},
		{0, 0, 0},

		{1, 1, 0},
		{-1, 1, 0},
		{0, 0, 0},
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
	if err := yellowProg.LinkFiles("assets/shaders/yellow_3d.vs", "assets/shaders/yellow_tri.fs"); err != nil {
		return fmt.Errorf("failed to set up yellow triangle program: %w", err)
	}

	vertexBuf = util.NewBuffer()
	vertexBuf.SetVec3([]mgl32.Vec3{
		{-1, -1, 0},
		{1, -1, 0},
		{0, 1, 0},
	})
	yellowProg.LinkVertexPointer("vPosition", 3, gl.FLOAT, vertexBuf, 0)

	var err error
	monkeyModel, err = util.NewModel("assets/meshes/monkey.obj")
	if err != nil {
		return fmt.Errorf("failed to load mesh: %w", err)
	}

	window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)

	return nil
}

func loop() {
	frames := 0
	lastFPS := glfw.GetTime()
	previousTime := lastFPS
	var d float32

	w, h := window.GetSize()
	fov := float32(45.0)

	camera := util.NewCamera(mgl32.Vec3{0, 2, 5}, mgl32.Vec2{-90, 0}, true)

	colorModel := mgl32.Translate3D(0, 2, 0)
	yellowModel := mgl32.HomogRotate3DX(mgl32.DegToRad(-90))

	window.SetKeyCallback(func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		switch key {
		case glfw.KeyEscape, glfw.KeyQ:
			window.Destroy()
		}
	})

	ocx, ocy := window.GetCursorPos()
	window.SetCursorPosCallback(func(w *glfw.Window, xpos, ypos float64) {
		dx := float32(xpos - ocx)
		dy := float32(ypos - ocy)

		camera.SetRotation(camera.Rotation().Add(mgl32.Vec2{
			mouseSensitivity * dx * d,
			-mouseSensitivity * dy * d,
		}))

		ocx = xpos
		ocy = ypos
	})

	for !window.ShouldClose() {
		// Update
		t := glfw.GetTime()
		d = float32(t - previousTime)

		util.KeyAction(window, glfw.KeyMinus, func() {
			fov -= 100 * d
		})
		util.KeyAction(window, glfw.KeyEqual, func() {
			fov += 100 * d
		})

		util.KeyAction(window, glfw.KeyW, func() {
			camera.MoveZ(movementSpeed * d)
		})
		util.KeyAction(window, glfw.KeyS, func() {
			camera.MoveZ(-movementSpeed * d)
		})

		util.KeyAction(window, glfw.KeyA, func() {
			camera.MoveX(-movementSpeed * d)
		})
		util.KeyAction(window, glfw.KeyD, func() {
			camera.MoveX(movementSpeed * d)
		})

		util.KeyAction(window, glfw.KeySpace, func() {
			camera.MoveY(movementSpeed * d)
		})
		util.KeyAction(window, glfw.KeyC, func() {
			camera.MoveY(-movementSpeed * d)
		})

		projection := mgl32.Perspective(mgl32.DegToRad(fov), float32(w)/float32(h), 0.1, 100)

		if t-lastFPS > 1 {
			log.Printf("FPS: %v", frames)
			log.Printf("FOV: %v", fov)
			log.Printf("Camera position: %v", camera.Position)
			log.Printf("Camera rotation: %v", camera.Rotation())

			frames = 0
			lastFPS = t
		}

		// Draw
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		colorProg.Use()
		colorProg.SetUniformMat4("projection", projection)
		colorProg.SetUniformMat4("camera", camera.Transform())
		colorProg.SetUniformMat4("model", colorModel)
		gl.DrawArrays(gl.TRIANGLES, 0, 6)

		yellowProg.Use()
		yellowProg.SetUniformMat4("projection", projection)
		yellowProg.SetUniformMat4("camera", camera.Transform())
		yellowProg.SetUniformMat4("model", yellowModel)
		gl.DrawArrays(gl.TRIANGLES, 0, 3)

		crosshairProg.Use()
		gl.DrawArrays(gl.LINES, 0, 4)

		window.SwapBuffers()

		// Post-update
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

	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Printf("OpenGL version: %v", version)

	gl.Enable(gl.DEBUG_OUTPUT)
	gl.DebugMessageCallback(func(source, gltype, id, severity uint32, length int32, message string, userParam unsafe.Pointer) {
		log.Printf("[SEV%v] %v", severity, message)
	}, nil)

	gl.Enable(gl.DEPTH_TEST)

	if err := setup(); err != nil {
		log.Fatalf("Setup failed: %v", err)
	}

	loop()
}
