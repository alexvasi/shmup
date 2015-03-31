package main

import (
	"github.com/go-gl/glfw/v3.1/glfw"
	mgl "github.com/go-gl/mathgl/mgl32"
)

func InitGLFW() error {
	err := glfw.Init()
	if err != nil {
		return err
	}

	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 2)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.Samples, 4)

	return nil
}

func NewWindow(width int, height int, title string) *glfw.Window {
	glfw.WindowHint(glfw.Visible, glfw.False)

	monitor := glfw.GetPrimaryMonitor()
	mode := monitor.GetVideoMode()

	window, err := glfw.CreateWindow(width, height, title, nil, nil)
	PanicOnError(err)

	window.SetPos((mode.Width-width)/2, (mode.Height-height)/2)
	window.Show()
	window.MakeContextCurrent()

	return window
}

func NewFullScreenWindow(title string) (*glfw.Window, mgl.Vec2) {
	monitor := glfw.GetPrimaryMonitor()
	mode := monitor.GetVideoMode()

	window, err := glfw.CreateWindow(mode.Width, mode.Height, title, monitor, nil)
	PanicOnError(err)
	window.MakeContextCurrent()

	return window, mgl.Vec2{float32(mode.Width), float32(mode.Height)}
}
