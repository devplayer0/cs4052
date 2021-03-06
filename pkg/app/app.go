package app

import (
	"fmt"
	"log"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/devplayer0/cs4052/pkg/object"
	"github.com/devplayer0/cs4052/pkg/util"
)

const (
	mouseSensitivity = 5
	movementSpeed    = 5
	boidCount        = 64
	brrLamp2Speed    = 3
)

type skinnedVSParams struct {
	DepthPass bool
}

// App represents the graphics application
type App struct {
	window *glfw.Window

	crosshair *Crosshair

	lighting *util.Lighting

	meshShader             *util.Program
	meshDepthShader        *util.Program
	skinnedMeshShader      *util.Program
	skinnedMeshDepthShader *util.Program
	skeletonShader         *util.Program

	skybox  *util.Skybox
	skybox2 *util.Skybox

	ground    *object.Mesh
	backpack  *object.Mesh
	scorpion  *object.Object
	tarantula *object.Object
	locust    *object.Object

	previousTime  float64
	animationTime float32

	lastDebug float64
	frames    uint32
	ocx, ocy  float64
	d         float32

	fov    float32
	paused bool

	projection mgl32.Mat4
	camera     *util.Camera

	depthMapsFirstPass bool
	brrLamp            util.Lamp
	brrLampOrbit       mgl32.Vec3
	brrLampAngle       float32
	brrLamp2           util.Lamp
	brrLamp2Dir        float32

	spotlight util.Spotlight

	scorpionTrans  mgl32.Mat4
	tarantulaTrans mgl32.Mat4
	locustTrans    mgl32.Mat4

	boids *object.Boids
}

