package main

import "github.com/go-gl/glfw/v3.0/glfw"

type Input struct {
	window *glfw.Window
	useJoy bool
}

func NewInput(w *glfw.Window) Input {
	return Input{
		window: w,
		useJoy: glfw.JoystickPresent(glfw.Joystick1),
	}
}

func (i *Input) Process() {
	glfw.PollEvents()

	if i.IsPressed(glfw.KeyEscape) {
		i.window.SetShouldClose(true)
	}
}

func (i Input) IsPressed(keys ...glfw.Key) bool {
	for _, key := range keys {
		if i.window.GetKey(key) == glfw.Press {
			return true
		}
	}
	return false
}
