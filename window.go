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
	glfw.WindowHint(glfw.RefreshRate, 60)
	glfw.WindowHint(glfw.Samples, 4)

	return nil
}

func NewWindow(w int, h int, title string, fullscreen bool) (*glfw.Window, mgl.Vec2) {
	monitor := glfw.GetPrimaryMonitor()
	mode := monitor.GetVideoMode()

	if fullscreen {
		w, h = mode.Width, mode.Height
	} else {
		glfw.WindowHint(glfw.Visible, glfw.False)
		monitor = nil
	}

	window, err := glfw.CreateWindow(w, h, title, monitor, nil)
	PanicOnError(err)

	if !fullscreen {
		window.SetPos((mode.Width-w)/2, (mode.Height-h)/2)
		window.Show()
	}

	window.MakeContextCurrent()
	return window, mgl.Vec2{float32(w), float32(h)}
}

func NewFullScreenWindow(title string) (*glfw.Window, mgl.Vec2) {
	monitor := glfw.GetPrimaryMonitor()
	mode := monitor.GetVideoMode()

	glfw.WindowHint(glfw.RedBits, mode.RedBits)
	glfw.WindowHint(glfw.GreenBits, mode.GreenBits)
	glfw.WindowHint(glfw.BlueBits, mode.BlueBits)

	window, err := glfw.CreateWindow(mode.Width, mode.Height, title, monitor, nil)
	PanicOnError(err)
	window.MakeContextCurrent()

	return window, mgl.Vec2{float32(mode.Width), float32(mode.Height)}
}
