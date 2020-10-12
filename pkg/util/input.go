package util

import "github.com/go-gl/glfw/v3.3/glfw"

// KeyAction shortcut to do something if a key is pressed
func KeyAction(w *glfw.Window, k glfw.Key, a func()) {
	if w.GetKey(k) == glfw.Press {
		a()
	}
}
