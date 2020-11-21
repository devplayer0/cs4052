package main

import (
	"log"
	"runtime"
	"unsafe"

	"github.com/devplayer0/cs4052/pkg/app"
	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

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
	window, err := glfw.CreateWindow(800, 600, "Midterm", nil, nil)
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

	app := app.NewApp(window)
	if err := app.Setup(); err != nil {
		log.Fatalf("Setup failed: %v", err)
	}

	for !window.ShouldClose() {
		app.Update()
	}
}
