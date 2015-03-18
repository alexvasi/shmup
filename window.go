package main

import (
	"fmt"

	"github.com/go-gl/glfw/v3.0/glfw"
)

func InitGLFW() {
	if !glfw.Init() {
		panic("Failed to initialize GLFW")
	}

	glfw.SetErrorCallback(func(code glfw.ErrorCode, desc string) {
		panic(fmt.Sprint("GLFW error", code, desc))
	})

	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 2)
	glfw.WindowHint(glfw.OpenglProfile, glfw.OpenglCoreProfile)
	glfw.WindowHint(glfw.OpenglForwardCompatible, glfw.True)
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.Samples, 4)
}

func NewWindow(width int, height int, title string) *glfw.Window {
	glfw.WindowHint(glfw.Visible, glfw.False)

	monitor, err := glfw.GetPrimaryMonitor()
	PanicOnError(err)

	mode, err := monitor.GetVideoMode()
	PanicOnError(err)

	window, err := glfw.CreateWindow(width, height, title, nil, nil)
	PanicOnError(err)

	window.SetPosition((mode.Width-width)/2, (mode.Height-height)/2)
	window.Show()
	window.MakeContextCurrent()

	return window
}
