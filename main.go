package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"time"

	"github.com/go-gl/gl/v3.2-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
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

	fullscreen := flag.Bool("fs", false, "fullscreen mode")
	flag.Parse()

	defer HandlePanic()

	PanicOnError(InitGLFW())
	defer glfw.Terminate()

	window, screenSize := NewWindow(width, height, title, *fullscreen)
	PanicOnError(gl.Init())

	InitSound()
	defer TerminateSound()
	LoadSoundAssets("*.wav")

	input := NewInput(window, *fullscreen)
	renderer := NewRenderer(width, height, screenSize)
	world := NewWorld(width, height)
	game := NewGame(world, input)
	timer := NewTimer()

	for !window.ShouldClose() {
		renderer.Clear()
		game.Update(timer.DT)
		world.Draw(renderer)
		renderer.Render()

		window.SwapBuffers()

		timer.Tick()
		if timer.TicksCount == 60 {
			window.SetTitle(timer.Stat())
			timer.ResetCounter()
		}

		input.Process()
		switch {
		case input.DebugToggled && input.Debug:
			glfw.SwapInterval(0)
		case input.DebugToggled:
			glfw.SwapInterval(1)
		case input.FullscreenToggled:
			//renderer.Cleanup()
			window.Destroy()
			window, screenSize = NewWindow(width, height, title,
				input.Fullscreen)
			input.SetWindow(window)
			renderer = NewRenderer(width, height, screenSize)
		}
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
