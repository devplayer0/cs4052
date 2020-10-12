package app

import (
	"fmt"
	"log"

	"github.com/devplayer0/cs4052/pkg/util"
	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	mouseSensitivity = 5
	movementSpeed    = 5
)

// App represents the graphics application
type App struct {
	window *glfw.Window

	crosshair *util.Program
	colorTri  *util.Program
	yellowTri *util.Program

	previousTime float64

	lastDebug float64
	frames    uint32
	ocx, ocy  float64
	d         float32

	fov float32

	projection              mgl32.Mat4
	camera                  *util.Camera
	colorModel, yellowModel mgl32.Mat4
}

// NewApp creates a new app for the window
func NewApp(w *glfw.Window) *App {
	a := &App{
		window: w,

		colorModel:  mgl32.Translate3D(0, 2, 0),
		yellowModel: mgl32.HomogRotate3DX(mgl32.DegToRad(-90)),

		fov:    45,
		camera: util.NewCamera(mgl32.Vec3{0, 2, 5}, mgl32.Vec2{-90, 0}, true),
	}
	a.ocx, a.ocy = w.GetCursorPos()
	a.updateProjection()

	return a
}

// Setup sets up the application (compile shaders, load models etc.)
func (a *App) Setup() error {
	a.window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
	a.window.SetKeyCallback(a.onKeyEvent)
	a.window.SetCursorPosCallback(a.onCursorMove)

	a.crosshair = util.NewProgram()
	if err := a.crosshair.LinkFiles("assets/shaders/crosshair.vs", "assets/shaders/crosshair.fs"); err != nil {
		return fmt.Errorf("failed to set up crosshair program: %w", err)
	}

	vertexBuf := util.NewBuffer()
	vertexBuf.SetVec2([]mgl32.Vec2{
		{-1, 0},
		{1, 0},

		{0, -1},
		{0, 1},
	})
	a.crosshair.LinkVertexPointer("vPosition", 2, gl.FLOAT, vertexBuf, 0)

	wi, hi := a.window.GetSize()
	w := float32(wi)
	h := float32(hi)
	a.crosshair.SetUniformMat4("projection", mgl32.Ortho2D(0, w, 0, h))
	a.crosshair.SetUniformMat4("model", mgl32.Translate3D(w/2, h/2, 0).Mul4(mgl32.Scale3D(10, 10, 1)))

	a.colorTri = util.NewProgram()
	if err := a.colorTri.LinkFiles("assets/shaders/color_3d.vs", "assets/shaders/color_tri.fs"); err != nil {
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
	a.colorTri.LinkVertexPointer("vPosition", 3, gl.FLOAT, vertexBuf, 0)

	colorBuf := util.NewBuffer()
	colorBuf.SetVec4([]mgl32.Vec4{
		{0, 1, 0, 1},
		{1, 0, 0, 1},
		{0, 0, 1, 1},

		{1, 0, 0, 1},
		{0, 0, 1, 1},
		{0, 1, 0, 1},
	})
	a.colorTri.LinkVertexPointer("vColor", 4, gl.FLOAT, colorBuf, 0)

	a.yellowTri = util.NewProgram()
	if err := a.yellowTri.LinkFiles("assets/shaders/yellow_3d.vs", "assets/shaders/yellow_tri.fs"); err != nil {
		return fmt.Errorf("failed to set up yellow triangle program: %w", err)
	}

	vertexBuf = util.NewBuffer()
	vertexBuf.SetVec3([]mgl32.Vec3{
		{-1, -1, 0},
		{1, -1, 0},
		{0, 1, 0},
	})
	a.yellowTri.LinkVertexPointer("vPosition", 3, gl.FLOAT, vertexBuf, 0)

	//var err error
	//monkeyModel, err = util.NewModel("assets/meshes/monkey.obj")
	//if err != nil {
	//	return fmt.Errorf("failed to load mesh: %w", err)
	//}

	gl.Enable(gl.DEPTH_TEST)

	return nil
}

func (a *App) onCursorMove(w *glfw.Window, xpos, ypos float64) {
	dx := float32(xpos - a.ocx)
	dy := float32(ypos - a.ocy)

	a.camera.SetRotation(a.camera.Rotation().Add(mgl32.Vec2{
		mouseSensitivity * dx * a.d,
		-mouseSensitivity * dy * a.d,
	}))

	a.ocx = xpos
	a.ocy = ypos
}
func (a *App) onKeyEvent(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	switch key {
	case glfw.KeyEscape, glfw.KeyQ:
		a.window.Destroy()
	}
}

func (a *App) updateProjection() {
	w, h := a.window.GetSize()
	a.projection = mgl32.Perspective(mgl32.DegToRad(a.fov), float32(w)/float32(h), 0.1, 100)
}
func (a *App) readInputs() {
	util.KeyAction(a.window, glfw.KeyMinus, func() {
		a.fov -= 50 * a.d
		a.updateProjection()
	})
	util.KeyAction(a.window, glfw.KeyEqual, func() {
		a.fov += 50 * a.d
		a.updateProjection()
	})

	util.KeyAction(a.window, glfw.KeyW, func() {
		a.camera.MoveZ(movementSpeed * a.d)
	})
	util.KeyAction(a.window, glfw.KeyS, func() {
		a.camera.MoveZ(-movementSpeed * a.d)
	})

	util.KeyAction(a.window, glfw.KeyA, func() {
		a.camera.MoveX(-movementSpeed * a.d)
	})
	util.KeyAction(a.window, glfw.KeyD, func() {
		a.camera.MoveX(movementSpeed * a.d)
	})

	util.KeyAction(a.window, glfw.KeySpace, func() {
		a.camera.MoveY(movementSpeed * a.d)
	})
	util.KeyAction(a.window, glfw.KeyC, func() {
		a.camera.MoveY(-movementSpeed * a.d)
	})
}

func (a *App) draw() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	a.colorTri.Use()
	a.colorTri.SetUniformMat4("projection", a.projection)
	a.colorTri.SetUniformMat4("camera", a.camera.Transform())
	a.colorTri.SetUniformMat4("model", a.colorModel)
	gl.DrawArrays(gl.TRIANGLES, 0, 6)

	a.yellowTri.Use()
	a.yellowTri.SetUniformMat4("projection", a.projection)
	a.yellowTri.SetUniformMat4("camera", a.camera.Transform())
	a.yellowTri.SetUniformMat4("model", a.yellowModel)
	gl.DrawArrays(gl.TRIANGLES, 0, 3)

	a.crosshair.Use()
	gl.DrawArrays(gl.LINES, 0, 4)

	a.window.SwapBuffers()
}

// Update updates the app state and draws to the screen
func (a *App) Update() {
	t := glfw.GetTime()
	a.d = float32(t - a.previousTime)

	if t-a.lastDebug > 1 {
		log.Printf("FPS: %v", a.frames)
		log.Printf("FOV: %v", a.fov)
		log.Printf("Camera position: %v", a.camera.Position)
		log.Printf("Camera rotation: %v", a.camera.Rotation())

		a.frames = 0
		a.lastDebug = t
	}

	a.readInputs()

	a.draw()

	// Post-update
	glfw.PollEvents()
	a.frames++
	a.previousTime = t
}
