package main

import (
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"time"

	"github.com/go-gl/gl/v3.2-core/gl"
	"github.com/go-gl/glfw/v3.0/glfw"
)

func init() {
	runtime.LockOSThread() // GLFW event handling must run on the main OS thread
	rand.Seed(time.Now().UnixNano())
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
	ship := NewShip(Human, 100, height/2)
	ship2 := NewShip(Others, 1000, height/2+30)
	ship2.dir[0] = -1

	world := NewWorld(width, height)
	world.AddShips(ship, ship2)

	timer := NewTimer(window)
	for !window.ShouldClose() {
		timer.Tick()
		renderer.Clear()

		input.Process()
		timer.ShowTimings(input.debug, 60)
		ship.Control(input.dir, input.fire)

		world.Update(timer.DT)
		world.Draw(renderer)
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
