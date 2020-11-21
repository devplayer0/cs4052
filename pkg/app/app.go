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

	crosshair *Crosshair

	lighting *util.Lighting
	backpack *util.Mesh

	previousTime float64

	lastDebug float64
	frames    uint32
	ocx, ocy  float64
	d         float32

	fov       float32
	wireframe bool

	projection mgl32.Mat4
	camera     *util.Camera

	brrLamp       util.Lamp
	brrLampOrbit  mgl32.Vec3
	brrLampAngle  float32
	backpackTrans mgl32.Mat4
}

// NewApp creates a new app for the window
func NewApp(w *glfw.Window) *App {
	a := &App{
		window: w,

		brrLampOrbit:  mgl32.Vec3{3, 8, 1},
		backpackTrans: mgl32.Translate3D(3, 2, 0),

		fov:    45,
		camera: util.NewCamera(mgl32.Vec3{0, 2, 10}, mgl32.Vec2{-90, 0}, true),
	}
	a.updateProjection()

	return a
}

// Setup sets up the application (compile shaders, load models etc.)
func (a *App) Setup() error {
	a.window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
	a.window.SetKeyCallback(a.onKeyEvent)
	a.window.SetCursorPosCallback(a.onCursorMove)

	var err error
	a.crosshair, err = NewCrosshair(a.window)
	if err != nil {
		return fmt.Errorf("failed to set up crosshair: %w", err)
	}

	a.brrLamp = util.Lamp{
		Position: a.brrLampOrbit,

		Ambient:     mgl32.Vec3{0.05, 0.00, 0.0},
		Diffuse:     mgl32.Vec3{0.8, 0.0, 0.0},
		Specular:    mgl32.Vec3{1, 0, 0},
		Attenuation: util.AttenuationParams{1, 0.09, 0.032},
	}

	a.lighting, err = util.NewLightingVSFile("assets/shaders/model.vs", []*util.Lamp{
		{
			Position: mgl32.Vec3{-2, 5, -2},

			Ambient:     mgl32.Vec3{0.01, 0.02, 0.04},
			Diffuse:     mgl32.Vec3{0.2, 0.3, 0.7},
			Specular:    mgl32.Vec3{0.18, 0.4, 0.83},
			Attenuation: util.AttenuationParams{1, 0.09, 0.032},
		},
		{
			Position: mgl32.Vec3{6, -2, 5},

			Ambient:     mgl32.Vec3{0.02, 0.05, 0.01},
			Diffuse:     mgl32.Vec3{0.3, 0.8, 0.15},
			Specular:    mgl32.Vec3{0.4, 1, 0.2},
			Attenuation: util.AttenuationParams{1, 0.09, 0.032},
		},
		&a.brrLamp,
	})
	if err != nil {
		return fmt.Errorf("failed to initialize lighting: %w", err)
	}

	a.backpack, err = util.NewMesh("assets/meshes/backpack.obj")
	if err != nil {
		return fmt.Errorf("failed to load mesh: %w", err)
	}
	a.backpack.Upload(a.lighting.Shader)

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LEQUAL)

	a.ocx, a.ocy = a.window.GetCursorPos()

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
	if action == glfw.Release {
		switch key {
		case glfw.KeyF:
			var m uint32
			m = gl.LINE
			if a.wireframe {
				m = gl.FILL
			}

			gl.PolygonMode(gl.FRONT_AND_BACK, m)
			a.wireframe = !a.wireframe
		}
	}

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

	a.backpack.Draw(a.lighting.Shader, a.projection, a.camera, a.backpackTrans)

	a.lighting.DrawCubes(a.projection, a.camera)

	a.crosshair.Draw()

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

	brrLampTransform := util.TransFromPos(a.brrLampOrbit).Mul4(mgl32.HomogRotate3DY(a.brrLampAngle)).Mul4(mgl32.Translate3D(0, 0, -5))
	a.brrLamp.Position = util.PosFromTrans(brrLampTransform)
	a.lighting.Update()

	a.brrLampAngle += 4 * a.d
	if a.brrLampAngle > 360 {
		a.brrLampAngle = 0
	}

	a.draw()

	// Post-update
	glfw.PollEvents()
	a.frames++
	a.previousTime = t
}
