package main

import (
	"github.com/go-gl/glfw/v3.0/glfw"
	mgl "github.com/go-gl/mathgl/mgl32"
)

type Input struct {
	window *glfw.Window
	useJoy bool
	debug  bool
	dir    mgl.Vec2
	fire   bool
}

type KeyCallback func(*glfw.Window, glfw.Key, int, glfw.Action, glfw.ModifierKey)

func NewInput(w *glfw.Window) *Input {
	i := &Input{
		window: w,
		useJoy: glfw.JoystickPresent(glfw.Joystick1),
	}
	w.SetKeyCallback(i.GetKeyCallback())
	return i
}

func (i *Input) Process() {
	glfw.PollEvents()

	if i.IsPressed(glfw.KeyEscape) {
		i.window.SetShouldClose(true)
	}

	i.dir = mgl.Vec2{0, 0}
	i.fire = false

	axes, err := glfw.GetJoystickAxes(glfw.Joystick1)
	if err == nil && len(axes) > 1 {
		i.dir[0] = axes[0]
		i.dir[1] = -1 * axes[1]
	}

	buttons, err := glfw.GetJoystickButtons(glfw.Joystick1)
	if err == nil && len(buttons) > 14 {
		if buttons[11] > 0 || buttons[14] > 0 {
			i.fire = true
		}
	}

	if i.IsPressed(glfw.KeyUp, glfw.KeyW) {
		i.dir[1] = 1
	} else if i.IsPressed(glfw.KeyDown, glfw.KeyS) {
		i.dir[1] = -1
	}

	if i.IsPressed(glfw.KeyLeft, glfw.KeyA) {
		i.dir[0] = -1
	} else if i.IsPressed(glfw.KeyRight, glfw.KeyD) {
		i.dir[0] = 1
	}

	if i.IsPressed(glfw.KeySpace, glfw.KeyZ) {
		i.fire = true
	}
}

func (i *Input) IsPressed(keys ...glfw.Key) bool {
	for _, key := range keys {
		if i.window.GetKey(key) == glfw.Press {
			return true
		}
	}
	return false
}

func (i *Input) GetKeyCallback() KeyCallback {
	cb := func(w *glfw.Window, key glfw.Key, scan int, action glfw.Action,
		m glfw.ModifierKey) {
		if key == glfw.KeyGraveAccent && action == glfw.Press {
			i.debug = !i.debug
		}
	}
	return cb
}
