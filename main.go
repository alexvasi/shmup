package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/go-gl/gl/v3.2-core/gl"
	"github.com/go-gl/glfw/v3.0/glfw"
)

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

func main() {
	const (
		title  = "Shmup"
		width  = 1366
		height = 768
	)

	defer HandlePanic()

	InitGLFW()
	defer glfw.Terminate()

	window := NewWindow(width, height, title)
	PanicOnError(gl.Init())

	input := NewInput(window)
	renderer := NewRenderer(width, height, width, height)
	ship := NewShip()
	missiles := []*Missile{}

	timer := NewTimer(window)
	for !window.ShouldClose() {
		timer.Tick()

		input.Process()
		timer.ShowTimings(input.debug, 60)

		ship.Thrust(input.dir, input.fire)
		newMissiles := ship.Update(timer.DT)
		if len(newMissiles) > 0 {
			missiles = append(missiles, newMissiles...)
		}

		for _, m := range missiles {
			m.Update(timer.DT)
		}

		renderer.Clear()
		ship.Draw(renderer)
		for _, m := range missiles {
			m.Draw(renderer)
		}
		renderer.Render()

		window.SwapBuffers()
	}
}

func HandlePanic() {
	if err := recover(); err != nil {
		fmt.Println(err)
		glfw.Terminate()
		os.Exit(1)
	}
}

func PanicOnError(err error) {
	if err != nil {
		panic(err)
	}
}
