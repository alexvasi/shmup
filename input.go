package main

import (
	"github.com/go-gl/glfw/v3.1/glfw"
	mgl "github.com/go-gl/mathgl/mgl32"
)

type Input struct {
	window *glfw.Window
	useJoy bool

	Dir               mgl.Vec2
	Fire              bool
	Debug             bool
	DebugToggled      bool
	Fullscreen        bool
	FullscreenToggled bool
}

func NewInput(w *glfw.Window, fullscreen bool) *Input {
	i := &Input{
		useJoy:     glfw.JoystickPresent(glfw.Joystick1),
		Fullscreen: fullscreen,
	}
	i.SetWindow(w)
	return i
}

func (i *Input) Process() {
	i.DebugToggled = false
	i.FullscreenToggled = false
	glfw.PollEvents()

	if i.IsPressed(glfw.KeyEscape) {
		i.window.SetShouldClose(true)
	}

	i.Dir = mgl.Vec2{0, 0}
	i.Fire = false

	axes := glfw.GetJoystickAxes(glfw.Joystick1)
	if len(axes) > 1 {
		i.Dir[0] = axes[0]
		i.Dir[1] = -1 * axes[1]
	}

	buttons := glfw.GetJoystickButtons(glfw.Joystick1)
	if len(buttons) > 14 {
		if buttons[11] > 0 || buttons[14] > 0 {
			i.Fire = true
		}
	}

	if i.IsPressed(glfw.KeyUp, glfw.KeyW) {
		i.Dir[1] = 1
	} else if i.IsPressed(glfw.KeyDown, glfw.KeyS) {
		i.Dir[1] = -1
	}

	if i.IsPressed(glfw.KeyLeft, glfw.KeyA) {
		i.Dir[0] = -1
	} else if i.IsPressed(glfw.KeyRight, glfw.KeyD) {
		i.Dir[0] = 1
	}

	if i.IsPressed(glfw.KeySpace, glfw.KeyZ) {
		i.Fire = true
	}
}

func (i *Input) SetWindow(w *glfw.Window) {
	i.window = w
	w.SetKeyCallback(i.GetKeyCallback())
	w.SetInputMode(glfw.CursorMode, glfw.CursorHidden)
}

func (i *Input) IsPressed(keys ...glfw.Key) bool {
	for _, key := range keys {
		if i.window.GetKey(key) == glfw.Press {
			return true
		}
	}
	return false
}

func (i *Input) GetKeyCallback() glfw.KeyCallback {
	cb := func(w *glfw.Window, key glfw.Key, scan int, action glfw.Action,
		m glfw.ModifierKey) {

		if key == glfw.KeyGraveAccent && action == glfw.Press {
			i.Debug = !i.Debug
			i.DebugToggled = true
		} else if key == glfw.KeyF && action == glfw.Press {
			i.Fullscreen = !i.Fullscreen
			i.FullscreenToggled = true
		}
	}
	return cb
}