// NewApp creates a new app for the window
func NewApp(w *glfw.Window) *App {
	a := &App{
		window: w,
		camera: util.NewCamera(mgl32.Vec3{0, 10, 11}, mgl32.Vec2{-90, -25}, true),

		brrLampOrbit: mgl32.Vec3{3, 10, 1},

		fov:    45,
		paused: false,

		depthMapsFirstPass: true,
		brrLamp2Dir:        float32(brrLamp2Speed),

		scorpionTrans:  mgl32.Translate3D(6, 0, 0).Mul4(mgl32.Scale3D(0.02, 0.02, 0.02)),
		tarantulaTrans: mgl32.Translate3D(0, 0, -2).Mul4(mgl32.Scale3D(0.04, 0.04, 0.04)),
		locustTrans:    mgl32.Translate3D(-8, 0, -2).Mul4(mgl32.Scale3D(0.01, 0.01, 0.01)),
	}

	wi, hi := w.GetSize()
	a.ocx, a.ocy = float64(wi)/2, float64(hi)/2
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

	a.skybox, err = util.NewSkybox("assets/skyboxes/mountains/")
	if err != nil {
		return fmt.Errorf("failed to set up skybox: %w", err)
	}
	a.skybox2, err = util.NewSkybox("assets/skyboxes/city/")
	if err != nil {
		return fmt.Errorf("failed to set up skybox 2: %w", err)
	}

	att := util.AttenuationParams{Constant: 1, Linear: 0.07, Quadratic: 0.017}
	a.brrLamp = util.Lamp{
		Position: a.brrLampOrbit,

		Ambient:     mgl32.Vec3{0.4, 0.0, 0.0},
		Diffuse:     mgl32.Vec3{1, 0.0, 0.0},
		Specular:    mgl32.Vec3{1, 0, 0},
		Attenuation: att,
	}
	a.brrLamp2 = util.Lamp{
		Position: mgl32.Vec3{6, 3, -1},

		Ambient:     mgl32.Vec3{0.02, 0.05, 0.01},
		Diffuse:     mgl32.Vec3{0.3, 0.8, 0.15},
		Specular:    mgl32.Vec3{0.4, 1, 0.2},
		Attenuation: att,
	}
	a.spotlight = util.Spotlight{
		Cutoff:      util.Cos(mgl32.DegToRad(12.5)),
		OuterCutoff: util.Cos(mgl32.DegToRad(15)),

		Ambient:  mgl32.Vec3{0, 0, 0},
		Diffuse:  mgl32.Vec3{1, 1, 1},
		Specular: mgl32.Vec3{2, 2, 2},

		Attenuation: att,
	}

	a.lighting, err = util.NewLighting([]*util.DirectionalLight{
		{
			Direction: mgl32.Vec3{-0.2, -1, -0.3},
			Ambient:   mgl32.Vec3{0.1, 0.1, 0.1},
			Diffuse:   mgl32.Vec3{0.7, 0.7, 0.7},
			Specular:  mgl32.Vec3{0.5, 0.5, 0.5},
		},
	}, []*util.Lamp{
		{
			Position: mgl32.Vec3{-2, 8, -2},

			Ambient:     mgl32.Vec3{0.01, 0.02, 0.04},
			Diffuse:     mgl32.Vec3{0.2, 0.3, 0.7},
			Specular:    mgl32.Vec3{0.18, 0.4, 0.83},
			Attenuation: att,
		},
		&a.brrLamp,
		&a.brrLamp2,
		{
			Position: mgl32.Vec3{-4, 6, 1},

			Ambient:     mgl32.Vec3{0.05, 0.05, 0.05},
			Diffuse:     mgl32.Vec3{0.8, 0.8, 0.8},
			Specular:    mgl32.Vec3{0.4, 0.4, 0.4},
			Attenuation: att,
		},

		{
			Position: mgl32.Vec3{-28, 2, -28},

			Diffuse:     mgl32.Vec3{0.9, 0.9, 0.9},
			Specular:    mgl32.Vec3{0.4, 0.4, 0.4},
			Attenuation: att,
		},
		{
			Position: mgl32.Vec3{28, 2, -28},

			Diffuse:     mgl32.Vec3{0.9, 0.9, 0.9},
			Specular:    mgl32.Vec3{0.4, 0.4, 0.4},
			Attenuation: att,
		},
		{
			Position: mgl32.Vec3{-28, 2, 28},

			Diffuse:     mgl32.Vec3{0.9, 0.9, 0.9},
			Specular:    mgl32.Vec3{0.4, 0.4, 0.4},
			Attenuation: att,
		},
	}, []*util.Spotlight{
		&a.spotlight,
	})
	if err != nil {
		return fmt.Errorf("failed to initialize lighting: %w", err)
	}

	a.meshShader, err = a.lighting.ProgramVSFile("assets/shaders/mesh.vs")
	if err != nil {
		return fmt.Errorf("failed to link mesh shaders: %w", err)
	}
	a.meshDepthShader, err = a.lighting.DepthProgramVSFile("assets/shaders/shadows_depth.vs")
	if err != nil {
		return fmt.Errorf("failed to link mesh depth pass shaders: %w", err)
	}

	groundDiffuse := util.NewTexture(gl.TEXTURE_2D)
	if err := groundDiffuse.LoadJPEGFile(gl.TEXTURE_2D, "assets/textures/brickwall.jpg"); err != nil {
		return fmt.Errorf("failed to load ground texture: %w", err)
	}
	groundDiffuse.Apply2DDefaults()

	groundNormal := util.NewTexture(gl.TEXTURE_2D)
	if err := groundNormal.LoadJPEGFile(gl.TEXTURE_2D, "assets/textures/brickwall_normal.jpg"); err != nil {
		return fmt.Errorf("failed to load ground normal map texture: %w", err)
	}
	groundNormal.Apply2DDefaults()

	a.ground, err = object.NewOBJMeshFile("assets/meshes/plane.obj", &object.Material{
		Specular: mgl32.Vec3{0.05, 0.05, 0.05},

		DiffuseTexture: groundDiffuse,
		NormalTexture:  groundNormal,
	})
	if err != nil {
		return fmt.Errorf("failed to load mesh: %w", err)
	}
	a.ground.Upload(a.meshShader).LinkDepthMap(a.meshDepthShader)

	a.backpack, err = object.NewOBJMeshFile("assets/meshes/backpack.obj", &object.Material{
		Diffuse:        mgl32.Vec3{0, 0.1, 0},
		Reflectiveness: 0.7,
	})
	if err != nil {
		return fmt.Errorf("failed to load mesh: %w", err)
	}
	a.backpack.Upload(a.meshShader).LinkDepthMap(a.meshDepthShader)

	a.skinnedMeshShader, err = a.lighting.ProgramVSTemplateFile("assets/shaders/mesh_skinned.vs", skinnedVSParams{false})
	if err != nil {
		return fmt.Errorf("failed to link skinned mesh shaders: %w", err)
	}
	a.skinnedMeshDepthShader, err = a.lighting.DepthProgramVSTemplateFile("assets/shaders/mesh_skinned.vs", skinnedVSParams{true})
	if err != nil {
		return fmt.Errorf("failed to link skinned mesh depth shaders: %w", err)
	}

	a.skeletonShader = util.NewProgram()
	if err := a.skeletonShader.LinkFiles("assets/shaders/generic_3d.vs", "assets/shaders/uniform_color.fs", ""); err != nil {
		return fmt.Errorf("failed to setup skeleton debug shader: %w", err)
	}
	a.skeletonShader.SetUniformVec3("color", mgl32.Vec3{1, 0, 1})

	a.scorpion, err = object.NewObjectFile("assets/objects/scorpion.sobj", a.skinnedMeshShader, a.skinnedMeshDepthShader, a.skeletonShader)
	if err != nil {
		return fmt.Errorf("failed to set up scorpion: %w", err)
	}
	a.tarantula, err = object.NewObjectFile("assets/objects/tarantula.sobj", a.skinnedMeshShader, a.skinnedMeshDepthShader, a.skeletonShader)
	if err != nil {
		return fmt.Errorf("failed to set up tarantula: %w", err)
	}
	a.locust, err = object.NewObjectFile("assets/objects/locust.sobj", a.skinnedMeshShader, a.skinnedMeshDepthShader, a.skeletonShader)
	if err != nil {
		return fmt.Errorf("failed to set up locust: %w", err)
	}

	a.boids = object.NewBoids(util.Bounds{
		Min: mgl32.Vec3{-32, 0, -32},
		Max: mgl32.Vec3{32, 0, 32},
	}, 0.02)
	for i := 0; i < boidCount; i++ {
		a.boids.Instances = append(a.boids.Instances, a.boids.MakeBoid())
	}

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.Enable(gl.CULL_FACE)

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
		case glfw.KeyM:
			object.MeshWireFrame = !object.MeshWireFrame
		case glfw.KeyE:
			a.scorpion.Debug = !a.scorpion.Debug
		case glfw.KeyP:
			a.paused = !a.paused
		case glfw.KeyN:
			object.DisableNormalMapping = !object.DisableNormalMapping
		case glfw.KeyZ:
			a.lighting.ShadowsEnabled = !a.lighting.ShadowsEnabled
		case glfw.KeyX:
			tmp := a.skybox
			a.skybox = a.skybox2
			a.skybox2 = tmp
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
	groundTrans := mgl32.Translate3D(0, 0, 0).Mul4(mgl32.Scale3D(32, 32, 32))
	backpackTrans := mgl32.Translate3D(7, 4, -8)

	// Depth map pass
	a.lighting.ShadowsDepthPass(func(dpa util.DepthMapParamsApplicator) {
		a.ground.DepthMapPass(a.meshDepthShader, groundTrans, dpa)
		a.backpack.DepthMapPass(a.meshDepthShader, backpackTrans, dpa)

		a.scorpion.DepthMapPass(a.scorpionTrans, dpa)
		a.tarantula.DepthMapPass(a.tarantulaTrans, dpa)
		a.locust.DepthMapPass(a.locustTrans, dpa)
	})

	// Drawing pass
	w, h := a.window.GetSize()
	gl.Viewport(0, 0, int32(w), int32(h))
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	a.ground.Draw(a.meshShader, a.projection, a.camera, groundTrans, a.skybox.Texture, a.lighting.DepthMaps)
	a.backpack.Draw(a.meshShader, a.projection, a.camera, backpackTrans, a.skybox.Texture, a.lighting.DepthMaps)

	a.scorpion.Draw(a.projection, a.camera, a.scorpionTrans, a.skybox.Texture, a.lighting.DepthMaps)
	a.tarantula.Draw(a.projection, a.camera, a.tarantulaTrans, a.skybox.Texture, a.lighting.DepthMaps)
	a.locust.Draw(a.projection, a.camera, a.locustTrans, a.skybox.Texture, a.lighting.DepthMaps)

	boidBase := mgl32.Scale3D(0.01, 0.01, 0.01)
	for _, b := range a.boids.Instances {
		angle := util.Atan2(b.Velocity.Z(), b.Velocity.X())
		trans := mgl32.Translate3D(b.Position.X(), 0, b.Position.Z()).Mul4(mgl32.HomogRotate3DY(angle)).Mul4(boidBase)

		a.scorpion.Update(a.projection, a.camera, trans, a.scorpion.Animations[4], a.animationTime)
		a.scorpion.Draw(a.projection, a.camera, trans, a.skybox.Texture, a.lighting.DepthMaps)
	}

	a.lighting.DrawCubes(a.projection, a.camera)

	a.skybox.Draw(a.projection, a.camera)

	a.crosshair.Draw()

	a.window.SwapBuffers()
}

// Update updates the app state and draws to the screen
func (a *App) Update() {
	t := glfw.GetTime()
	a.d = float32(t - a.previousTime)
	if !a.paused {
		a.animationTime += float32(a.d)
		a.boids.Update()
	}

	if t-a.lastDebug > 1 {
		log.Printf("FPS: %v", a.frames)
		log.Printf("FOV: %v", a.fov)
		log.Printf("Camera position: %v", a.camera.Position)
		log.Printf("Camera rotation: %v", a.camera.Rotation())
		log.Printf("Camera direction: %v", a.camera.Direction())

		a.frames = 0
		a.lastDebug = t
	}

	a.readInputs()

	brrLampTransform := util.TransFromPos(a.brrLampOrbit).Mul4(mgl32.HomogRotate3DY(a.brrLampAngle)).Mul4(mgl32.Translate3D(0, 0, -5))
	a.brrLamp.Position = util.PosFromTrans(brrLampTransform)
	a.spotlight.Position = a.camera.Position
	a.spotlight.Direction = a.camera.Direction()

	a.brrLamp2.Position = mgl32.Vec3{a.brrLamp2.Position.X() + a.brrLamp2Dir*a.d, a.brrLamp2.Position.Y(), a.brrLamp2.Position.Z()}

	a.lighting.SetViewPos(a.camera.Position)
	a.lighting.Update(a.meshShader, a.skinnedMeshShader)

	if a.depthMapsFirstPass {
		a.lighting.UpdateLamps(a.meshShader, a.skinnedMeshShader)
		a.depthMapsFirstPass = false
	} else {
		a.lighting.UpdateLamp(&a.brrLamp, a.meshShader, a.skinnedMeshShader)
		a.lighting.UpdateLamp(&a.brrLamp2, a.meshShader, a.skinnedMeshShader)
	}

	a.brrLampAngle += 4 * a.d
	if a.brrLampAngle > 360 {
		a.brrLampAngle = 0
	}
	if a.brrLamp2.Position.X() > 12 {
		a.brrLamp2Dir = -float32(brrLamp2Speed)
	} else if a.brrLamp2.Position.X() < 6 {
		a.brrLamp2Dir = float32(brrLamp2Speed)
	}

	a.scorpion.Update(a.projection, a.camera, a.scorpionTrans, a.scorpion.Animations[0], a.animationTime)
	a.tarantula.Update(a.projection, a.camera, a.tarantulaTrans, nil, 0)
	a.locust.Update(a.projection, a.camera, a.locustTrans, nil, 0)

	a.draw()

	// Post-update
	glfw.PollEvents()
	a.frames++
	a.previousTime = t
}
